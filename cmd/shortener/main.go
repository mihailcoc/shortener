package main

import (
	"github.com/mihailcoc/shortener/cmd/shortener/config"
)

func main() {
	// Загружаем переменные окружения.
	c := config.NewConfig()
	// Формирует структуру для старта сервера
	serv := NewServer(c.ServerAddress, c)

	serv.StartServer()
}
