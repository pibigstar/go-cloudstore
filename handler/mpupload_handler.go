package handler

import (
	"fmt"
	"github.com/pibigstar/go-cloudstore/db"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pibigstar/go-cloudstore/db/redis"
	"github.com/pibigstar/go-cloudstore/utils"
	goRedis "github.com/gomodule/redigo/redis"
)

// 分块上传Handler

type MultipartUploadInfo struct {
	FileHash   string
	UploadID   string
	FileSize   int
	ChunkSize  int //分块大小
	ChunkCount int //分块数量
}

const (
	ChunkSize = 5 * 1024 * 1024
	ChunkDataDIR = "D://data/"
)

// 初始化分块上传
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		fmt.Printf("Failed to parse file size to int,err:%s\n", err.Error())
		return
	}
	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	// 生成分块上传的初始化信息
	uploadInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  ChunkSize, // 5M
		ChunkCount: int(math.Ceil(float64(filesize / ChunkSize))),
	}
	// 将初始化信息写入到redis缓存
	rConn.Do("HSET", "MP_"+uploadInfo.UploadID, "chunkcount", uploadInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+uploadInfo.UploadID, "filehash", uploadInfo.FileHash)
	rConn.Do("HSET", "MP_"+uploadInfo.UploadID, "filesize", uploadInfo.FileSize)

	// 将初始化信息返回给客户端
	resp := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: uploadInfo,
	}
	w.Write(resp.JSONBytes())
}

// 上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	//username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	// 文件分块索引
	chunkIndex := r.Form.Get("index")

	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	// 获取文件句柄，用户存储分块内容
	filePath := ChunkDataDIR + uploadID + "/" + chunkIndex
	os.MkdirAll(filePath,0744)
	file, err := os.Create(filePath)
	if err != nil {
		w.Write(utils.NewRespMsg(-1,"upload part failed",nil).JSONBytes())
		return
	}
	defer file.Close()
	buff := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buff)
		if err!=nil{
			break
		}
		file.Write(buff[:n])
	}
	// 更新redis缓存状态
	rConn.Do("HSET","MP_"+uploadID,"chkidx_"+chunkIndex,1)

	// 返回客户端
	w.Write(utils.NewRespMsg(0,"OK",nil).JSONBytes())
}
// 上传合并
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	username := r.Form.Get("username")
	uploadid := r.Form.Get("uploadid")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	values, err := goRedis.Values(rConn.Do("HGETALL", "MP_"+uploadid))
	if err != nil {
		w.Write(utils.NewRespMsg(-1,"complete upload failed",nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i:=0;i<len(values);i+=2 {
		k := string(values[i].([]byte))
		v := string(values[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		}else if strings.HasPrefix(k,"chkidx_") && v =="1"{
			chunkCount+=1
		}
	}
	if totalCount!=chunkCount {
		w.Write(utils.NewRespMsg(-1,"invalid request",nil).JSONBytes())
		return
	}
	// 更新文件表和用户文件表
	db.InsertFile(filehash,filename,"",int64(filesize))
	db.CreateUserFile(username,filehash,filename,filesize)

	w.Write(utils.NewRespMsg(0,"ok",nil).JSONBytes())
}
