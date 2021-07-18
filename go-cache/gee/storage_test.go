package gee
import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
)
// https://geektutu.com/post/quick-go-test.html
func handleError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("failed", err)
	}
}

func TestRunServers(t *testing.T) {
	servers := []string{"localhost:8089", "localhost:8090", "localhost:8091", "localhost:8092", "localhost:8088"}
	/*
	addr := ":8091"
	Run(addr, servers)
	*/
	
	for i, addr := range servers{
		name := fmt.Sprintf("server %d", i)
		addr := addr
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			Run(addr, servers) // 如何并行启用多个服务并且阻塞在那里
		})
		
	}
	/* go Run(addr, servers)
	for {

	}*/
}

func TestConn(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	handleError(t, err)
	defer ln.Close()

	http.HandleFunc("/hello", helloHandler)
	go http.Serve(ln, nil)

	resp, err := http.Get("http://" + ln.Addr().String() + "/hello")
	handleError(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handleError(t, err)

	if string(body) != "hello world" {
		t.Fatal("expected hello world, but got", string(body))
	}
}

func TestHttp(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8089/", nil)
	w := httptest.NewRecorder()
	helloHandler(w, req)
	bytes, _ := ioutil.ReadAll(w.Result().Body)

	if string(bytes) != "hello world" {
		t.Fatal("expected hello world, but got", string(bytes))
	}
}
