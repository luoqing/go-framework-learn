package consistenthash

import (
	"sync"
	"fmt"
	"hash/crc32"
	"sort"
	"errors"
	"strings"
)

// 问题
// how to solve hash conflict hash冲突
// concurrencent safe 并发安全
// remove elem from slice 
// https://github.com/golang/groupcache/blob/master/consistenthash/consistenthash.go
/*
type struct Data {
	key string
	value interface{}
}

type struct Node struct {
	data map[string]Data
	index int
	ip string // port, hostname
	weight double // 按权重初始化进行
}*/
// IsEmpty() bool
// 研究下有哪些hash 函数，一致性hash如何保证不碰撞

type HashHandler func(data []byte) uint32
type ConsistentHash struct {
	replicas int
	hashFunc HashHandler
	//hashMap map[int]*Node
	hashMap map[int]string // 每个虚拟节点的hash对应哪台ip， 这个Node 可以简化为string
	hashKeys []int // 每个虚拟节点的hash值，这个hash值应该是int
	sync.RWMutex
}


//https://studygolang.com/articles/4017 --- 考虑了Node， 也考虑了并发安全

func New(replicas int, hashFunc HashHandler) *ConsistentHash{
	if hashFunc == nil {
		hashFunc = crc32.ChecksumIEEE
	}

	return &ConsistentHash{
		replicas: replicas,
		hashFunc: hashFunc,
		hashMap: make(map[int]string),
	}
	
}

func (c *ConsistentHash)Add(ips... string) {
	c.Lock()
	defer c.Unlock()
	//c.nodes = append(c.nodes, nodes...)
	for _, ip := range ips {
		// node + 
		i := 0
		for i < c.replicas {
			key := fmt.Sprintf("%s_%d", ip, i)
			hash := int(c.hashFunc([]byte(key))) // hash会不会碰撞？
			/* 这个碰撞了不行，所以这里没有考虑，尽量使用不会conflict的
			for n, ok := c.hashMap[hash]; ok && n != ip{
				// key 
				// 再加一个随机数，重新计算hash
				hash = c.hashFunc[key]
			}*/
			c.hashKeys = append(c.hashKeys, hash)
			c.hashMap[hash] = key
			i ++
		}
	}
	sort.Slice(c.hashKeys, func(i, j int) bool {
		return c.hashKeys[i] < c.hashKeys[j]
	})
	fmt.Println(c.hashKeys)

}


func (c *ConsistentHash)Remove(ip string) {
	c.Lock()
	defer c.Unlock()
	/* 为啥要去掉node的逻辑，因为这个耦合的不好
	for _, n := range c.nodes { 
		if n.ip == ip {
			// del
			c.nodes  = append(c.nodes[:i], c.nodes[i+1:]...)
		}
	}*/
	
	i := 0
	for i < c.replicas {
		key := fmt.Sprintf("%s_%d", ip, i)
		hash := int(c.hashFunc([]byte(key)))
		/* 这个碰撞了不行，所以这里没有考虑，尽量使用不会conflict的
		for n, ok := c.hashMap[hash]; !ok || n != ip {
			// 重新计算hash，直至满足上述条件
		}*/
		delete(c.hashMap, hash)
		j := sort.Search(len(c.hashKeys), func(j int)bool{
			return c.hashKeys[j] >= hash
		})
		c.hashKeys = append(c.hashKeys[:j], c.hashKeys[j+1:]...)
	}
	/*
	// node上的数据全部迁移到新的位置
	for k, v := range node.data {
		newNode, _ := c.Hash(k)
		newNode.data[k] = v;
	}*/

}

func (c *ConsistentHash)Get(key string) (string, error){
	c.RLock()
	defer c.RUnlock()
	hash := int(c.hashFunc([]byte(key))) // 这些都不考虑碰撞，这个就算碰撞还好
	index := sort.Search(len(c.hashKeys), func(i int)bool{
		return c.hashKeys[i] >= hash
	})
	hashKey := c.hashKeys[index % len(c.hashKeys)]
	if node, ok := c.hashMap[hashKey]; ok {
		tmp := strings.Split(node, "_")
		ip := tmp[0]
		return ip, nil
	} else {
		return "", errors.New("ERROR HASH MAP")
	}
}


// 指定hash函数，初始化多少个虚拟节点
// new(replicas, hashfunc)
// add(node)---计算每个node对应的虚拟节点对应第几个node， 每个node的虚拟节点的hash值
// remove(node)---将node上的所有key重新计算hash， 将node对应的hashkeys进行移除，将node也进行移除，
// get(key) --- 计算key的hash值， 然后二分查找离哪个hashKey最近，然后找到其对应的node的index，返回node的信息

/*
func main() {

	c := New()

	for i := 0; i < 10; i++ {
		si := fmt.Sprintf("%d", i)
		c.Add(NewNode("172.18.1."+si)
	}

	for k, v := range c.hashMap {
		fmt.Println("Hash:", k, " IP:", v)
	}

	ipCount := make(map[string]int, 0)
	for i := 0; i < 1000; i++ {
		si := fmt.Sprintf("key%d", i)
		ip := c.Get(si)
		if _, ok := ipCount[ip]; ok {
			ipCount[ip] += 1
		} else {
			ipCount[ip] = 1
		}
	}

	for k, v := range ipCount {
		fmt.Println("Node IP:", k, " count:", v)
	}

}
*/