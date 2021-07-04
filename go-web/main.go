package main
import(
	"fmt"
	"go-web/gee"
)

func main() {
	g := gee.New()
	g.Use(gee.Logger()) // 1个

	v1 := g.Group("/v1", func(c *gee.Context) {
		fmt.Fprintf(c.Writer, "v1 URL.Path = %q\n", c.Req.URL.Path)
	})
	v2 := v1.Group("/v2", func(c *gee.Context) {
		fmt.Fprintf(c.Writer, "v2 URL.Path = %q\n", c.Req.URL.Path)
	})

	v1.Get("/", func(c *gee.Context) {
		fmt.Fprintf(c.Writer, "URL.Path = %q\n", c.Req.URL.Path)
	})

	v2.Post("/hello", func(c *gee.Context){
		word := c.Req.FormValue("word")
		c.Writer.Write([]byte(word))
	})

	v1.Get("/param/:word", func(c *gee.Context){
		word := c.Req.FormValue("word")
		c.Writer.Write([]byte(word))
	})

	g.Get("/any/*", func(c *gee.Context) {
		fmt.Fprintf(c.Writer, "any URL.Path = %q\n", c.Req.URL.Path)
	})

	g.Use(gee.Recovery()) // 能很好捕捉异常吗？

	g.Run(":8090")
	fmt.Println("hello,wold")
}