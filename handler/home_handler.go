package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GoHomeHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "static/view/home.html")
}

//bytes, err := ioutil.ReadFile("./static/view/home.html")
//if err != nil {
//w.WriteHeader(http.StatusNotFound)
//return
//}
//io.WriteString(w, string(bytes))
