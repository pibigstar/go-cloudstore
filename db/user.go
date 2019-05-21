package db

import (
	"fmt"
	"github.com/pibigstar/go-cloudstore/db/mysql"
)

func UserSignup(username string, passwd string) bool {
	stmt, err := mysql.DBConn().Prepare("insert tbl_user set user_name=?,user_pwd=?")
	if err != nil {
		fmt.Printf("Failed to prepare sql,err:%s\n", err.Error())
		return false
	}
	result, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Printf("Failed to exec sql,err:%s\n", err.Error())
		return false
	}

	// 判断是否插入了,有时候即使SQL查询成功了，但并没有影响行记录
	if row, err := result.RowsAffected(); err == nil {
		if row <= 0 {
			fmt.Println("the row has in the sql")
			return false
		}
	}
	return true

}

type TableUser struct {
	Username string
	Password string
	Phone    string
	SignupAt string
	Email    string
	Status   int
}

func UserLogin(username string, password string) (*TableUser, error) {
	stmt, err := mysql.DBConn().Prepare("select user_name,phone from tbl_user where user_name=? and user_pwd=?")
	if err != nil {
		fmt.Printf("Failed to prepare sql,err:%s\n", err.Error())
		return nil, err
	}
	user := TableUser{}
	err = stmt.QueryRow(username, password).Scan(&user.Username, &user.Phone)
	if err != nil {
		fmt.Printf("Failed to exec sql,err:%s\n", err.Error())
		return nil, err
	}

	return &user, nil
}

func UpdateUserToken(username, token string) bool {
	stmt, err := mysql.DBConn().Prepare("insert tbl_user_token set user_name=?,user_token=?")
	if err != nil {
		fmt.Printf("Failed to prepare sql,err:%s\n", err.Error())
		return false
	}
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Printf("Failed to exec sql,err:%s\n", err.Error())
		return false
	}
	return true
}

func GetUserInfo(username string) (*TableUser, error) {
	stmt, err := mysql.DBConn().Prepare("select user_name,signup_at from tbl_user where user_name=?")
	if err != nil {
		fmt.Printf("Failed to prepare sql,err:%s\n", err.Error())
		return nil, err
	}
	user := TableUser{}

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Printf("Failed to exec sql,err:%s\n", err.Error())
		return nil, err
	}
	return &user, nil
}
