package db

import (
	"database/sql"
	"fmt"
	"github.com/pibigstar/go-cloudstore/db/mysql"
	"log"
)

// 插入文件表
func InsertFile(filehash, filename, fileaddr string, filesize int64) bool {
	db := mysql.DBConn()
	stmt, err := db.Prepare("insert tbl_file set file_sha1=?,file_name=?,file_addr=?,file_size=?,status=?")
	if err != nil {
		fmt.Printf("Failed to prepare sql,err:%s\n", err.Error())
		return false
	}
	defer stmt.Close()

	result, err := stmt.Exec(filehash, filename, fileaddr, filesize, 1)
	if err != nil {
		fmt.Printf("Failed to exec sql,err:%s\n", err.Error())
		return false
	}

	// 判断是否插入了,有时候即使SQL查询成功了，但并没有影响行记录
	if row, err := result.RowsAffected(); err == nil {
		if row <= 0 {
			fmt.Printf("File with hash:%s has been uploaded \n", filehash)
		}
	}
	return true
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// 根据Hash从mysql中获取元数据信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, _ := mysql.DBConn().Prepare("select file_sha1,file_addr,file_size,file_name from tbl_file where file_sha1=? and status=1")
	defer stmt.Close()
	tfile := TableFile{}
	err := stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileSize, &tfile.FileName)
	if err != nil {
		fmt.Printf("Failed to exec sql,err:%s\n", err.Error())
		return nil, err
	}
	return &tfile, nil
}

func RenameFile(filehash, name string) bool {
	stmt, err := mysql.DBConn().Prepare("update tbl_user_file set file_name=? where file_sha1=?")
	if err != nil {
		log.Println(err.Error())
		return false
	}
	_, err = stmt.Exec(name, filehash)
	if err != nil {
		return false
	}
	return true
}
