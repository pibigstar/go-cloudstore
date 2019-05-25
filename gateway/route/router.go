package route

import (
	"github.com/gin-gonic/gin"
	"github.com/pibigstar/go-cloudstore/gateway/gateway"
)

// 网关API路由
func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static", "./static")

	router.GET("/user/signup", gateway.ToUserSignupHandler)
	router.POST("/user/signup", gateway.DoUserSignupHandler)

	return router
}
