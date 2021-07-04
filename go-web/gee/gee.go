package gee
import(
	"net/http"
//	"path"
//	"fmt"
)

type HandlerFunc func(ctx *Context)
type Engine struct {
	route *Router
	*RouterGroup // for  middlewares
	//groups []*RouterGroup // for groups with prefix , with *RouterGroup for middlewares
}

// 如何将group的信息写入到engine呢
type RouterGroup struct {
	engine *Engine // all groups share a Engine instance
	prefix string
	middlewares []HandlerFunc
}


func (r *RouterGroup) Use(middlewares ...HandlerFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *RouterGroup) Group(prefix string, middlewares ...HandlerFunc) *RouterGroup{
	grp := &RouterGroup{
		prefix: r.prefix + prefix,
		middlewares: append(r.middlewares, middlewares...), // 有groups此处不拼接，防止重复
		engine: r.engine,
	}
	//r.engine.groups = append(r.engine.groups, grp)
	return grp;
}

// engine继承了routegroup，可以直接调用routegroup的函数，同时routegroup有一个engine，进行大家share
func (r *RouterGroup)addRoute(method, pattern string, handler HandlerFunc) error {
	handlers := append(r.middlewares, handler)
	return r.engine.route.addRoute(method, r.prefix + pattern, handlers...)
}

func (r *RouterGroup)handle(method, pattern string, handler HandlerFunc) error {
	// 在此处将middlewarres添加到router
	//handlers := r.combineHandlers([]HandlerFunc{handler})
	return r.addRoute(method, pattern, handler)
}


func (r *RouterGroup)Get(pattern string, handler HandlerFunc) {
	err := r.handle("GET", pattern, handler)
	if err != nil {
		panic(err)
	}
}

func (r *RouterGroup)Post(pattern string, handler HandlerFunc) {
	err := r.handle("POST", pattern, handler)
	if err != nil {
		panic(err)
	}
}

// 这个地方有点奇怪，就是group和engine如何联动
func New() *Engine {
	g := &Engine{
		route:NewRouter(),
	}
	g.RouterGroup = &RouterGroup{engine: g}
	//g.groups = []*RouterGroup{g.RouterGroup} // 这个将全局的也加进来了, 因为有的可能是没有
	return g
}
func (g *Engine)ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// get method and pattern
	c := NewContext(w, req)
	/*
	// 获取所有的handlers， 先获取path，然后判断path是否有所有groups的middlewares
	for _, grp := range g.groups {
		// 这个要去重，获取最接近的那个就行
		if strings.HasPrefix(c.Path, grp.Prefix) {
			c.handlers = append(c.handlers, grp.Middlewares...)
		}
	}*/
	g.route.handle(c)
}

func (g *Engine)Run(addr string) (err error){
	return http.ListenAndServe(addr, g)
}


/*
func (group *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(group.middlewares) + len(handlers)
	if finalSize >= 256 {
		panic("too many handlers")
	}
	mergedHandlers := make([]HandlerFunc, finalSize)
	copy(mergedHandlers, group.middlewares)
	copy(mergedHandlers[len(group.middlewares):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.prefix, relativePath)
}
func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) byte {
	l := len(str)
	
	return str[l-1]
}*/