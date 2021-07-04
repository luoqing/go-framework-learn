package gee
import(
	"testing"
	"fmt"
//	"os"
)

func TestParsePattern(t *testing.T) {
	r := NewRouter();
	//pattern := os.Args[1]
	//pattern := "/a/b/c"
	//pattern := "/a/:b/:c"
	//pattern := "/a/:b/*"
	pattern := "/a/*"
	fmt.Println(pattern)
	parts, params, err := r.parsePattern(pattern)
	fmt.Println(parts)
	fmt.Println(params)
	fmt.Println(err)
}