package lru

import(
	"testing"
)

func TestNew(t *testing.T) {
	lru := New(2)
	lru.Set("name", "tsing")
	lru.Get("name")
	lru.Set("age", 1)
	lru.Get("age")
	lru.Set("name", "lo")
	lru.Get("name")
	lru.Set("Work", "coding")
	lru.Get("Work")
	lru.Get("name")
	lru.Get("age")


}

func TestGet(t *testing.T) {

}

func TestSet(t *testing.T) {

}


func TestRemoveOldest(t *testing.T) {

}