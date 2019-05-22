package middleware

import (
	"github.com/pibigstar/go-cloudstore/utils"
	"net/http"
)

// http请求拦截器
func HttpInterceptor(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		username := r.Form.Get("username")
		token := r.Form.Get("token")
		if !checkToken(username, token) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		f(w, r)
	})
}

func checkToken(username, token string) bool {
	genToken := utils.GenToken(username)
	if genToken == token {
		return true
	} else {
		return false
	}
}
