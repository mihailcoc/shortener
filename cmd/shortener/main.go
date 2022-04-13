package main

import (
	"github.com/gin-gonic/gin"
)

var (
	addr    = "localhost:8080"
	scheme  = "http"
	baseURL = scheme + "://" + addr
)

func main() {
	server := gin.Default()
	server.GET(
		"/:key",
		handlerGet,
	)
	server.POST("/", handlerPost)
	server.POST("/api/shorten", handlerPostAPI)
	server.Run(addr)
}
