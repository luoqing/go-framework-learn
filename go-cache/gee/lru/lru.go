package lru
import (
	"container/list"
	"fmt"
	"errors"
)

// 并发安全
// https://github.com/golang/groupcache/blob/master/lru/lru.go

// least recently use
// cache
// entry, key 也是interface{}
// Cache
type Key interface{}

type LRU struct {
	size int
	len int
	htable map[Key]*list.Element
	queue *list.List
}
type Entry struct {
	key Key
	value interface{}
}

func New(maxSize int) *LRU{
	return &LRU{
		size:maxSize,
		len:0,
		htable:make(map[Key]*list.Element),
		queue:list.New(),
	}
}

func (r *LRU)Get(key string) (interface{}, error) {
	if elem, ok := r.htable[key]; ok {
		r.queue.MoveToFront(elem)
		fmt.Printf("key:%s value:%v\n", key, elem.Value)
		return elem.Value.(*Entry).value, nil
	}
	fmt.Printf("key:%s NOT found!\n", key)
	return nil, errors.New("NOT found!")
}

func (r *LRU)Set(key Key, v interface{}) error{
	n := &Entry{
		key: key,
		value: v,
	}
	if elem, ok := r.htable[key]; ok {
		r.queue.MoveToFront(elem)
		// how to update value
		elem.Value = n
	} else if (r.len < r.size) {
		elem = r.queue.PushFront(n)
		r.htable[key] = elem
		r.len ++
	} else {
		// removeOldest
		r.RemoveOldest()
		elem = r.queue.PushFront(n)
		r.htable[key] = elem
	}
	return nil
}

func (r *LRU)RemoveOldest() {
	elem := r.queue.Back()
	n := elem.Value.(*Entry)
	
	delete(r.htable, n.key) // how to get key
	r.queue.Remove(elem)
	// hashtable
}

func (r *LRU)Len() int{
	return r.len
}

// 重置就行了
func (r *LRU)Clear(){
}

