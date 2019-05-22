package handler

import (
	"encoding/json"
	"fmt"
	"github.com/pibigstar/go-cloudstore/constant"
	"github.com/pibigstar/go-cloudstore/db"
	"github.com/pibigstar/go-cloudstore/utils"
	"io"
	"io/ioutil"
	"net/http"
)

// 用户注册
func UserSignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		bytes, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(bytes)
		return
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		enc_pwd := utils.Sha1([]byte(password + cont.PASSWORD_SALT))
		b := db.UserSignup(username, enc_pwd)
		if b {
			io.WriteString(w, "Success!")
		} else {
			io.WriteString(w, "Fail!")
		}
	}
}

type LoginResponse struct {
	UserName string `json:"Username"`
	Token    string `json:"Token"`
	Location string `json:"Location"`
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == http.MethodGet {
		bytes, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		io.WriteString(w, string(bytes))
	} else if r.Method == http.MethodPost {
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		enc_pwd := utils.Sha1([]byte(password + cont.PASSWORD_SALT))
		_, err := db.UserLogin(username, enc_pwd)
		if err != nil {
			io.WriteString(w, "Failed！")
			return
		}
		token := utils.GenToken(username)
		db.UpdateUserToken(username, token)
		response := LoginResponse{
			Location: fmt.Sprintf("http://%s/static/view/home.html", r.Host),
			UserName: username,
			Token:    token,
		}
		resp := utils.RespMsg{
			Code: 0,
			Msg:  "OK",
			Data: response,
		}
		bytes, _ := json.Marshal(resp)
		io.WriteString(w, string(bytes))
	}
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")

	user, err := db.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

