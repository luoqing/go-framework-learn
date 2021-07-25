package gee

import (
	"reflect"
	"testing"
)

func testSelect(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "User", []string{"*"})
	clause.Set(WHERE, "Name = ?", "Tom")
	clause.Set(ORDERBY, "Age ASC")
	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

type User struct {
	ID   int
	Name string
	Age  int
}

func testInsert(t *testing.T) {
	/*
		user1 := User{
			ID:   1,
			Name: "tsing",
			Age:  8,
		}
		user2 := User{
			ID:   2,
			Name: "lo",
			Age:  10,
		}
		users := []interface{}{user1, user2} // 要支持是按照struct的方式
		// 为什么要抽象下面这种，因为上面的struct到values依赖于struct的数据结构
	*/
	var clause Clause
	values1 := []interface{}{1, "tsing", 8}
	values2 := []interface{}{2, "lo", 10}
	values := []interface{}{values1, values2}
	fields := []string{"ID", "Name", "Age"}
	clause.Set(INSERT, "User", fields)
	clause.Set(VALUES, values...)
	sql, vars := clause.Build(INSERT, VALUES)
	t.Log(sql, vars)
	if sql != "INSERT INTO User (ID,Name,Age) VALUES (?, ?, ?), (?, ?, ?)" {
		t.Fatal("failed to build INSERT SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{1, "tsing", 8, 2, "lo", 10}) {
		t.Fatal("failed to build INSERT SQLVars")
	}

}

func testUpdate(t *testing.T) {
	var clause Clause
	clause.Set(UPDATE, "User", []interface{}{"Name", "lq"}, []interface{}{"Age", 15})
	clause.Set(WHERE, "ID = ? AND Name != ?", 1, "tsing")
	sql, vars := clause.Build(UPDATE, WHERE)
	t.Log(sql, vars)
	if sql != "UPDATE User SET Name = ?,Age = ? WHERE ID = ? AND Name != ?" {
		t.Fatal("failed to build UPDATE SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"lq", 15, 1, "tsing"}) {
		t.Fatal("failed to build UPDATE SQLVars")
	}

}

func testDelete(t *testing.T) {
	var clause Clause
	clause.Set(DELETE, "User")
	clause.Set(WHERE, "ID = ?", 1)
	sql, vars := clause.Build(DELETE, WHERE)
	t.Log(sql, vars)
	if sql != "DElETE FROM User WHERE ID = ?" {
		t.Fatal("failed to build DELETE SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{1}) {
		t.Fatal("failed to build DELETE SQLVars")
	}

}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
	t.Run("insert", func(t *testing.T) {
		testInsert(t)
	})
	t.Run("update", func(t *testing.T) {
		testUpdate(t)
	})
	t.Run("delete", func(t *testing.T) {
		testDelete(t)
	})
}
