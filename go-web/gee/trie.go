package gee

import (
	"fmt"
	//"strings"
	"container/list"
)

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
	//handlers []HandlerFunc
	handler HandlerFunc
}

//https://github.com/geektutu/7days-golang/blob/master/gee-web/day3-router/gee/trie.go
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 找到从 index 到 len 的path符合的node
// 这个search肯定有问题
func (n *node)search(parts []string, index int) (*node) {
	if (index == len(parts)) {
		return n
	}

	// match children
	child := n.matchChildren(parts[index])
	if child != nil {
		fmt.Printf("index:%d, part:%s, node：%v\n", index, parts[index], child.part)
		return child.search(parts, index + 1) // 这个有问题， 为啥会请求 /favicon.ico？
	}
	fmt.Printf("error index:%d, part:%s, node：%v\n", index, parts[index], child)
	return nil
}

// 先找到合适的位置，先找到parent
// 然后parent的children
// 返回最后那个节点，要设置handlers和isWild， 还要设置pattern
func (n *node)insert(parts []string, index int)(*node) {
	if index == len(parts){
		return n
	}
	child := n.matchChildren(parts[index])
	if child == nil {
		child = &node{
			part: parts[index],
		}
		n.children = append(n.children, child)
	}
	return child.insert(parts, index+1)
}

func (n *node)matchChildren(part string) (*node){
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}

	return nil
}

// 把traverse的数据都写在list中
// 也可以使用slice
// 广度优先搜索
func (n *node)traverse(l *list.List) {
	l.PushBack(n)
	for _, child := range n.children {
		child.traverse(l)
	}
}

// insert
// search
// addroute
// getroute