//一致性hash还需要一直ping机器，如果机器down掉，要remove ---todo
// 还需要去check配置中的ips，如果有增删机器，都要add和remove ---改成配置的方式，有一个常驻的timer
// remove注意要做数据的迁移---机器down掉了，内存数据是无法或得到，除非持久化，这个时候pick的机器上没有这个数据了，要么内存持久化恢复，要么从db恢复

// 2021.07.10 todo
// 1.sinker 可以改成写文件再试试看
// 2.pickpeer 的hash函数改为轮询，然后peergetter 改成grpc请求看看
// 3.有一个疑问，这个写到本地的时候是如何中转的，比如其是起了一个api的接口，然后这个api的接口会转成/group/name
// 但是这个group和name并没有 写入。peers一直有呀
// 还是需要测试看看


package gee
import (
	"go-cache/lru"
	"sync"
	"errors"
)

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

type Group struct {
	name string
	storage Sinker
	// local cache
	localCache *lru.Cache
	maxCacheSize int
	// pick peer ---including hash and get
	peers PeerPicker
}

func NewGroup(name string, cacheSize int, sinker Sinker)*Group {
	g := &Group{
		name: name,
		maxCacheSize: cacheSize,
		localCache: lru.New(cacheSize),
		storage: sinker,
	}
	mu.Lock()
	defer mu.Unlock()
	groups[name] = g // 如果重名怎么办
	return g
}

func GetGroup(name string) (g *Group){
	mu.Lock()
	defer mu.Unlock()
	g, ok := groups[name]
	if !ok {
		return nil
	}
	return g
}

func (g *Group)RegisterPeers(peers PeerPicker) {
	g.peers = peers
}

func (g *Group)Get(key string) ([]byte, error){
	// 先在localcache上查，然后再getFromPeer
	value, err := g.getFromCache(key)
	if err == nil {
		return value, err
	}
	value, err = g.getFromPeer(key)
	if err == nil {
		return value, err
	}
	return g.storage.Get(key)
}


// todo:最好是返回byte
func (g *Group)getFromCache(key string) ([]byte, error){
	value, err := g.localCache.Get(key)
	if err != nil {
		return nil, err
	}
	return value.([]byte), nil
}

func (g *Group)getFromPeer(key string) ([]byte, error){
	if g.peers == nil {
		return nil, errors.New("peers empty")
	}
	peer, err := g.peers.PickPeer(key) // key是否需要全局唯一，还是在某个group中唯一即可
	if err != nil {
		return nil, err
	}
	value, err := peer.Get(g.name, key)
	return value, err
}

func (g *Group)Set(key string, value []byte) (error){
	// 先在localcache上查，然后再getFromPeer
	err := g.setFromCache(key, value)
	if err != nil {
		return err
	}
	err = g.setFromPeer(key, value)
	if err != nil {
		return err
	}
	return g.storage.Set(key, value)
}

func (g *Group)setFromCache(key string, value []byte) (error){
	err := g.localCache.Set(key, value)
	return err
}

func (g *Group)setFromPeer(key string, value []byte) (error){
	if g.peers == nil {
		return errors.New("peers empty")
	}
	peer, err := g.peers.PickPeer(key)
	if err != nil {
		return err
	}
	err = peer.Set(g.name, key, value)
	return err
}





