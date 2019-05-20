package db

import (
	"database/sql"
	"fmt"
	"github.com/pibigstar/go-cloudstore/db/mysql"
)

func UserSignup(username string,passwd string) bool  {
	stmt, err := mysql.DBConn().Prepare("insert tbl_user set user_name=?,user_pwd=?")
	if err != nil {
		fmt.Printf("Failed to prepare sql,err:%s\n",err.Error())
		return false
	}
	result, err := stmt.Exec(username, passwd)
	if err!=nil {
		fmt.Printf("Failed to prepare sql,err:%s\n",err.Error())
		return false
	}

	// 判断是否插入了,有时候即使SQL查询成功了，但并没有影响行记录
	if row, err := result.RowsAffected();err == nil {
		if row <= 0 {
			fmt.Println("the row has in the sql")
			return false
		}
	}
	return  true

}

type TableUser struct {
	UserName sql.NullString
	Password sql.NullString
}

func UserLogin(username string,password string) *TableUser  {
	mysql.DBConn().Prepare("select user_name,user_pwd where user_name=? and user_pwd")
	return nil
}
