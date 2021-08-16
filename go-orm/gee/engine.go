package gee

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type Engine struct {
	db      *sql.DB
	dialect Dialect
}

func init() {
	RegisterDialect("mysql", &MysqlDialect{})
}

func NewEngine(dns, source string) *Engine {
	//db, err := sql.Open("mysql", "root:123654@tcp(127.0.0.1:3306)/video_test?charset=utf8")
	db, err := sql.Open(source, dns)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	dialect, ok := GetDialect(source)
	if !ok {
		panic(errors.New("get dialect failed"))
	}
	return &Engine{
		db:      db,
		dialect: dialect,
	}
}

func (g *Engine) Close() {
	g.db.Close()
}

func (g *Engine) NewSession() *Session {
	return &Session{
		db:      g.db,
		dialect: g.dialect,
	}
}

// 这种写法相对来说更加通用一些，因为update等返回不一定
type TxFunc func(*Session) (interface{}, error)

func (g *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := g.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}
