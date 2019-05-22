package utils

import (
	"github.com/pibigstar/go-cloudstore/constant"
	"os"
	"time"
)

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FormatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 生成一个40位字符的token
func GenToken(username string) string {
	tokenPrefix := MD5([]byte(username+cont.TOKEN_SALT))
	return tokenPrefix
}
