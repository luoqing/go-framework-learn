package gee

import (
	"runtime"
	"net/http"
	"log"
	"fmt"
	"strings"
)

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller, get call stack

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc) // function
		file, line := fn.FileLine(pc) // file line
		str.WriteString(fmt.Sprintf("\n\t%s:%d----%v", file, line, fn))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				// do something or record log
				log.Println("exec panic error: ", err)
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
				// log.Println(debug.Stack())
			}
		}()
		c.Next()
		

	}
	
}