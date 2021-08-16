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

func TestTransaction(t *testing.T) {
	s := g.NewSession()
	var conf = AppConf{
		AppID:      11,
		AppName:    "update 1",
		InsertTime: time.Now(),
	}
	var appID int32 = 11
	fn := func(s *Session) error {
		_, err := s.Where("app_id = ?", appID).Update(conf)
		if err != nil {
			return err
		}
		conf.AppName = "update 2"
		_, err = s.Where("app_id = ?", appID).Update(conf)
		if err != nil {
			return err
		}
		return nil
	}
	err := s.Transaction(fn)
	if err != nil {
		t.Errorf("tx error:%v", err)
	}
}

// 这些考虑过并发吗
// 最好是一个事务一个独立的sesssion， 这种方法会好点
// 上面session粒度如果有两个事务公用一个s调用Transaction，则tx有可能出现问题，fn中使用select有默认的limit，这个就有可能混淆。
// 综上，下面的这种事务实现Transaction与单独使用begin等不会混淆，所以tx还是要单独new session
// https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
func TestEngineTransaction(t *testing.T) {
	var conf = AppConf{
		AppID:      91,
		AppName:    "dfdfd89fdfd",
		InsertTime: time.Now(),
	}

	fn := func(s *Session) (interface{}, error) {
		res, err := s.Insert(conf)
		if err != nil {
			return res, err
		}
		conf.AppName = "update 2"
		res, err = s.Where("app_id = ?", conf.AppID).Update(conf)
		if err != nil {
			return res, err
		}
		return res, nil
	}
	res, err := g.Transaction(fn)
	if err != nil {
		t.Errorf("engine tx error:%v", err)
	}
	fmt.Println(res)
}
