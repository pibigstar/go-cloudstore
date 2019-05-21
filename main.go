package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pibigstar/go-cloudstore/handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/success", handler.UploadSuccessHandler)
	http.HandleFunc("/file/meta", handler.GetFileMeta)
	http.HandleFunc("/file/query", handler.QueryFileHandler)
	http.HandleFunc("/file/download", handler.DownloadFileHandler)
	http.HandleFunc("/file/update", handler.UpdateFileMetaHandler)
	http.HandleFunc("/file/delete", handler.DeleteFileHandler)
	http.HandleFunc("/user/signup", handler.UserSignupHandler)
	http.HandleFunc("/user/signin", handler.UserLoginHandler)
	http.HandleFunc("/user/info", handler.GetUserInfoHandler)
	http.HandleFunc("/home", handler.GoHomeHandler)

	// 静态资源配置
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("server is started...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("starter server error")
	}
}
