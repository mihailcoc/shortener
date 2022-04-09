package main

import (
	"github.com/gin-gonic/gin"
	"main/cmd/shortener/handler"
)

var urls = make(map[string]string)

var (
	addr    = "localhost:8080"
	scheme  = "http"
	baseURL = scheme + "://" + addr
)

func main() {
	server := gin.Default()
	server.GET("/:key", handler.handlerGet)
	server.POST("/", handler.handlerPost)
	server.Run(addr)
}
