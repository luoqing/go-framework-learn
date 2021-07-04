package gee
import(
	"strings"
	"fmt"
	"errors"
)

// todo, tree.handlers, remove groups
// todo, *, isWild， 这个正则匹配，在search的时候也比较麻烦parsePattern
// todo, : 后面是带参数的params，这个比较麻烦

type Router struct {
	methodTrees map[string]*node
}

func NewRouter() *Router {
	return &Router{
		methodTrees: make(map[string]*node),
	}
}

// 支持/xxx/xxx, 支持/xx/:id 支持/xxx/*
func (r *Router)parsePattern(pattern string) (parts []string) {
	// todo考虑*
	if len(pattern) == 0 {
		return
	}
	parts = strings.Split(pattern[1:], "/")
	for i, part := range parts {
		parts[i] =  "/" + part
	}
	return 
}

// 往这里面添加
// todo：使用handlers，这样就不用去比对group的prefix
func (r *Router) addRoute(method, pattern string, handler HandlerFunc) {
	tree, ok := r.methodTrees[method]
	if !ok {
		r.methodTrees[method] = &node{
			part: "",
		}
		tree, _ = r.methodTrees[method]; 
	}
	parts := r.parsePattern(pattern)
	n := tree.insert(parts, 0)
	if n == nil {
		fmt.Println("insert error")
	}
	n.handler = handler
	n.pattern = pattern
	// todo 打印出来traverse
}

func (r *Router) getRoute(method, path string) (handler HandlerFunc, err error){
	tree, ok := r.methodTrees[method]
	if !ok {
		r.methodTrees[method] = &node{
			part: "",
		}
		tree, _ = r.methodTrees[method]; 
	}
	parts := r.parsePattern(path)
	n := tree.search(parts, 0) // 这一行有问题
	if n == nil {
		return nil, errors.New("NOT found")
	} else {
		return n.handler, nil
	}
}

func (r *Router)handle(c *Context) {
	if handler, err := r.getRoute(c.Method, c.Path); err == nil {
		c.handlers = append(c.handlers, handler)
		c.Next()
	} else {
		fmt.Println(err)
		// output: 404 not found
		output := "404 not found"
		c.Writer.Write([]byte(output))
	}
}