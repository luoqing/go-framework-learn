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

type Schema struct {
	//fields []*Filed
	Model      interface{} // 存储struct
	FieldNames []string
	name2Field map[string]*Field // 这个为啥要存成private，是为了不让修改是吗
	TableName  string            // table的名称
}

func (s *Schema) GetField(name string) (*Field, bool) {
	field, ok := s.name2Field[name]
	return field, ok
}

func (s *Schema) SetField(name string, field *Field) {
	s.name2Field[name] = field
}

func StructToTable(tbStruct interface{}, dialect Dialect) *Schema {
	r := &Schema{
		Model:      tbStruct,
		name2Field: make(map[string]*Field),
	}
	// 通过反射获取数据的filed，获取其name，value，tag
	typ := reflect.TypeOf(tbStruct)
	tb := typ.Name()
	r.TableName = Camel2Case(tb)
	val := reflect.ValueOf(tbStruct)
	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i).Name
		if alias, ok := typ.Field(i).Tag.Lookup("db"); ok {
			fieldName = alias
		} else {
			// 转化为下划线的
			fieldName = Camel2Case(fieldName)
		}
		dbType := dialect.DataTypeOf(val.Field(i))

		//dbType := goKind2DBType(typeName)
		//typeName := val.Type().Field(i).Type.Name()
		//dbType := go2DBType(typeName)
		if dbType == "varchar" {
			dbType = "varchar(255)"
			if len, ok := typ.Field(i).Tag.Lookup("len"); ok {
				dbType = fmt.Sprintf("varchar(%s)", len)
			}
		}
		// val.Field(i)是Value
		// typ.Field(i)是structField
		f := &Field{
			Name:  fieldName,
			Value: val.Field(i).Interface(),
			Tag:   typ.Field(i).Tag,
			Type:  dbType,
		}
		r.FieldNames = append(r.FieldNames, fieldName)
		r.SetField(fieldName, f)
	}
	return r
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
