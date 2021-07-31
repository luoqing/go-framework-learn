package gee

import "sync"

/*
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
*/
// 自己写的，和别人写的区别主要有以下
// 1.不是针对一个请求一个wg
// 2. make的部分没加锁，而是在new哪里做的
// 3. 就是防止缓存击穿，是支队sinker做了加锁，而不是group.Get整个函数---这点确实做的不够
// 缓存雪崩，所有key同时失效（可能因为设置相同失效时间），导致穿透
// 缓存击穿，某一个key失效时很多次请求，导致击穿
// 缓存穿透，某一个key很多次请求，但是该key在缓存不存在，

type call struct {
	val interface{}
	err error
	wg  sync.WaitGroup
}
type SingleFlight struct {
	mu sync.Mutex
	m  map[string]*call // 这个和mu搭配使用是为了锁定key
}

func NewSingleFlight() *SingleFlight {
	return &SingleFlight{}
}

func (s *SingleFlight) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	s.mu.Lock()
	if s.m == nil {
		s.m = make(map[string]*call) // 这个还不能在NewSingleFlight进行
	}
	if c, ok := s.m[key]; ok {
		s.mu.Unlock() // lock和unlock配套使用
		c.wg.Wait()   // wait等处理完，获取到数据, 此处不wait，wal拿到可能是nil
		return c.val, c.err
	}
	s.mu.Unlock() // 这里要用defer吗

	s.mu.Lock()
	c := new(call)
	c.wg.Add(1)
	s.m[key] = c
	s.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	s.mu.Lock()
	delete(s.m, key)
	s.mu.Unlock()
	return c.val, c.err
}
