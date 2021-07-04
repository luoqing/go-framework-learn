package gee
import(
	"strings"
	"fmt"
	"errors"
	"net/url"
)

// todo, tree.handlers, remove groups---done
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
func (r *Router)parsePattern(pattern string) (parts []string, params []string, err error) {
	// todo考虑*
	if len(pattern) == 0 || pattern[0] != '/'{
		return parts, params, fmt.Errorf("invalid pattern:%s", pattern)
	}
	s := 0
	isWild := false
	end := false
	for i := 1; i < len(pattern); i ++ {
		if pattern[i] == ' '{
			return parts, params, fmt.Errorf("invalid pattern:%s", pattern)
		}
		if pattern[i] == '*'{
			part := pattern[s:i+1]
			parts = append(parts, part)
			end = true
			break
		}
		if pattern[i] == '/' {
			part := pattern[s:i]
			if part != "/" {
				parts = append(parts, part)
				if part[1] == ':' {
					// :后面必须要有数据， 获取keys
					if len(part) < 3 {
						return parts, params, fmt.Errorf("invalid pattern:%s", pattern)
					}
					params = append(params, part[2:])
				}
				// isWild后面必须都是带params
				if isWild && i < len(pattern)-1 && pattern[i+1] != ':'{
					return parts, params, fmt.Errorf("invalid pattern:%s", pattern)
				}
			}
			s = i
		} else if pattern[i] == ':' {
			// :前面必须是/
			if pattern[i-1] != '/' {
				return parts, params, fmt.Errorf("invalid pattern:%s", pattern)
			} else {
				// 获取params
				isWild = true
			}
		}
	}
	if !end {
		part := pattern[s:]
		if isWild && (len(part)< 3 || part[1] != ':'){
			return parts, params, fmt.Errorf("invalid pattern:%s", pattern)
		}
		if len(part) > 1 {
			parts = append(parts, part)
			if isWild {
				params = append(params, part[2:])
			}
		}
		
	}
	
	return
}

func (r *Router)parsePath(path string) (parts []string) {
	if path == "/" {
		parts = append(parts, path)
		return
	}
	tmpParts := strings.Split(path[1:], "/")
	for _, part := range tmpParts {
		if part == "" {
			continue
		}
		// 如果有:, * 就报错，找不到
		part =  "/" + part
		parts = append(parts, part)
	}
	return 
}

// 往这里面添加
// todo：使用handlers，这样就不用去比对group的prefix
func (r *Router) addRoute(method, pattern string, handlers ...HandlerFunc) error {
	fmt.Printf("addRoute method:%s, pattern:%s handlers cnt:%d\n", method, pattern, len(handlers))

	tree, ok := r.methodTrees[method]
	if !ok {
		r.methodTrees[method] = &node{
			part: "",
		}
		tree, _ = r.methodTrees[method]; 
	}
	parts, _, err := r.parsePattern(pattern)
	if err != nil {
		return err
	}
	n := tree.insert(parts, 0)
	if n == nil {
		return errors.New("insert error")
	}
	n.handlers = append(n.handlers, handlers...)
	n.pattern = pattern
	// todo 打印出来traverse
	return nil
}

func (r *Router) getRoute(method, path string) (params url.Values, handlers []HandlerFunc, err error){
	tree, ok := r.methodTrees[method]
	if !ok {
		r.methodTrees[method] = &node{
			part: "",
		}
		tree, _ = r.methodTrees[method]; 
	}
	parts := r.parsePath(path) // parsePath
	n, params := tree.search(parts, 0) // 这一行有问题
	if n == nil {
		return params, nil, errors.New("NOT found")
	} else {
		return params, n.handlers, nil
	}
}

func (r *Router)handle(c *Context) {
	if params, handlers, err := r.getRoute(c.Method, c.Path); err == nil {
		c.handlers = handlers
		c.Params = params
		c.Next()
	} else {
		fmt.Println(err)
		// output: 404 not found
		output := "404 not found"
		c.Writer.Write([]byte(output))
	}
}