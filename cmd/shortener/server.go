package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mihailcoc/shortener/cmd/shortener/compressor"
	"github.com/mihailcoc/shortener/cmd/shortener/configs"
	"github.com/mihailcoc/shortener/cmd/shortener/handler"
)

type server struct {
	addr   string
	config configs.Config
}

func NewServer(addr string, config configs.Config) *server {
	return &server{
		addr:   addr,
		config: config,
	}
}

func (s *server) StartServer() {
	// Создаем новый handler с переменными окружения.
	h := handler.NewHandler(s.config)
	// Создаем новый роутер.
	router := mux.NewRouter()

	router.HandleFunc("/{url}", h.HandlerGet).Methods("GET")
	router.HandleFunc("/", h.HandlerPost).Methods("POST")
	router.HandleFunc("/api/shorten", h.HandlerPostAPI).Methods("POST")
	//log.Printf("ServerAddress перед запуском сервера %s", h.configs.ServerAddress)
	//log.Printf("FileStoragePath перед запуском сервера %s", h.configs.FileStoragePath)
	log.Fatal(http.ListenAndServe(s.addr, compressor.GzipHandle(router)))

}
