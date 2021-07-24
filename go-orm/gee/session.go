package gee

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Session struct {
	db      *sql.DB
	sqlStmt strings.Builder
	sqlVars []interface{}
	table   *RefTable
	dialect Dialect
}

func (s *Session) Reset() {
	s.sqlStmt.Reset()
	s.sqlVars = s.sqlVars[0:0]
}

func (s *Session) Model(tbStruct interface{}) *Session {
	s.table = StructToTable(tbStruct, s.dialect)
	return s
}

func (s *Session) Create(tbStruct interface{}) error {
	// 这个需要拼接create table的语句
	/* tableName, fieldsname, type
	Create Table tb (
		app_id int(11),
		app_key varchar(64),
		app_name vachar(128),
	)
	*/
	s.Reset()
	if s.table == nil {
		s.table = StructToTable(tbStruct, s.dialect)
	}
	s.sqlStmt.WriteString("CREATE TABLE ")
	s.sqlStmt.WriteString(s.table.TableName)
	s.sqlStmt.WriteString("(\n")
	for i, name := range s.table.FieldNames {
		field, ok := s.table.Name2Field[name]
		if !ok {
			return errors.New("error table Name2Field")
		}
		var col string
		if i == len(s.table.FieldNames)-1 {
			col = fmt.Sprintf("%s %s\n", name, field.Type)
		} else {
			col = fmt.Sprintf("%s %s,\n", name, field.Type)
		}
		s.sqlStmt.WriteString(col)
		i++
	}
	s.sqlStmt.WriteString(")\n")
	fmt.Println(s.sqlStmt.String())
	_, err := s.db.Exec(s.sqlStmt.String(), s.sqlVars...)
	return err

}

func (s *Session) Insert(tbStruct interface{}) (sql.Result, error) {
	s.Reset()
	if s.table == nil {
		s.table = StructToTable(tbStruct, s.dialect)
	}
	s.sqlStmt.WriteString("INSERT INTO ")
	s.sqlStmt.WriteString(s.table.TableName)
	fieldsname := fmt.Sprintf("(%s)", strings.Join(s.table.FieldNames, ","))
	s.sqlStmt.WriteString(fieldsname)
	s.sqlStmt.WriteString(" VALUES(")
	for i := 0; i < len(s.table.FieldNames)-1; i++ {
		s.sqlStmt.WriteString("?, ")
		name := s.table.FieldNames[i]
		field, ok := s.table.Name2Field[name]
		if !ok {
			return nil, errors.New("error table Name2Field")
		}
		s.sqlVars = append(s.sqlVars, field.Value)
	}
	s.sqlStmt.WriteString("?)")
	name := s.table.FieldNames[len(s.table.FieldNames)-1]
	field, ok := s.table.Name2Field[name]
	if !ok {
		return nil, errors.New("error table Name2Field")
	}
	s.sqlVars = append(s.sqlVars, field.Value)
	fmt.Println(s.sqlStmt.String())
	return s.db.Exec(s.sqlStmt.String(), s.sqlVars...)
}

func (s *Session) Exec(stmt string, args ...interface{}) (sql.Result, error) {
	s.Reset()
	s.sqlStmt.WriteString(stmt)
	for _, arg := range args {
		s.sqlVars = append(s.sqlVars, arg)
	}
	return s.db.Exec(s.sqlStmt.String(), s.sqlVars...)
}

func (s *Session) Query(stmt string, args ...interface{}) (*sql.Rows, error) {
	s.Reset()
	s.sqlStmt.WriteString(stmt)
	for _, arg := range args {
		s.sqlVars = append(s.sqlVars, arg)
	}
	return s.db.Query(s.sqlStmt.String(), s.sqlVars...)
}
