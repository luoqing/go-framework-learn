package gee

import (
	"fmt"
	"log"
	"testing"
)

func TestEngine(t *testing.T) {
	source := "mysql"
	dns := "root:123654@tcp(127.0.0.1:3306)/video_test?charset=utf8"
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("connect db:%v", err)
		}
	}()
	g := NewEngine(dns, source)
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
