package main

import (
	//	"encoding/json"
	//	"geerpc/codec"
	"encoding/json"
	"fmt"
	"geerpc/codec"
	"geerpc/server"
	"log"
	"net"
	"time"
)

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	server.Accept(l)
}

// day1 with codec--god or json
func main() {
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple server client
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	// send options
	//server.DefaultOption.CodecType = codec.GobType
	server.DefaultOption.CodecType = codec.JsonType
	_ = json.NewEncoder(conn).Encode(server.DefaultOption)
	//cc := codec.NewGobCodec(conn)
	cc := codec.NewJsonCodec(conn)

	// send request & receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		err := cc.Write(h, fmt.Sprintf("server req %d", h.Seq))
		if err != nil {
			fmt.Println(err)
			return
		}
		err = cc.ReadHeader(h)
		if err != nil {
			fmt.Println(err)
			return
		}
		var reply string
		_ = cc.ReadBody(&reply)
		if err != nil {
			fmt.Println(err)
			return
		}
		log.Println("reply:", reply)
	}
}

/* day2 with client
// 可以指定option，使用不同的codec
// 封装了read和write，不用考虑读取不完整，或者写入部分
func main() {
    log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)
    //client, _ := geerpc.Dial("tcp", <-addr)
	client, _ := Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}
*/
