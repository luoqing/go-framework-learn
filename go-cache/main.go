package main

import(
	"fmt"
	"flag"
	"net/http"
	"log"
	"go-cache/gee"
	"os"
)

// mock测试
// https://geektutu.com/post/quick-gomock.html
// 

func startCacheRpcServer(addr string, svrs []string, sinker gee.Sinker) {
	pool := gee.NewRpcPool(addr, 5, svrs)
	
	name := "testg"
	g := gee.NewGroup(name, 2<<10, sinker) // sinker 本地内存， name是内存的名称
	g.RegisterPeers(pool)
	pool.Run()

}

func startApiServer(apiAddr string, server string) {
	http.Handle("/api/get", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			group := r.URL.Query().Get("group")
			key := r.URL.Query().Get("key")
			client := gee.NewPpcGetter(server)
			value, err := client.Get(group, key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(value)
		}))

	http.Handle("/api/set", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			value := r.URL.Query().Get("value")
			group := r.URL.Query().Get("group")
			key := r.URL.Query().Get("key")
			client := gee.NewPpcGetter(server)
			b := []byte(value)
			err := client.Set(group, key, b)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			fmt.Fprintln(w, "success")
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr, nil))

}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	var addr, sinker, apiAddr string
	flag.StringVar(&addr, "rpc_server", "localhost:8090", "cache server")
	flag.StringVar(&sinker, "rpc_sinker", "map", "cache sinker")
	flag.StringVar(&apiAddr, "api_server", "localhost:8190", "cache api server")

	flag.Parse()

	servers := []string{"localhost:8089", "localhost:8090", "localhost:8091", "localhost:8092", "localhost:8088"}
	fmt.Printf("start addr:%s\n", addr)
	dstDir := "./" + addr
	exist, err := PathExists(dstDir)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}

	if !exist {
		fmt.Printf("no dir![%v]\n", dstDir)
		// 创建文件夹
		err := os.Mkdir(dstDir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
			return
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}

	var getter gee.Sinker
	if sinker == "map" {
		getter = gee.NewSMap()
	} else if sinker == "file" {
		getter = gee.NewFileSinker(dstDir)
	} else {
		log.Fatalf("error sinker:%s! sinker should be map or file!", sinker)
	}
	go startApiServer(apiAddr, addr)
	startCacheRpcServer(addr, servers, getter)
	
}
