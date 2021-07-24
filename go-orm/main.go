package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // 应该是调用init
)

func main() {
	var db *sql.DB
	db, err := sql.Open("mysql", "root:123654@tcp(127.0.0.1:3306)/video_test?charset=utf8")
	defer db.Close()
	if err != nil {
		log.Fatalf("db connect failed")
	}
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping failed")
	}
	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT Fapp_key, Fapp_name FROM t_access_app_conf WHERE Fapp_id = ?", 12)
	if err != nil {

	}
	defer rows.Close()
	var appkey, appname string
	for rows.Next() {
		if err := rows.Scan(&appkey, &appname); err != nil {
			log.Fatalf("row scan error")
		}
		fmt.Println(appkey)
		fmt.Println(appname)
	}
}
