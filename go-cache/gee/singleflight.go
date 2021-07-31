package gee

import "sync"

type SingleFlight struct {
	Wg  sync.WaitGroup
	Mu  sync.Mutex
	Map map[string]bool // 这个和mu搭配使用是为了锁定key
}

func NewSingleFlight() *SingleFlight {
	return &SingleFlight{
		Map: make(map[string]bool),
		Wg:  sync.WaitGroup{},
	}
}
