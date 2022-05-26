package main

import (
	"github.com/mihailcoc/shortener/cmd/shortener/configs"
	"github.com/mihailcoc/shortener/internal/app/servers"
)

func main() {
	// Загружаем переменные окружения.
	c := configs.NewConfig()
	// Формирует структуру для старта сервера
	serv := servers.NewServer(c.ServerAddress, c)

	serv.StartServer()
}
