package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pibigstar/go-cloudstore/db"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	goRedis "github.com/gomodule/redigo/redis"
	"github.com/pibigstar/go-cloudstore/db/redis"
	"github.com/pibigstar/go-cloudstore/utils"
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
	ChunkSize    = 5 * 1024 * 1024
	ChunkDataDIR = "D://data/"
)

// 初始化分块上传
func InitialMultipartUploadHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
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
	c.Data(http.StatusOK,"application/json",resp.JSONBytes())
}

// 上传文件分块
func UploadPartHandler(c *gin.Context) {

	//username := r.Form.Get("username")
	uploadID := c.Request.FormValue("uploadid")
	// 文件分块索引
	chunkIndex := c.Request.FormValue("index")

	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	// 获取文件句柄，用户存储分块内容
	filePath := ChunkDataDIR + uploadID + "/" + chunkIndex
	os.MkdirAll(filePath, 0744)
	file, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"code": -1,
			"msg": "upload part failed",
		})
		return
	}
	defer file.Close()
	buff := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buff)
		if err != nil {
			break
		}
		file.Write(buff[:n])
	}
	// 更新redis缓存状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 返回客户端
	c.JSON(http.StatusOK, gin.H{
		"msg": "OK!",
		"code": 0,
	})
}

// 上传合并
func CompleteUploadHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	uploadid := c.Request.FormValue("uploadid")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	values, err := goRedis.Values(rConn.Do("HGETALL", "MP_"+uploadid))
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"code": 0,
			"msg": "complete upload failed",
		})
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(values); i += 2 {
		k := string(values[i].([]byte))
		v := string(values[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount += 1
		}
	}
	if totalCount != chunkCount {
		c.JSON(http.StatusInternalServerError,gin.H{
			"code": -1,
			"msg": "invalid request",
		})
		return
	}
	// 更新文件表和用户文件表
	db.InsertFile(filehash, filename, "", int64(filesize))
	db.CreateUserFile(username, filehash, filename, filesize)

	c.JSON(http.StatusOK,gin.H{
		"code": 0,
		"msg": "OK",
	})
}
