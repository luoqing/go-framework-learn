package gee

import (
	"fmt"
	"reflect"
	"unicode"
)

// 将struct转化为 table fields

type Field struct {
	Name  string
	Tag   reflect.StructTag // 解析tag
	Value interface{}       // 需要将golang的类型和db的类型对应起来，尤其是在create table
	Type  string
}

type RefTable struct {
	//fields []*Filed
	FieldNames []string
	Name2Field map[string]*Field
	TableName  string // table的名称
}

func Camel2Case(name string) string {
	var buffer []rune
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer = append(buffer, '_')
			}
			buffer = append(buffer, unicode.ToLower(r))
		} else {
			buffer = append(buffer, unicode.ToLower(r))
		}
	}
	return string(buffer)
}

func go2DBType(typeName string) string {
	switch typeName {
	case "bool":
		return "tinyint(4)"
	case "int8":
		return "tinyint(4)"
	case "int16":
		return "int(11)"
	case "int32":
		return "int(11)"
	case "int64":
		return "bigint"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "Time":
		return "datetime"
	case "string": // text
		return "varchar"

	}
	return ""
}

func goKind2DBType(kind reflect.Kind) string {
	switch kind {
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

	}
	return ""

}

func StructToTable(tbStruct interface{}) *RefTable {
	r := &RefTable{}
	// 通过反射获取数据的filed，获取其name，value，tag
	typ := reflect.TypeOf(tbStruct)
	tb := typ.Name()
	r.TableName = Camel2Case(tb)
	r.Name2Field = make(map[string]*Field)
	val := reflect.ValueOf(tbStruct)
	for i := 0; i < typ.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		if alias, ok := typ.Field(i).Tag.Lookup("db"); ok {
			fieldName = alias
		} else {
			// 转化为下划线的
			fieldName = Camel2Case(fieldName)
		}
		//typeName := val.Type().Field(i).Type.Kind()

		//dbType := goKind2DBType(typeName)
		typeName := val.Type().Field(i).Type.Name()
		dbType := go2DBType(typeName)
		if dbType == "varchar" {
			dbType = "varchar(255)"
			if len, ok := typ.Field(i).Tag.Lookup("len"); ok {
				dbType = fmt.Sprintf("varchar(%s)", len)
			}
		}

		f := &Field{
			Name:  fieldName,
			Value: val.Field(i).Interface(),
			Tag:   typ.Field(i).Tag,
			Type:  dbType,
		}
		r.FieldNames = append(r.FieldNames, fieldName)
		r.Name2Field[fieldName] = f
	}
	return r
}
