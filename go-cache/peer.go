package gee

// 这个是分布式机器hash算法---不一定使用一致性hash
type  PeerPicker interface {
	PickPeer(key string) (PeerGetter, error)

}


type  PeerGetter interface{
	Get(group, key string) ([]byte, error)
	Set(group, key string, value []byte) (error)
}
