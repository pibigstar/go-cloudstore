package process

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pibigstar/go-cloudstore/mq"
	"github.com/pibigstar/go-cloudstore/store/oss"
)

func ProcessTransfer(msg []byte) bool {
	// 解析msg
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	fmt.Printf("%+v\n", pubData)
	// 根据临时存储文件路径，创建文件句柄
	fin, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	// 通过文件句柄将文件内容读出来，并上传到OSS
	err = oss.Bucket().PutObject(
		pubData.DestLocation,
		bufio.NewReader(fin))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	// 更新文件的存储路径到文件表
	return true
}
