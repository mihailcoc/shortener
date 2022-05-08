package main

//"github.com/mihailcoc/shortener/cmd/shortener/handler"
//"github.com/mihailcoc/shortener/internal/app/server"
//"github.com/mihailcoc/shortener/internal/app/storage"

func main() {
	// Загружаем переменные окружения.
	c := NewConfig()
	// Формирует структуру для старта сервера
	serv := NewServer(c.ServerAddress, c)

	serv.StartServer()
}
