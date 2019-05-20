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


	fmt.Println("server is started...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("starter server error")
	}
}
