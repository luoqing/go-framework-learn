package gee
import(
	"fmt"
)

type Sinker interface {
	Set(key string, value []byte) error

	Get(key string) ([]byte, error)
}

type SMap struct {
	data map[string][]byte
}

func NewSMap() *SMap {
	return &SMap{
		data: make(map[string][]byte),
	}
}

func (s *SMap)Set(key string, value []byte) error {
	s.data[key] = value
	return nil
}

func (s *SMap)Get(key string) ([]byte, error) {
	value, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("key:%s not existed", key)
	}
	return value, nil
}

