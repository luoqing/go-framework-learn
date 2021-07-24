package gee

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var g *Engine

func init() {
	source := "mysql"
	dns := "root:123654@tcp(127.0.0.1:3306)/video_test?charset=utf8"
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("connect db:%v", err)
		}
	}()
	g = NewEngine(dns, source)
}
func TestEngine(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("connect db:%v", err)
		}
	}()
	defer g.Close()
	s := g.NewSession()
	stmt := "SELECT Fapp_key, Fapp_name FROM t_access_app_conf WHERE Fapp_id = ?"
	rows, err := s.Query(stmt, 12)
	if err != nil {
		t.Fatalf("query rows failed:%v", err)
	}
	defer rows.Close()
	var appkey, appname string
	for rows.Next() {
		if err := rows.Scan(&appkey, &appname); err != nil {
			log.Fatalf("row scan error")
		}
		fmt.Println(appkey)
		fmt.Println(appname)
	}

}

func TestCreateTable(t *testing.T) {
	type AppConf struct {
		AppID      int32  `db:"app_id"`
		AppName    string `len:"128"`
		InsertTime time.Time
	}
	s := g.NewSession()
	var conf AppConf
	err := s.Create(conf)
	if err != nil {
		t.Errorf("create table failed：%v", err)
	}
}

func TestInsertTable(t *testing.T) {
	type AppConf struct {
		AppID      int32  `db:"app_id"`
		AppName    string `len:"128"`
		InsertTime time.Time
	}
	s := g.NewSession()
	var conf = AppConf{
		AppID:      11,
		AppName:    "dfdfd",
		InsertTime: time.Now(),
	}
	_, err := s.Insert(conf)
	if err != nil {
		t.Errorf("insert table failed：%v", err)
	}

}
