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
func TestQuery(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("connect db:%v", err)
		}
	}()
	defer g.Close()
	s := g.NewSession()
	stmt := "SELECT app_name FROM app_conf WHERE app_id = ?"
	rows, err := s.Query(stmt, 11)
	if err != nil {
		t.Fatalf("query rows failed:%v", err)
	}
	defer rows.Close()
	var appname string
	for rows.Next() {
		if err := rows.Scan(&appname); err != nil {
			log.Fatalf("row scan error")
		}
		fmt.Println(appname)
	}

}

func TestCreateTable(t *testing.T) {

	type AppConf2 struct {
		AppID      int32  `db:"app_id"`
		AppName    string `len:"128"`
		InsertTime time.Time
	}
	s := g.NewSession()
	var conf AppConf2
	err := s.Drop(conf)
	if err != nil {
		t.Errorf("drop table failed：%v", err)
	}

	err = s.Create(conf)
	if err != nil {
		t.Errorf("create table failed：%v", err)
	}
}

type AppConf struct {
	AppID      int32  `db:"app_id"`
	AppName    string `len:"128"`
	InsertTime time.Time
}

func TestInsert2Table(t *testing.T) {

	s := g.NewSession()
	var conf = AppConf{
		AppID:      11,
		AppName:    "dfdfd",
		InsertTime: time.Now(),
	}
	_, err := s.Insert2(conf)
	if err != nil {
		t.Errorf("insert table failed：%v", err)
	}
}

// has one problem
func TestSelect(t *testing.T) {
	s := g.NewSession()
	type Conf struct {
		AppID      int32  `db:"app_id"`
		AppName    string `len:"128"`
		InsertTime string
	}
	var confs []Conf
	var conf AppConf
	var appID int32 = 11
	// 无法获取insertTime类型为time.Time的字段
	// 这里必须先model，因为select需要tableName
	err := s.Model(conf).Select("app_id, app_name,insert_time").Where("app_id = ?", appID).Limit(2).Find(&confs)
	if err != nil {
		t.Errorf("select error:%v", err)
	}
	Info("%v", confs)
}

func TestSelectOne(t *testing.T) {
	s := g.NewSession()
	type Conf struct {
		AppID      int32  `db:"app_id"`
		AppName    string `len:"128"`
		InsertTime string
	}
	var conf AppConf
	var appID int32 = 11
	var rst Conf
	// 无法获取insertTime类型为time.Time的字段
	err := s.Model(conf).Select("app_id, app_name,insert_time").Where("app_id = ?", appID).First(&rst) // 必须传指针，否则无法根据Address()来进行赋值
	if err != nil {
		t.Errorf("select error:%v", err)
	}
	Info("%v", rst)
}

func TestMultiInsert(t *testing.T) {
	var conf1 = AppConf{
		AppID:      21,
		AppName:    "dfdfd89fdfd",
		InsertTime: time.Now(),
	}

	var conf2 = AppConf{
		AppID:      22,
		AppName:    "dfdfd",
		InsertTime: time.Now(),
	}
	s := g.NewSession()
	_, err := s.Insert(conf1, conf2)
	if err != nil {
		t.Errorf("insert error:%v", err)
	}
}

func TestUpdate(t *testing.T) {
	var conf = AppConf{
		AppID:      11,
		AppName:    "fdfdfdfdupdate",
		InsertTime: time.Now(),
	}
	var appID int32 = 11
	s := g.NewSession()
	_, err := s.Where("app_id = ?", appID).Update(conf)
	if err != nil {
		t.Errorf("update error:%v", err)
	}

}

func TestDelete(t *testing.T) {
	s := g.NewSession()
	var conf AppConf
	var appID int32 = 11
	_, err := s.Where("app_id = ?", appID).Delete(conf)
	if err != nil {
		t.Errorf("delete error:%v", err)
	}
}

// todo 去写一个日志类，日志类会有不同的颜色
// level，writter(console, file)
// 封装现有的日志，屏幕输出显示不同的颜色
