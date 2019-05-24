package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pibigstar/go-cloudstore/constant"
	"github.com/pibigstar/go-cloudstore/db"
	"github.com/pibigstar/go-cloudstore/utils"
	"net/http"
)

// 去注册页面
func ToUserSignupHandler(c *gin.Context)  {
	c.Redirect(http.StatusFound, "static/view/signup.html")
}
// 用户注册
func DoUserSignupHandler(c *gin.Context){
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	enc_pwd := utils.Sha1([]byte(password + cont.PASSWORD_SALT))
	b := db.UserSignup(username, enc_pwd)
	if b {
		c.JSON(http.StatusOK,gin.H{
			"msg": "Success",
			"code": 0,
		})
	} else {
		c.JSON(http.StatusOK,gin.H{
			"msg": "Fail",
			"code": -1,
		})
	}
}

type LoginResponse struct {
	UserName string `json:"Username"`
	Token    string `json:"Token"`
	Location string `json:"Location"`
}

// 去用户登录页面
func ToUserLoginHandler(c *gin.Context) {
	c.Redirect(http.StatusFound,"static/view/signin.html")
}
// 处理用户登录请求
func DoUserLoginHandler(c *gin.Context)  {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	enc_pwd := utils.Sha1([]byte(password + cont.PASSWORD_SALT))
	_, err := db.UserLogin(username, enc_pwd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "Login Failed!",
		})
		return
	}
	token := utils.GenToken(username)
	db.UpdateUserToken(username, token)
	response := LoginResponse{
		Location: "/static/view/home.html",
		UserName: username,
		Token:    token,
	}
	resp := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: response,
	}
	c.Data(http.StatusOK, "application/json",resp.JSONBytes())
}

func GetUserInfoHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	user, err := db.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "get user info failed!",
		})
		return
	}
	resp := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	c.JSON(http.StatusOK, resp)
}
