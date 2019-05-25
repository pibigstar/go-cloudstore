package route

import (
	"github.com/gin-gonic/gin"
	"github.com/pibigstar/go-cloudstore/handler"
	"github.com/pibigstar/go-cloudstore/middleware"
)

func Router() *gin.Engine {

	router := gin.Default()
	// 处理静态资源
	router.Static("/static/", "./static")

	//无需验证
	router.GET("/user/signup", handler.ToUserSignupHandler)
	router.POST("/user/signup", handler.DoUserSignupHandler)
	router.GET("user/signin", handler.ToUserLoginHandler)
	router.POST("user/signin", handler.DoUserLoginHandler)

	//添加拦截器, Use之后的所有Handler都会经过拦截器校验
	router.Use(middleware.HttpInterceptor())
	router.POST("user/info", handler.GetUserInfoHandler)

	router.GET("/home", handler.GoHomeHandler)
	//上传相关接口
	router.GET("/file/upload", handler.ToUploadHandler)
	router.POST("/file/upload", handler.DoUploadHandler)
	//分块上传
	router.POST("/file/mpupload/init", handler.InitialMultipartUploadHandler)
	router.POST("/file/mpupload/uppart", handler.UploadPartHandler)
	router.POST("/file/mpupload/complete", handler.CompleteUploadHandler)

	//秒传
	router.POST("/file/fastupload", handler.TryFastUploadHandler)
	//查询
	router.POST("/file/meta", handler.GetFileMeta)
	router.POST("/file/query", handler.QueryFileHandler)
	//下载
	router.GET("/file/download", handler.DownloadFileHandler)
	//更新与删除
	router.POST("/file/update", handler.UpdateFileMetaHandler)
	router.POST("/file/delete", handler.DeleteFileHandler)

	return router
}
