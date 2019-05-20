package handler

import (
	"github.com/pibigstar/go-cloudstore/constant"
	"github.com/pibigstar/go-cloudstore/db"
	"github.com/pibigstar/go-cloudstore/utils"
	"io"
	"io/ioutil"
	"net/http"
)
// 用户注册
func UserSignupHandler(w http.ResponseWriter,r *http.Request)  {
	if r.Method == http.MethodGet {
		bytes, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(bytes)
		return
	}else if r.Method == http.MethodPost {
		r.ParseForm()

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		enc_pwd := utils.Sha1([]byte(password+cont.SALT))
		b := db.UserSignup(username, enc_pwd)
		if b {
			io.WriteString(w,"Success!")
		}else {
			io.WriteString(w,"Fail!")
		}
	}
}
