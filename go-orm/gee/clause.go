package gee

// todo：mock
// todo:singleton

import (
	"fmt"
	"strings"
)

// 语句又做了一层抽象
type generator func(values ...interface{}) (string, []interface{})
type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
)

type Clause struct {
	sqls map[Type]string
	vars map[Type][]interface{}
}

func (c *Clause) Set(typ Type, values ...interface{}) {
	if c.sqls == nil {
		c.sqls = make(map[Type]string)
		c.vars = make(map[Type][]interface{})
	}
	gen, ok := generators[typ]
	if !ok {
		panic("unknown generator")
	}
	stmt, vars := gen(values...)
	c.sqls[typ] = stmt
	c.vars[typ] = vars
}

// 将多行组成一行
// Find, First, Scan, Create, Insert,Update, Delete

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sqls[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.vars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars

}

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderby
	generators[UPDATE] = _update
	generators[DELETE] = _delete
}

func _select(values ...interface{}) (string, []interface{}) {
	if len(values) == 0 {
		panic("empty slect fields")
	}
	fmt.Println(values[0])
	stmt := fmt.Sprintf("SELECT %s FROM %s", strings.Join(values[1].([]string), ","), values[0])
	var vars []interface{}
	return stmt, vars

}

// 这些都没做字段检查，比如这个字段要是一个证书
func _limit(values ...interface{}) (string, []interface{}) {
	stmt := "LIMIT ?"
	vars := []interface{}{values[0]}
	if len(values) == 2 {
		stmt += ", ?"
		vars = append(vars, values[1])
	}
	return stmt, vars
}

func _where(values ...interface{}) (string, []interface{}) {

	stmt := fmt.Sprintf("WHERE %s", values[0])
	var vars []interface{}
	vars = append(vars, values[1:]...)
	return stmt, vars
}

func _orderby(values ...interface{}) (string, []interface{}) {
	if len(values) != 1 {
		panic("error orderby")
	}
	stmt := fmt.Sprintf("ORDER BY %s", values[0])
	var vars []interface{}
	return stmt, vars
}
func _insert(values ...interface{}) (string, []interface{}) {
	fields := strings.Join(values[1].([]string), ",")
	stmt := fmt.Sprintf("INSERT INTO %s (%s)", values[0], fields)

	return stmt, values[1:]
}
func placeHolders(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

// 插入多行
func _values(values ...interface{}) (string, []interface{}) {
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = placeHolders(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _update(values ...interface{}) (string, []interface{}) {
	stmt := fmt.Sprintf("UPDATE %s SET", values[0])

	return stmt, values[1:]
}

func _delete(values ...interface{}) (string, []interface{}) {
	stmt := fmt.Sprintf("DElETE FROM %s", values[0])
	var vars []interface{}
	return stmt, vars
}
