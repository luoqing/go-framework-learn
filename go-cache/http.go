package gee
import (
	"go-cache/consistenthash"
	"net/http"
	"net/url"
	"errors"
	"fmt"
	"io/ioutil"
	"bytes"
	"strings"
	"io"
)

type HttpPool struct {
	self string // 分布式和本地保持统一
	basePath string // base path
	replicas int 
	peersHash *consistenthash.ConsistentHash
	//map[string]*HttpGetter // 为了映射，其实这个可以不要的， 一个ip对应一个连接可以
	//mu *sync.Mutex
}

type HttpGetter struct{
	server string
	basePath string
}

func NewHttpPool(self string, replicas int, ips []string) *HttpPool{
	p := &HttpPool{
		self: self,
		replicas: replicas,
	}
	p.peersHash = consistenthash.New(replicas, nil)
	p.peersHash.Add(ips...)
	return p
}

/*
func (p *HttpPool)Set(peers...string) {
	p.peersHash.Add(peers...)
}*/

func (p *HttpPool)PickPeer(key string) (PeerGetter, error) {
	// 选择出一台服务器
	// 然后拼好server和path，然后返回HttpGetter
	ip, err := p.peersHash.Get(key)
	if err != nil {
		return nil, err
	}
	if ip == p.self {
		return nil, errors.New("local")
	}
	getter := &HttpGetter{
		server: ip,
		basePath: "/api",
	}
	return getter, nil

}

func (g *HttpGetter)Get(group, key string)([]byte, error){
	// 根据server和path，拼出url，然后请求获取到数据
	baseURL := fmt.Sprintf("%s%s", g.server, g.basePath)
	u := fmt.Sprintf(
		"%v/get/%v/%v",
		baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

func (g *HttpGetter)Set(group, key string, value []byte)(error){
	// 根据server和path，拼出url，然后请求获取到数据
	baseURL := fmt.Sprintf("%s%s", g.server, g.basePath)
	u := fmt.Sprintf(
		"%v/set/%v/%v",
		baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	contentType := "application/octet-stream"
	r := bytes.NewReader(value)
	res, err := http.Post(u, contentType, r)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	return nil
}

func (p* HttpPool)apiGet(w http.ResponseWriter, r *http.Request) {
	basePath := "/api/get"
	parts := strings.SplitN(r.URL.Path[len(basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view)
}


func (p* HttpPool)apiSet(w http.ResponseWriter, r *http.Request) {
	basePath := "/api/set"
	parts := strings.SplitN(r.URL.Path[len(basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	defer r.Body.Close()
	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "value body error ", http.StatusInternalServerError)
		return
	}
	
	err = group.Set(key, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "success")
	//w.Write("success")
}


// 之前给一个服务器启动了一个storage
// 这个给每个服务器启动一个storage吗，
// 给每个服务器启动一个groupcache，一个groupcache对应自己的
func (p *HttpPool)ServeHTTP(w http.ResponseWriter, req *http.Request) {
	/*
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}*/
	// 两个路由本地get和set的
	// /api/set
	// /api/get
	path := req.URL.Path
	
	switch path {
	case "/api/get":
		p.apiGet(w, req)
	case "/api/set":
		p.apiSet(w, req)
	}

	// 根据不同路由获取不同的数据
	// 查询的逻辑---分布式查询和本地查询能统一吗
	// 获取group，key参数
	// 根据group.get(key)
	
	// 写入的逻辑---
	// 获取group，key，value参数
	// 根据group.set(key，value)
	// 注意自己

}

