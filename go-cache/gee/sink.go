package gee
import(
	"fmt"
	"os"
	"io/ioutil"
)

// 也测试下FileSinker
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

type FileSinker struct {
	dataPath string
}

func NewFileSinker(dataPath string) *FileSinker {
	return &FileSinker{
		dataPath: dataPath,
	}
}

func (s *FileSinker)Set(key string, value []byte) error {
	file := s.dataPath + "/" + key
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	_, err = f.Write(value)
	if err != nil {
		return err
	}
	
	// 按照key写文件
	return nil
}

func (s *FileSinker)Get(key string) ([]byte, error) {
	// 按照文件读取
	file := s.dataPath + "/" + key
	f, err := os.Open(file)
	value, err := ioutil.ReadAll(f)
	if err != nil {
		return value, err
	}
	return value, nil
}
