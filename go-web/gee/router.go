package gee

type Router struct {
	route map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		route: make(map[string]HandlerFunc),
	}
}

// 往这里面添加
func (r *Router) addRoute(method, pattern string, handler HandlerFunc) {
	key := method + "_" + pattern
	r.route[key] = handler
}

func (r *Router)handle(c *Context) {
	method := c.Req.Method
	pattern := c.Req.URL.Path
	key := method + "_" + pattern
	if handler, ok := r.route[key]; ok {
		c.handlers = append(c.handlers, handler)
		c.Next()
	} else {
		// output: 404 not found
		output := "404 not found"
		c.Writer.Write([]byte(output))
	}
}