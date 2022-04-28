package main

func main() {
	// Загружаем переменные окружения.
	c := NewConfig()
	// Формирует структуру для старта сервера
	serv := NewServer(c.ServerAddress, c)

	serv.StartServer()
}
