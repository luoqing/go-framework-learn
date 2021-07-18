package gee

import(
	"testing"
	"fmt"
	"flag"
)

// mock测试
// https://geektutu.com/post/quick-gomock.html
// 

func startCacheServer(addr string, svrs []string, sinker Sinker) {
	pool := NewHttpPool(addr, 5, svrs)
	
	name := "testg"
	g := NewGroup(name, 2<<10, sinker) // sinker 本地内存， name是内存的名称
	g.RegisterPeers(pool)
	pool.Run()

}
var addr, sinker string
func init() {
	testing.Init()
	flag.StringVar(&addr, "server", "localhost:8090", "cache server")
	flag.StringVar(&sinker, "sinker", "map", "cache sinker")
	flag.Parse()
}
// set curl -i -X POST "http://localhost:8092/api/set/testg/key" -d "value"
// get curl "http://localhost:8088/api/get/testg/key"
func TestGroupCache(t *testing.T){

	servers := []string{"localhost:8089", "localhost:8090", "localhost:8091", "localhost:8092", "localhost:8088"}
	fmt.Printf("start addr:%s\n", addr)

	var getter Sinker
	if sinker == "map" {
		getter = NewSMap()
	} else if sinker == "file" {
		getter = NewFileSinker("./")
	} else {
		t.Fatalf("error sinker:%s! sinker should be map or file!", sinker)
	}
	startCacheServer(addr, servers, getter)
	/*
	for i, addr := range servers{
		name := fmt.Sprintf("server %d", i)
		addr := addr
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			startCacheServer(addr, servers) // 这样测试是不行的，因为有一个groups的全局的变量
		})
		
	}*/

}

