package main

import (
	"github.com/pibigstar/go-cloudstore/gateway/route"
)

// 启动gateway路由，提供User路由服务
func main() {
	r := route.Router()
	r.Run(":8080")
}
