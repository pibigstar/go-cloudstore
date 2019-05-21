package handler

import (
	"io"
	"io/ioutil"
	"net/http"
)

func GoHomeHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadFile("./static/view/home.html")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.WriteString(w, string(bytes))
}
