package main

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	ServerAddress string `env:"localhost:8080"`
	BaseURL       string `env:"localhost:8080"`
}

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
