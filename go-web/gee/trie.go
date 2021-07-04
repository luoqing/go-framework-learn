package gee

import (
	"fmt"
	//"strings"
	"container/list"
	"net/url"
)

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
	handlers []HandlerFunc
	//handler HandlerFunc
}

//https://github.com/geektutu/7days-golang/blob/master/gee-web/day3-router/gee/trie.go
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 找到从 index 到 len 的path符合的node
// 这个search肯定有问题
func (n *node)search(parts []string, index int) (*node, url.Values) {
	params := make(url.Values)
	if (index == len(parts) || (len(n.part) > 1 && n.part[1] == '*')) {
		return n, params
	}

	// match children
	children := n.matchChildren(parts[index])
	for _, child := range children {
		fmt.Printf("index:%d, part:%s, node：%v\n", index, parts[index], child.part)
		if child.isWild {
			key := child.part[2:]
			value := parts[index][1:]
			params[key] = []string{value}
		}
		return child.search(parts, index + 1) // 这个有问题， 为啥会请求 /favicon.ico？
	}
	fmt.Printf("error index:%d, part:%s, node：%v\n", index, parts[index], n)
	return nil, nil
}

// 先找到合适的位置，先找到parent
// 然后parent的children
// 返回最后那个节点，要设置handlers和isWild， 还要设置pattern
func (n *node)insert(parts []string, index int)(*node) {
	if index == len(parts){
		return n
	}
	child := n.matchChild(parts[index])

	if child == nil {
		isWild := false
		part := parts[index]
		if len(part) > 1 && part[1] == ':' {
			isWild = true
		}
		child = &node{
			part: part,
			isWild: isWild,
		}
		n.children = append(n.children, child)
	}
	fmt.Printf("insert index:%d, part:%s, node：%v\n", index, parts[index], child)

	return child.insert(parts, index+1)
}

func (n *node)matchChild(part string) (*node){
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}

	return nil
}

func (n *node)matchChildren(part string) ([]*node){
	var children []*node
	for _, child := range n.children {
		if child.part == part || child.isWild || (len(child.part) > 1 && child.part[1] == '*'){
			children = append(children, child)
		}
	}

	return children
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