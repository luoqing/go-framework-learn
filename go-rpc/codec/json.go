package codec

// 2021.09.19 todo:
// 1.先run起来，然后实现json的序列化的方式
// 2.server部分改造，trpc的server端如何处理各个连接
// 读取请求 readRequest
// 处理请求 handleRequest
// 回复请求 sendResponse
// 3.服务注册——反射的方法和使用生成自动化代码的方式，使用反射的方法又如何实现的呢？——实现路由
// client.invoke(connect, sendrequest)
// 4.高性能客户端是解决什么问题呢？——实现client
// 5.如何进行超时处理？
// 6.如何支持http协议，如果你来支持，你会如何支持？
import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

func init() {
	Register(JsonType, NewJsonCodec)
}

type JsonCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *json.Decoder
	enc  *json.Encoder
}

func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JsonCodec{
		conn: conn,
		buf:  buf,
		dec:  json.NewDecoder(conn),
		enc:  json.NewEncoder(buf),
	}
}

func (c *JsonCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *JsonCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *JsonCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec: json error encoding header:", err)
		return err
	}
	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: json error encoding body:", err)
		return err
	}
	return nil
}

func (c *JsonCodec) Close() error {
	return c.conn.Close()
}
