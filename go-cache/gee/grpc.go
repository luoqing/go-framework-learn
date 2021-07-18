package gee

// grpc :https://zhuanlan.zhihu.com/p/85508384
import(
	"go-cache/consistenthash"
	"fmt"
	"google.golang.org/grpc"
    pb "go-cache/cacher"
	"errors"
	"context"
	"net"
	"time"
)

type RpcPool struct {
	self string // 分布式和本地保持统一
	replicas int 
	peersHash *consistenthash.ConsistentHash // 也可以是轮询的方式，主要实现的是Get(key) (ip, err) 根据key来hash到某个ip， 还有就是Add(...ips), 将服务器添加到均衡器
	pb.UnimplementedCacherServer
}

func NewRpcPool(self string, replicas int, ips []string) *RpcPool{
	p := &RpcPool{
		self: self,
		replicas: replicas,
	}
	p.peersHash = consistenthash.New(replicas, nil)
	p.peersHash.Add(ips...)
	return p
}

func (p *RpcPool)PickPeer(key string) (PeerGetter, error) {
	ip, err := p.peersHash.Get(key)
	if err != nil {
		return nil, err
	}
	fmt.Printf("ip:%s self:%s\n", ip, p.self)
	if ip == p.self {
		return nil, errors.New("local")
	}
	getter := &PpcGetter{
		server: ip,
	}
	return getter, nil
}
// rpc handler
func (p *RpcPool) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	groupName := in.GetGroup()
	key := in.GetKey()
	group := GetGroup(groupName)
	rep := &pb.GetReply{
		Code: 0,
		Message: "success",
	}
	if group == nil {
		msg := fmt.Sprintf("no such group:%s", groupName)
		rep.Code = -1
		rep.Message = msg
		return rep, errors.New(msg)
	}

	view, err := group.Get(key)
	//fmt.Println(string(view))
	if err != nil {
		msg := fmt.Sprintf("group get err %v", err)
		rep.Code = -3
		rep.Message = msg
		return rep, errors.New(msg)
	}
	rep.Value = view
	return rep, nil
}
// rpc handler
func (p *RpcPool) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	name := in.GetGroup()
	key := in.GetKey()
	value := in.GetValue()
	group := GetGroup(name)
	rep := &pb.SetReply{
		Code: 0,
		Message: "success",
	}
	if group == nil {
		msg := fmt.Sprintf("no such group:%s", name)
		rep.Code = -1
		rep.Message = msg
		return rep, errors.New(msg)
	}

	//fmt.Printf("value is %s\n", string(value))
	
	err := group.Set(key, value)
	if err != nil {
		msg := fmt.Sprintf("group set err %v", err)
		rep.Code = -2
		rep.Message = msg
		return rep, errors.New(msg)
	}
	return rep, nil
}

func (p *RpcPool)Run() error{
	// rpc server
	lis, err := net.Listen("tcp", p.self)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCacherServer(s, p)
	fmt.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

type PpcGetter struct{
	server string
}

func NewPpcGetter(server string) *PpcGetter{
	return &PpcGetter{
		server: server,
	}
}

func (g *PpcGetter)Get(group, key string)([]byte, error){
	// rpc client
	// 拼request，请求service
	conn, err := grpc.Dial(g.server, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCacherClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.GetRequest{
		Group: group,
		Key: key,
	}
	rep, err := c.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("could not get value:  %v", err)
	}
	return rep.Value, nil
}

func (g *PpcGetter)Set(group, key string, value []byte)(error){
	// rpc client
	// 拼request，请求service
	conn, err := grpc.Dial(g.server, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCacherClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.SetRequest{
		Group: group,
		Key: key,
		Value: value,
	}
	_, err = c.Set(ctx, req)
	if err != nil {
		return fmt.Errorf("could not get value:  %v", err)
	}
	return nil
}