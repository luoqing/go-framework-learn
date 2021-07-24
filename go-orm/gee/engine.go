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
