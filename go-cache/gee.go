// http 请求 
// 先pick一台机器，然后将数据写入

// 先在本地cache查能否查到，如果不行，pick一台机器，进行数据查询
// 先支持http接口，数据写入和查询，写入一个

// 可以在这里定义路由 pick机器，如果是本机
// func (g *Engine)ServeHTTP(w http.ResponseWriter, req *http.Request) {
// http.ListenAndServe(addr, g)

package gee
import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"go-cache/lru"
	"go-cache/consistenthash"
	"net/url"
	"log"
	"io"
)

/*
type PeerPicker interface {
	func PickPeer(key string) (string error) {

	}
}*/

type Storage struct{
	data map[string]string
	localCache *lru.Cache
	maxCacheSize int
	peers []string
	replicas int
	peersHash *consistenthash.ConsistentHash
}

//一致性hash还需要一直ping机器，如果机器down掉，要remove ---todo
// 还需要去check配置中的ips，如果有增删机器，都要add和remove ---改成配置的方式，有一个常驻的timer
// remove注意要做数据的迁移---机器down掉了，内存数据是无法或得到，除非持久化，这个时候pick的机器上没有这个数据了，要么内存持久化恢复，要么从db恢复

// replicas int, hashFunc HashHandler
func Run(addr string, servers []string) {

	localCacheSize := 256
	replicas := 10
	s := NewStorage(localCacheSize, replicas, servers)
	/*
	//http.HandleFunc("/", helloHandler)
	http.HandleFunc("/get", s.getValue)
	http.HandleFunc("/set", s.setValue)
	http.HandleFunc("/api/get", s.get) // 分布式
	http.HandleFunc("/api/set", s.set) // 分布式
	*/
	fmt.Printf("start server：%s\n", addr)
	err := http.ListenAndServe(addr, s)
	if err != nil {
		fmt.Println("ListenAndServe error: ", err.Error())
	}
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}


func NewStorage(localCacheSize int, replicas int, ips []string) *Storage{
	s := &Storage{
		data: make(map[string]string),
		maxCacheSize: localCacheSize,
		replicas: replicas,
	}
	//var replicas int, hashFunc consistenthash.HashHandler;
	
	s.localCache = lru.New(localCacheSize)
	s.peers = ips
	s.peersHash = consistenthash.New(replicas, nil)
	s.peersHash.Add(s.peers...)
	return s
	
}

func (s *Storage)ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	switch path {
	case "/get":
		s.getValue(w, req)
	case "/set":
		s.setValue(w, req)
	case "/api/get":
		s.get(w, req)
	case "/api/set":
		s.set(w, req)
	}
}


func (s *Storage)save2LocalCache(key string, value string) (err error){
	err = s.localCache.Set(key, value)
	return
}

func (s *Storage)getFromLocalCache(key string) (value string, err error){
	v, err := s.localCache.Get(key)
	if err != nil {
		return
	}
	return v.(string), nil
}

func (s *Storage)pickPeer(key string) (ip string, err error){
	ip, err = s.peersHash.Get(key)
	return
}

func ForwardHandler(proxyUrl string, writer http.ResponseWriter, request *http.Request) {
    u, err := url.Parse(proxyUrl)
    if nil != err {
        log.Println(err)
        return
    }

    proxy := httputil.ReverseProxy{
        Director: func(request *http.Request) {
            request.URL = u
        },
    }

    proxy.ServeHTTP(writer, request)
}

// 先在本地查询
// 本地查询不到再pick一台机器，请求这台机器去拿数据
func (s *Storage)get(w http.ResponseWriter, req *http.Request) {
	key := req.FormValue("k")
	if key == ""{
		fmt.Fprintln(w, "params error")
		return
	}
	value, err := s.getFromLocalCache(key)
	if err == nil {
		fmt.Fprintln(w, value)
		return
	}
	ip, err := s.pickPeer(key)
	if err != nil {
		return
	}
	// 拼url，http.DoRequest
	url := fmt.Sprintf("http://%s/get?k=%s", ip, key) // urlencode
	resp, err := http.Get(url)
	if err != nil {
		errmsg := fmt.Sprintf("set peer failed:%v", err)
		fmt.Fprintln(w, errmsg)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Fprintln(w, string(body))
	//ForwardHandler(url, w, req)
	
}

func (s *Storage)set(w http.ResponseWriter, req *http.Request) {
	key := req.FormValue("k")
	value := req.FormValue("v")
	if key == "" || value == ""{
		fmt.Fprintln(w, "params error")
		return
	}
	err := s.save2LocalCache(key, value)
	if err != nil {
		// 本地local和远程如何保持一致性
		fmt.Fprintln(w, "set local failed")
		return
	}
	ip, err := s.pickPeer(key)
	if err != nil {
		fmt.Fprintln(w, "pick peer failed")
		return
	}
	fmt.Printf("key:%s pick peer:%s\n", key, ip)
	url := fmt.Sprintf("http://%s/set?k=%s&v=%s", ip, key, value)
	resp, err := http.Get(url)
	if err != nil {
		errmsg := fmt.Sprintf("set peer failed:%v", err)
		fmt.Fprintln(w, errmsg)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Fprintln(w, string(body))
	//ForwardHandler(url, w, req)
	// 其实可以使用roundtrip转发会更好一点
}

func (s *Storage)getValue(w http.ResponseWriter, req *http.Request) {
	key := req.FormValue("k")
	if key == ""{
		fmt.Fprintln(w, "params error")
		return
	}
	v, ok := s.data[key]; 
	if !ok {
		fmt.Fprintln(w, "key not existed")
		return
	}
	value := fmt.Sprintf("hello,%s", v)
	fmt.Fprintln(w, value)
}

func (s *Storage)setValue(w http.ResponseWriter, req *http.Request) {
	key := req.FormValue("k")
	value := req.FormValue("v")
	if key == "" || value == ""{
		fmt.Fprintln(w, "params error")
		return
	}
	s.data[key] = value
	fmt.Fprintln(w, "set success")
}