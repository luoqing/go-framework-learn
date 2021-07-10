package consistenthash


import (
	"testing"
	"fmt"
)

func TestConsistentHash(t *testing.T) {
	c := New(5, nil)

	for i := 0; i < 10; i++ {
		si := fmt.Sprintf("%d", i)
		c.Add("172.18.1." + si)
	}

	for k, v := range c.hashMap {
		fmt.Println("Hash:", k, " IP:", v)
	}

	ipCount := make(map[string]int, 0)
	// 数据不均衡，数据有倾斜
	// 一致性hash有没有更好的hash算法，不让数据倾斜
	// todo：md5也算一种hash
	for i := 0; i < 1000; i++ {
		si := fmt.Sprintf("key%d", i)
		ip, err := c.Get(si)
		if err != nil {
			fmt.Println(err)
			continue
		}
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