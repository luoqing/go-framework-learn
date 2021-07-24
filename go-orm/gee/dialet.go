package gee

import (
	"fmt"
	"reflect"
	"time"
)

// golang type 转化其他数据类型，比如mysql,sqlite,pg的时候这个地方不一样
// open的datasource也不一样
var dialectsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}

type MysqlDialect struct {
}

// dialect不太懂如何去抽象
func (d *MysqlDialect) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "tinyint(4)"
	case reflect.Int8:
		return "tinyint(4)"
	case reflect.Int:
		return "int(11)"
	case reflect.Int16:
		return "int(11)"
	case reflect.Int32:
		return "int(11)"
	case reflect.Int64:
		return "bigint"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	case reflect.String:
		return "varchar"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))

}

func (d *MysqlDialect) TableExistSQL(tableName string) (string, []interface{}) {
	stmt := "SELECT table_name FROM information_schema.TABLES WHERE table_name = ?"
	vals := []interface{}{tableName}
	return stmt, vals
}
