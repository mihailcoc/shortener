package main

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	ServerAddress string `env:"localhost:8080"`
	BaseURL       string `env:"localhost:8080"`
}

func main() {
	server := gin.Default()
	server.GET(
		"/:key",
		handlerGet,
	)
	server.GET(
		"/api/shorten:key",
		handlerGetAPI,
	)
	server.POST("/", handlerPost)
	server.POST("/api/shorten", handlerPostAPI)
	server.Run(addr)
}
