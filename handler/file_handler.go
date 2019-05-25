package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pibigstar/go-cloudstore/config"
	"github.com/pibigstar/go-cloudstore/constant"
	"github.com/pibigstar/go-cloudstore/db"
	"github.com/pibigstar/go-cloudstore/meta"
	"github.com/pibigstar/go-cloudstore/mq"
	"github.com/pibigstar/go-cloudstore/store/oss"
	"github.com/pibigstar/go-cloudstore/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 去上传页面
func ToUploadHandler(c *gin.Context)  {
	c.Redirect(http.StatusFound, "static/view/upload.html")
}
// 处理文件上传
func DoUploadHandler(c *gin.Context) {
	//POST请求是上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get file data, err:%s\n", err.Error())
		return
	}
	defer file.Close()
	// 2. 把文件内容转为[]byte
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Printf("Failed to get file data, err:%s\n", err.Error())
		return
	}
	// 判断临时上传本地路径是否存在
	exist, err := utils.PathExists(cont.UploadFilePath)
	if !exist {
		err := os.Mkdir(cont.UploadFilePath, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create file dir, err:%s\n", err.Error())
			return
		}
	}
	filePath := fmt.Sprintf("%s%s", cont.UploadFilePath, header.Filename)
	// 文件元数据
	fileMeta := meta.FileMeta{
		FileName: header.Filename,
		FileSha1: utils.Sha1(buf.Bytes()),
		FileSize: int64(len(buf.Bytes())),
		FilePath: cont.UploadFilePath,
		Location: filePath,
		UploadAt: utils.FormatTime(),
	}
	fmt.Println("file sha1:", fileMeta.FileSha1)

	// 创建一个新文件
	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to create new file, err:%s\n", err.Error())
		return
	}
	//defer newFile.Close()
	// 将上传的文件内容复制到新文件中
	nByte, err := newFile.Write(buf.Bytes())
	if int64(nByte) != fileMeta.FileSize || err != nil {
		log.Printf("Failed to save data into file, writtenSize:%d, err:%s\n", nByte, err.Error())
		return
	}
	// 将文件上传到OSS中
	// 游标重新回到文件头部
	newFile.Seek(0, 0)
	ossPath := config.OSSRootDir + fileMeta.FileSha1
	fileMeta.Location = ossPath

	if !config.AsyncTransferEnable {
		// 同步
		err = oss.Bucket().PutObject(ossPath, newFile)
		if err != nil {
			log.Println(err.Error())
			return
		}
		fileMeta.Location = ossPath
	} else {
		// 写入异步转移任务队列
		data := mq.TransferData{
			FileHash:      fileMeta.FileSha1,
			CurLocation:   filePath,
			DestLocation:  ossPath,
			DestStoreType: config.StoreOSS,
		}
		pubData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to marsha1 transfer data,err:%s\n",err.Error())
			return
		}
		pubSuc := mq.Publish(
			config.TransExchangeName,
			config.TransOSSRoutingKey,
			pubData,
		)
		if !pubSuc {
			// TODO: 当前发送转移信息失败，稍后重试
			fmt.Println("失败")
		}
	}

	// 将文件信息保存到mysql中
	meta.UpdateFileMetaDB(fileMeta)

	//更新用户文件表记录
	username := c.Request.FormValue("username")
	suc := db.CreateUserFile(username, fileMeta.FileSha1, fileMeta.FileName, int(fileMeta.FileSize))
	if suc {
		// 重定向路由
		c.Redirect(http.StatusFound,"/static/view/home.html")
	} else {
		c.JSON(http.StatusOK,gin.H{
			"msg": "Upload Failed!",
			"code": -1,
		})
	}
}

// 尝试使用秒传
func TryFastUploadHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	// 从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	//TODO: 秒传失败之后应该调用正常上传逻辑
	if err != nil {
		c.Status(http.StatusNotFound)
		//UploadHandler(w,r)
		return
	}
	if fileMeta == nil {
		resp := utils.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		c.JSON(http.StatusOK,resp)
		return
	}

	// 上传过了，触发秒传，将文件信息写入用户文件表
	suc := db.CreateUserFile(username, filehash, filename, filesize)
	if suc {
		resp := utils.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		c.JSON(http.StatusOK,resp)
		return
	} else {
		c.JSON(http.StatusOK,gin.H{
			"msg": "秒传失败",
			"code": -1,
		})
	}

}

// 获取文件元数据信息
func GetFileMeta(c *gin.Context) {
	hash := c.Request.FormValue("filehash")
	fileMeta, err := meta.GetFileMetaDB(hash)
	bytes, err := json.Marshal(fileMeta)
	if err != nil {
		fmt.Printf("Failed to parse fileMeta,err:%s\n", err.Error())
		c.Status(http.StatusInternalServerError)
	}
	c.Data(http.StatusOK,"application/json",bytes)
}

// 批量查询文件元数据信息
func QueryFileHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	limit, _ := strconv.Atoi(c.Request.FormValue("limit"))
	//fileMeta := meta.GetLastFileMeta(limit)
	files, err := db.QueryUserFileMetas(username, limit)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(files)
	if err != nil {
		fmt.Printf("Failed to parse fileMeta,err:%s\n", err.Error())
		c.Status(http.StatusInternalServerError)
	}
	c.Data(http.StatusOK,"application/json",bytes)
}
// 下载文件
func DownloadFileHandler(c *gin.Context) {

	hash := c.Request.FormValue("filehash")
	fileMeta,err := meta.GetFileMetaDB(hash)
	if err!= nil {
		c.Status(http.StatusInternalServerError)
	}
	// 从oss中下载
	if strings.HasPrefix(fileMeta.Location, config.OSSRootDir){
		hash := c.Request.FormValue("filehash")
		fileMeta,err := meta.GetFileMetaDB(hash)
		if err != nil {
			log.Println(err.Error())
			return
		}
		url := oss.DownloadURL(fileMeta.Location)
		fmt.Println(url)
		c.JSON(http.StatusOK,gin.H{
			"code": 0,
			"msg": url,
		})
	} else {
		//从本地下载
		file, err := os.Open(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to open the file, err:%s\n", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Printf("Failed to read the file, err:%s\n", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Header("Content-Type", "application/octect-stream")
		c.Header("Content-Disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
		c.Data(http.StatusOK, "application/json", bytes)
	}
}

// 重命名
func UpdateFileMetaHandler(c *gin.Context) {
	// 操作类型
	opType := c.Request.FormValue("op")
	hash := c.Request.FormValue("filehash")
	newFileName := c.Request.FormValue("filename")
	// 目前仅仅支持重命名，如果不是，则返回403
	if opType != "0" {
		c.Status(http.StatusForbidden)
		return
	}
	b := db.RenameFile(hash, newFileName)
	if !b {
		c.JSON(http.StatusOK,gin.H{
			"code": -1,
			"msg": "更新失败",
		})
	} else {
		c.JSON(http.StatusOK,gin.H{
			"code": 0,
			"msg": "OK",
		})
	}
}

// 删除文件
func DeleteFileHandler(c *gin.Context) {
	hash := c.Request.FormValue("filehash")

	fileMeta := meta.GetFileMeta(hash)
	err := os.Remove(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to delete the file, err:%s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	// 将此文件从元数据集合中删除
	meta.RemoveFileMeta(hash)

	c.JSON(http.StatusOK,gin.H{
		"msg": "OK!",
		"code": 0,
	})
}
