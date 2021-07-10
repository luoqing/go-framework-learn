package lru

import(
	"testing"
)

// table-drived
// t.Errorf, t.Fatalf
// t.Helper()
// setup() teardown()
func TestNew(t *testing.T) {
	cache := New(2)
	cases := []struct{
		Key, Value string
	}{
		{"name", "tsing"},
		{"age", "1"},
		{"Work", "coding"},
	}
	for _, c := range cases {
		t.Run("set", func(t *testing.T){
			cache.Set(c.Key, c.Value)
		})
		t.Run("get", func(t *testing.T){
			v, err := cache.Get(c.Key)
			if err != nil {
				t.Errorf("get key:%s value failed:%v", c.Key, err)
				return
			}
			if v != c.Value {
				t.Errorf("key:%s value:%s is not expectd:%s", c.Key, v, c.Value)
				return
			}
		})
	}

	for _, c := range cases {
		t.Run("getAgain", func(t *testing.T){
			v, err := cache.Get(c.Key)
			if err != nil {
				t.Errorf("get key:%s value failed:%v", c.Key, err)
				return
			}
			if v != c.Value {
				t.Errorf("key:%s value:%s is not expectd:%s", c.Key, v, c.Value)
				return
			}
		})
	}
}

func TestGet(t *testing.T) {

}

func TestSet(t *testing.T) {

}


func TestRemoveOldest(t *testing.T) {

}