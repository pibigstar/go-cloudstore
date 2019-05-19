package handler

import (
	"encoding/json"
	"fmt"
	"github.com/pibigstar/go-cloudstore/meta"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/pibigstar/go-cloudstore/utils"
)

// 文件上传后的存放路径
const uploadFilePath = "temp/"

// 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// GET请求是去上传页面
	if r.Method == "GET" {
		bytes, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(bytes))
	} else if r.Method == "POST" {
		//POST请求是上传文件
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get file data, err:%s\n", err.Error())
			return
		}
		defer file.Close()
		exist, err := utils.PathExists(uploadFilePath)
		if !exist {
			err := os.Mkdir(uploadFilePath, os.ModePerm)
			if err != nil {
				fmt.Printf("Failed to create file dir, err:%s\n", err.Error())
				return
			}
		}
		filePath := fmt.Sprintf("%s%s", uploadFilePath, header.Filename)
		// 文件元数据
		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			FilePath: uploadFilePath,
			Location: filePath,
			UploadAt: utils.FormatTime(),
		}
		// 创建一个新文件
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create new file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()
		// 将上传的文件内容复制到新文件中
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
			return
		}
		// 计算文件的sha1值
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = utils.FileSha1(newFile)
		fmt.Println("sha1:",fileMeta.FileSha1)
		meta.UpdateFileMeta(fileMeta)

		// 重定向路由
		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
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
	fileMeta := meta.GetFileMeta(hash)
	bytes, err := json.Marshal(fileMeta)
	if err != nil {
		fmt.Printf("Failed to parse fileMeta,err:%s\n",err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(bytes)
}

// 批量查询文件元数据信息
func QueryFileHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	limit, _ := strconv.Atoi(r.Form.Get("limit"))
	fileMeta := meta.GetLastFileMeta(limit)
	bytes, err := json.Marshal(fileMeta)
	if err != nil {
		fmt.Printf("Failed to parse fileMeta,err:%s\n",err.Error())
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
		fmt.Printf("Failed to open the file, err:%s\n",err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := ioutil.ReadAll(file)
	if err!=nil{
		fmt.Printf("Failed to read the file, err:%s\n",err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("Content-Disposition","attachment;filename=\""+fileMeta.FileName+"\"")
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
		fmt.Printf("Failed to rename the file, err:%s\n",err.Error())
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

func DeleteFileHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	hash := r.Form.Get("filehash")

	fileMeta := meta.GetFileMeta(hash)
	err := os.Remove(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to delete the file, err:%s\n",err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 将此文件从元数据集合中删除
	meta.RemoveFileMeta(hash)

	io.WriteString(w,"OK!")
}