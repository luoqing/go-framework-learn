package gee
import(
	"net/http"
	"strings"
//	"fmt"
)

type HandlerFunc func(ctx *Context)
type Engine struct {
	route *Router
	*RouterGroup // for global middlewares
	groups []*RouterGroup // for groups with prefix
}

// 如何将group的信息写入到engine呢
type RouterGroup struct {
	engine *Engine // all groups share a Engine instance
	Prefix string
	Middlewares []HandlerFunc
}


func (r *RouterGroup) Use(middlewares ...HandlerFunc) {
	r.Middlewares = append(r.Middlewares, middlewares...)
}

func (r *RouterGroup) Group(prefix string, middlewares ...HandlerFunc) *RouterGroup{
	grp := &RouterGroup{
		Prefix: r.Prefix + prefix,
		Middlewares: middlewares, // 此处不拼接，防止重复
		engine: r.engine,
	}
	r.engine.groups = append(r.engine.groups, grp)
	return grp;
}

// engine继承了routegroup，可以直接调用routegroup的函数，同时routegroup有一个engine，进行大家share
func (r *RouterGroup)addRoute(method, pattern string, handler HandlerFunc) {
	r.engine.route.addRoute(method, r.Prefix + pattern, handler)
}

func (r *RouterGroup)Get(pattern string, handler HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

func (r *RouterGroup)Post(pattern string, handler HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

// 这个地方有点奇怪，就是group和engine如何联动
func New() *Engine {
	g := &Engine{
		route:NewRouter(),
	}
	g.RouterGroup = &RouterGroup{engine: g}
	g.groups = []*RouterGroup{g.RouterGroup} // 这个将全局的也加进来了, 因为有的可能是没有
	return g
}
func (g *Engine)ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// get method and pattern
	c := NewContext(w, req)
	// 获取所有的handlers， 先获取path，然后判断path是否有所有groups的middlewares
	for _, grp := range g.groups {
		// 这个要去重，获取最接近的那个就行
		if strings.HasPrefix(c.Path, grp.Prefix) {
			c.handlers = append(c.handlers, grp.Middlewares...)
		}
	}
	g.route.handle(c)
}

func (g *Engine)Run(addr string) (err error){
	return http.ListenAndServe(addr, g)
}


