package main
import(
	"fmt"
	"go-web/gee"
)

func main() {
	g := gee.New()
	g.Use(gee.Logger())
	v1 := g.Group("/v1", func(c *gee.Context) {
		fmt.Fprintf(c.Writer, "start URL.Path = %q\n", c.Req.URL.Path)
	})
	v2 := v1.Group("/v2", func(c *gee.Context) {
		fmt.Fprintf(c.Writer, "second URL.Path = %q\n", c.Req.URL.Path)
	})

	v1.Get("/", func(c *gee.Context) {
		fmt.Fprintf(c.Writer, "URL.Path = %q\n", c.Req.URL.Path)
	})

	v2.Post("/hello", func(c *gee.Context){
		word := c.Req.FormValue("word")
		c.Writer.Write([]byte(word))
	})

	g.Get("/v1/v2/group", func(c *gee.Context) {
		fmt.Fprintf(c.Writer, "group URL.Path = %q\n", c.Req.URL.Path)
	})

	g.Run(":8090")
	fmt.Println("hello,wold")
}