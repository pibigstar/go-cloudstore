package main

import (
	"fmt"
	"github.com/pibigstar/go-cloudstore/route"
)

func main() {

	fmt.Println("server is started...")

	router := route.Router()
	router.Run(":8080")
}
