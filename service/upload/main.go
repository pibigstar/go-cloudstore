package main

import (
	"fmt"
	"github.com/pibigstar/go-cloudstore/middleware"
	"log"
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
	http.HandleFunc("/user/signup", handler.UserSignupHandler)
	http.HandleFunc("/user/signin", handler.UserLoginHandler)
	// 添加了拦截器
	http.HandleFunc("/user/info", middleware.HttpInterceptor(handler.GetUserInfoHandler))
	http.HandleFunc("/home", handler.GoHomeHandler)

	// 静态资源配置
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("server is started...")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("starter server error")
	}
}
