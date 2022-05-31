package configs

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

// Создаем структуру для загрузки переменных окружения.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"storage.json"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             []byte
}

func checkExists(f string) bool {
	return flag.Lookup(f) == nil
}

// Задаём функцию новых конфигураций
func NewConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Задаем флаг для b переменной окружения URL
	if checkExists("b") {
		flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "BaseUrl")
	}
	// Задаем флаг для a переменной окружения ServerAddress
	if checkExists("a") {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "ServerAddress")
	}
	// Задаем флаг для f переменной окружения FileStoragePath
	if checkExists("f") {
		flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "FileStoragePath")
	}
	// Парсим флаги
	flag.Parse()

	return cfg
}
