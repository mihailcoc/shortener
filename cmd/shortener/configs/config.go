package configs

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
	"github.com/mihailcoc/shortener/internal/app/service"
)

// Создаем структуру для загрузки переменных окружения.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"storage.json"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             []byte
	Workers         int `env:"WORKERS" envDefault:"10"`
	WorkersBuffer   int `env:"WORKERS_BUFFER" envDefault:"100"`
}

func checkExists(f string) bool {
	return flag.Lookup(f) == nil
}

func NewConfig() Config {

	var cfg Config

	random, err := service.GenerateRandom(16)
	if err != nil {
		log.Fatal(err)
	}

	cfg.Key = random

	err = env.Parse(&cfg)
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
	if checkExists("d") {
		flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "DatabaseDSN")
	}
	if checkExists("w") {
		flag.IntVar(&cfg.Workers, "w", cfg.Workers, "Workers")
	}

	if checkExists("wb") {
		flag.IntVar(&cfg.WorkersBuffer, "wb", cfg.WorkersBuffer, "WorkersBuffer")
	}
	flag.Parse()

	return cfg
}
