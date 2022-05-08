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
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func checkExists(f string) bool {
	return flag.Lookup(f) == nil
}

func NewConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if checkExists("b") {
		flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "BaseUrl")
	}

	if checkExists("a") {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "ServerAddress")
	}
	if checkExists("f") {
		flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "FileStoragePath")
	}
	flag.Parse()

	return cfg
}