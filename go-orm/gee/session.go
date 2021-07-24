package gee

import (
	"database/sql"
	"strings"
)

type Session struct {
	db      *sql.DB
	sqlStmt strings.Builder
	sqlVars []interface{}
}

func (s *Session) Reset() {
	s.sqlStmt.Reset()
	s.sqlVars = s.sqlVars[0:0]
}

func (s *Session) Exec(stmt string, args ...interface{}) (sql.Result, error) {
	s.sqlStmt.WriteString(stmt)
	for _, arg := range args {
		s.sqlVars = append(s.sqlVars, arg)
	}
	return s.db.Exec(s.sqlStmt.String(), args...)
}

func (s *Session) Query(stmt string, args ...interface{}) (*sql.Rows, error) {
	s.Reset()
	s.sqlStmt.WriteString(stmt)
	for _, arg := range args {
		s.sqlVars = append(s.sqlVars, arg)
	}
	return s.db.Query(s.sqlStmt.String(), args...)
}
