package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3307)/fileserver?charset=utf8")
	// 设置最大活动链接数
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Printf("Failed to connect the mysql,err:%s\n", err.Error())
		os.Exit(1)
	}
}

// 返回数据库链接对象
func DBConn() *sql.DB {
	return db
}
