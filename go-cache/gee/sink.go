package gee

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
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

func (s *SMap) Set(key string, value []byte) error {
	s.data[key] = value
	return nil
}

func (s *SMap) Get(key string) ([]byte, error) {
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

func (s *FileSinker) Set(key string, value []byte) error {
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

func (s *FileSinker) Get(key string) ([]byte, error) {
	// 按照文件读取
	file := s.dataPath + "/" + key
	f, err := os.Open(file)
	value, err := ioutil.ReadAll(f)
	if err != nil {
		return value, err
	}
	return value, nil
}

type DBSinker struct {
	db         *sql.DB
	cacheTable string // 创建一个key-value的table，然后可以根据key去search到table
}

// db open
// 甚至可以连接redis，去获取数据。
func NewDBSinker(dns string) *DBSinker {
	//db, err := sql.Open("mysql", "root:123654@tcp(127.0.0.1:3306)/video_test?charset=utf8")
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	return &DBSinker{
		db: db,
	}
}

// db set
func (s *DBSinker) Set(key string, value []byte) error {
	return nil
}

// db get
func (s *DBSinker) Get(key string) ([]byte, error) {
	return nil, nil
}
