package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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
)

// 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// GET请求是去上传页面
	if r.Method == http.MethodGet {
		bytes, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(bytes))
	} else if r.Method == http.MethodPost {
		//POST请求是上传文件
		file, header, err := r.FormFile("file")
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
				CurLocation:   fileMeta.Location,
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
		username := r.Form.Get("username")
		suc := db.CreateUserFile(username, fileMeta.FileSha1, fileMeta.FileName, int(fileMeta.FileSize))
		if suc {
			// 重定向路由
			http.Redirect(w, r, "/home", http.StatusFound)
		} else {
			io.WriteString(w, "Upload Failed!")
		}
	}
}

// 尝试使用秒传
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	//TODO: 秒传失败之后应该调用正常上传逻辑
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		//UploadHandler(w,r)
		return
	}
	if fileMeta == nil {
		resp := utils.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONBytes())
		return
	}

	// 上传过了，触发秒传，将文件信息写入用户文件表
	suc := db.CreateUserFile(username, filehash, filename, filesize)
	if suc {
		resp := utils.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	} else {
		io.WriteString(w, "秒传失败")
	}

}

// 上传成功
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload success")
}

// 获取文件元数据信息
func GetFileMeta(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hash := r.Form.Get("filehash")
	fileMeta, err := meta.GetFileMetaDB(hash)
	bytes, err := json.Marshal(fileMeta)
	if err != nil {
		fmt.Printf("Failed to parse fileMeta,err:%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(bytes)
}

// 批量查询文件元数据信息
func QueryFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	limit, _ := strconv.Atoi(r.Form.Get("limit"))
	//fileMeta := meta.GetLastFileMeta(limit)
	files, err := db.QueryUserFileMetas(username, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(files)
	if err != nil {
		fmt.Printf("Failed to parse fileMeta,err:%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(bytes)
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	hash := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(hash)

	file, err := os.Open(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to open the file, err:%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Failed to read the file, err:%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	w.Write(bytes)
}

func UpdateFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// 操作类型
	opType := r.Form.Get("type")
	hash := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	// 目前仅仅支持重命名，如果不是，则返回403
	if opType != "1" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//if r.Method != "POST" {
	//	w.WriteHeader(http.StatusMethodNotAllowed)
	//	return
	//}
	fileMeta := meta.GetFileMeta(hash)
	fileMeta.FileName = newFileName
	fileMeta.UploadAt = utils.FormatTime()
	meta.UpdateFileMeta(fileMeta)

	err := os.Rename(fileMeta.Location, fileMeta.FilePath+newFileName)
	if err != nil {
		fmt.Printf("Failed to rename the file, err:%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	hash := r.Form.Get("filehash")

	fileMeta := meta.GetFileMeta(hash)
	err := os.Remove(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to delete the file, err:%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 将此文件从元数据集合中删除
	meta.RemoveFileMeta(hash)

	io.WriteString(w, "OK!")
}
