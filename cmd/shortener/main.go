package main

import (
	"github.com/gin-gonic/gin"
)

var urls = make(map[string]string)

var (
	addr    = "localhost:8080"
	scheme  = "http"
	baseURL = scheme + "://" + addr
)

func main() {
	server := gin.Default()
	server.GET(
		"/:key",
		main.handler.handlerGet,
	)
	server.POST("/", main.handler.handlerPost)
	server.Run(addr)
}
