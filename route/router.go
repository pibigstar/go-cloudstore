package route

import (
	"github.com/gin-gonic/gin"
	"github.com/pibigstar/go-cloudstore/handler"
	"github.com/pibigstar/go-cloudstore/middleware"
)

func Router() *gin.Engine  {

	router := gin.Default()
	// 处理静态资源
	router.Static("/static/","./static")

	//无需验证
	router.GET("/user/signup", handler.ToUserSignupHandler)
	router.POST("/user/signup", handler.DoUserSignupHandler)
	router.GET("user/signin", handler.ToUserLoginHandler)
	router.POST("user/signin", handler.DoUserLoginHandler)

	//添加拦截器, Use之后的所有Handler都会经过拦截器校验
	router.Use(middleware.HttpInterceptor())
	router.POST("user/info", handler.GetUserInfoHandler)

	return router
}