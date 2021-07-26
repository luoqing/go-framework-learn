package gee

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Session struct {
	db      *sql.DB
	sqlStmt strings.Builder
	sqlVars []interface{}
	table   *Schema
	dialect Dialect
	clause  Clause
}

func (s *Session) Reset() {
	s.sqlStmt.Reset()
	s.sqlVars = s.sqlVars[0:0]
}

func (s *Session) Model(tbStruct interface{}) *Session {
	if s.table == nil || reflect.TypeOf(tbStruct) != reflect.TypeOf(s.table.Model) {
		s.table = StructToTable(tbStruct, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *Schema {
	return s.table
}

func (s *Session) TableExist(tableName string) (bool, error) {
	s.Reset()
	var stmt string
	stmt, s.sqlVars = s.dialect.TableExistSQL(tableName)
	s.sqlStmt.WriteString(stmt)
	row, err := s.db.Query(s.sqlStmt.String(), s.sqlVars...)
	if err != nil {
		return false, err
	}
	if row.Next() {
		return true, nil
	}
	return false, nil
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
	exist, err := s.TableExist(s.table.TableName)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("table exists")
	}
	s.Reset()

	var cols []string
	for _, name := range s.table.FieldNames {
		field, ok := s.table.GetField(name)
		if !ok {
			return errors.New("error table Name2Field")
		}
		col := fmt.Sprintf("%s %s", name, field.Type) // 此处可以根据tag null default comemnt来补充createtable的语句
		cols = append(cols, col)
	}
	fields := strings.Join(cols, ",\n")

	s.sqlStmt.WriteString("CREATE TABLE ")
	s.sqlStmt.WriteString(s.table.TableName)
	s.sqlStmt.WriteString("(\n")
	s.sqlStmt.WriteString(fields)
	s.sqlStmt.WriteString(")\n")
	fmt.Println(s.sqlStmt.String())
	_, err = s.db.Exec(s.sqlStmt.String(), s.sqlVars...)
	return err

}
func (s *Session) Select(values ...string) *Session {
	s.clause.Set(SELECT, s.table.TableName, values)
	return s
}

func (s *Session) Where(values ...string) *Session {
	s.clause.Set(WHERE, values)
	return s
}

func (s *Session) Limit(n int) *Session {
	s.clause.Set(LIMIT, n)
	return s
}

func (s *Session) OrderBy(orderRule string) *Session {
	s.clause.Set(ORDERBY, orderRule)
	return s
}

func (s *Session) Find(values []interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem() // 获取类型
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	sql, vars := s.clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	rows, err := s.db.Query(sql, vars...)
	if err != nil {
		return err
	}
	// 进行数据的scan
	for rows.Next() {
		// new一个类型的实际例子
		dest := reflect.New(destType).Elem()

		var fields []interface{}
		for _, name := range table.FieldNames {
			fields = append(fields, dest.FieldByName(name).Addr().Interface())
		}
		/*
			val := reflect.ValueOf(dest)
			typ := reflect.TypeOf(dest)
			for i := 0; i < typ.NumField(); i++ {
				field := val.Field(i).Interface()
				fields = append(fields, field)
			}
		*/
		if err := rows.Scan(fields...); err != nil {
			return err
		}
		values = append(values, dest)
		//destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()

}

/*
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	//destType := destSlice.Type().Elem()
	//dest := reflect.New(destType).Elem()
	destType := dest.Type().Elem()
	destSlice := reflect.New(reflect.SliceOf(destType)).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}*/

func struct2record(tbStruct interface{}) (vars []interface{}) {
	typ := reflect.TypeOf(tbStruct)
	val := reflect.ValueOf(tbStruct)
	for i := 0; i < typ.NumField(); i++ {
		value := val.Field(i).Interface()
		vars = append(vars, value)
	}
	return
}

func (s *Session) Insert(tbStruct ...interface{}) (sql.Result, error) {
	// 一个struct对应一个record
	s = s.Model(tbStruct[0])
	fields := s.table.FieldNames
	s.clause.Set(INSERT, s.table.TableName, fields)
	var values []interface{}
	for _, st := range tbStruct {
		record := struct2record(st)
		values = append(values, record)
	}
	s.clause.Set(VALUES, values...)
	sql, vars := s.clause.Build(INSERT, VALUES)
	return s.db.Exec(sql, vars...)
}

func (s *Session) Update(tbStruct interface{}) (sql.Result, error) {
	s = s.Model(tbStruct)
	fields := s.table.FieldNames
	values := struct2record(tbStruct)
	updateVars := []interface{}{s.table.TableName}
	for i, value := range values {
		field := fields[i]
		v := []interface{}{field, value}
		updateVars = append(updateVars, v)
	}
	s.clause.Set(UPDATE, updateVars...)
	sql, vars := s.clause.Build(UPDATE, WHERE, LIMIT)
	return s.db.Exec(sql, vars...)
}

// 找到主键, 先用where
func (s *Session) Delete(tbStruct interface{}) (sql.Result, error) {
	s = s.Model(tbStruct)
	s.clause.Set(DELETE, s.table.TableName)
	sql, vars := s.clause.Build(DELETE, WHERE, LIMIT)
	return s.db.Exec(sql, vars...)
}

func (s *Session) Insert2(tbStruct interface{}) (sql.Result, error) {
	s.Reset()
	if s.table == nil {
		s.table = StructToTable(tbStruct, s.dialect)
	}
	// 以下都是build的
	s.sqlStmt.WriteString("INSERT INTO ") // _insert
	s.sqlStmt.WriteString(s.table.TableName)
	fieldsname := fmt.Sprintf("(%s)", strings.Join(s.table.FieldNames, ","))
	s.sqlStmt.WriteString(fieldsname)
	s.sqlStmt.WriteString(" VALUES(") // _values
	for i := 0; i < len(s.table.FieldNames)-1; i++ {
		s.sqlStmt.WriteString("?, ")
		name := s.table.FieldNames[i]
		field, ok := s.table.GetField(name)
		if !ok {
			return nil, errors.New("error table Name2Field")
		}
		s.sqlVars = append(s.sqlVars, field.Value)
	}
	s.sqlStmt.WriteString("?)")
	name := s.table.FieldNames[len(s.table.FieldNames)-1]
	field, ok := s.table.GetField(name)
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
