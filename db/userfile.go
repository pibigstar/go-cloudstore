package db

import (
	"fmt"
	"github.com/pibigstar/go-cloudstore/db/mysql"
	"time"
)

type UserFile struct {
	UserName string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdated string
}
// 插入用户文件表
func CreateUserFile(username, filehash, filename string,filesize int) bool {
	stmt, err := mysql.DBConn().Prepare("insert tbl_user_file set user_name=?,file_name=?,file_sha1=?,file_size=?,upload_at=?")
	if err != nil {
		fmt.Printf("Failed to prepare sql,err:%s\n", err.Error())
		return false
	}
	_, err = stmt.Exec(username, filename, filehash, filesize, time.Now())
	if err != nil {
		fmt.Printf("Failed to exec sql,err:%s\n", err.Error())
		return false
	}
	return true
}
// 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mysql.DBConn().Prepare("select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name=? limit ?")
	if err != nil {
		fmt.Printf("Failed to prepare dd sql,err:%s\n", err.Error())
		return nil,err
	}
	rows, err := stmt.Query(username, limit)
	if err != nil {
		fmt.Printf("Failed to exec sql,err:%s\n", err.Error())
		return nil,err
	}
	var userFiles []UserFile
	for rows.Next(){
		userfile := UserFile{}
		err := rows.Scan(&userfile.FileHash, &userfile.FileName, &userfile.FileSize, &userfile.UploadAt, &userfile.LastUpdated)
		if err != nil {
			fmt.Printf("Failed to scan rows,err:%s\n", err.Error())
		}
		userFiles = append(userFiles, userfile)
	}
	return userFiles,nil
}
