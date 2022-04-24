package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

type server struct {
	addr   string
	config Config
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

	flag.Parse()

	return cfg
}

func NewServer(addr string, config Config) *server {
	return &server{
		addr:   addr,
		config: config,
	}
}

func (s *server) Start() {
	h := main.NewServer(s.config)

	router := mux.NewRouter()

	router.HandleFunc("/{url}", handlerGet).Methods("GET")
	router.HandleFunc("/", handlerPost).Methods("POST")
	router.HandleFunc("/api/shorten", handlerPostAPI).Methods("POST")

	log.Fatal(http.ListenAndServe(s.addr, handlers.GzipHandle(router)))
}
