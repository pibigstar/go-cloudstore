package main

import (
	"fmt"
	"github.com/pibigstar/go-cloudstore/route"
	"net/http"

	"github.com/pibigstar/go-cloudstore/handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/fastupload", handler.TryFastUploadHandler)
	http.HandleFunc("/file/upload/success", handler.UploadSuccessHandler)
	http.HandleFunc("/file/meta", handler.GetFileMeta)
	http.HandleFunc("/file/query", handler.QueryFileHandler)
	http.HandleFunc("/file/download", handler.DownloadFileHandler)
	http.HandleFunc("/file/downloadurl", handler.DownloadFileByUrlHandler)
	http.HandleFunc("/file/update", handler.UpdateFileMetaHandler)
	http.HandleFunc("/file/delete", handler.DeleteFileHandler)

	// 分块上传
	http.HandleFunc("/file/mpupload/init", handler.InitialMultipartUploadHandler)
	http.HandleFunc("/file/mpupload/uppart", handler.UploadPartHandler)
	http.HandleFunc("/file/mpupload/complete", handler.CompleteUploadHandler)

	// 用户相关接口
	http.HandleFunc("/home", handler.GoHomeHandler)

	// 静态资源配置
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("server is started...")


	router := route.Router()
	router.Run(":8080")
}
