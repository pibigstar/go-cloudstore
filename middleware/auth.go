package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pibigstar/go-cloudstore/utils"
	"net/http"
)

// http请求拦截器
func HttpInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")
		if !checkToken(username, token) {
			//终止此请求链
			c.Abort()
			resp := utils.NewRespMsg(http.StatusForbidden, "token无效", nil)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.Next()
	}
}

func checkToken(username, token string) bool {
	genToken := utils.GenToken(username)
	if genToken == token {
		return true
	} else {
		return false
	}
}
