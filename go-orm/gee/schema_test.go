package gee

import (
	"fmt"
	"testing"
	"time"
)

func TestCamel2Case(t *testing.T) {
	tests := []struct {
		Name string
		Want string
	}{
		{"abc", "abc"},
		{"Abc", "abc"},
		{"AppName", "app_name"},
		{"appId", "app_id"}, // id 要特殊处理
		{"Is_lo", "is_lo"},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			after := Camel2Case(tt.Name)
			if after != tt.Want {
				t.Errorf("error convert:%s - %s", tt.Name, after)
			}
		})
	}
}

func TestCase2Camel(t *testing.T) {
	tests := []struct {
		Want string
		Name string
	}{
		{"Abc", "abc"},
		{"AppName", "app_name"},
		{"AppId", "app_id"}, // id 要特殊处理
		{"IsLo", "is_lo"},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			after := Case2Camel(tt.Name)
			if after != tt.Want {
				t.Errorf("error convert:%s - %s", tt.Name, after)
			}
		})
	}

}

func TestStructToTable(t *testing.T) {
	// 将struct转化为table
	type AppConf struct {
		AppID      int32  `db:"app_id"`
		AppName    string `len:"128"`
		InsertTime time.Time
	}
	var conf AppConf
	dialect := &MysqlDialect{}
	r := StructToTable(conf, dialect)
	fmt.Println(*r)
}
