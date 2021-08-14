package codec

import (
	"io"
	"sync"
)

type Header struct {
	ServiceMethod string // format "Service.Method"
	Seq           uint64 // sequence number chosen by client
	Error         string
}

// codec 也可以自定义协议的encode和decode，frameBuilder就是如何读取，serialization序列化：pb,json,xml等
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

const (
	GobType  string = "application/gob"
	JsonType string = "application/json" // not implemented
)

var (
	NewCodecFuncMap map[string]NewCodecFunc
	lock            sync.RWMutex
)

func init() {
	NewCodecFuncMap = make(map[string]NewCodecFunc)
	
}

func Register(name string, f NewCodecFunc) {
	lock.Lock()
	NewCodecFuncMap[name] = f
	lock.Unlock()
}
