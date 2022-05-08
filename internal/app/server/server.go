package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mihailcoc/shortener/cmd/shortener/config"
	"github.com/mihailcoc/shortener/internal/app/compressor"
	"github.com/mihailcoc/shortener/internal/app/handler"
)

type server struct {
	addr   string
	config config.Config
}

func NewServer(addr string, config config.Config) *server {
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

	router.HandleFunc("/{url}", h.handlerGet).Methods("GET")
	router.HandleFunc("/", h.handlerPost).Methods("POST")
	router.HandleFunc("/api/shorten", h.handlerPostAPI).Methods("POST")
	log.Printf("ServerAddress перед запуском сервера %s", h.config.ServerAddress)
	log.Printf("FileStoragePath перед запуском сервера %s", h.config.FileStoragePath)
	log.Fatal(http.ListenAndServe(s.addr, compressor.GzipHandle(router)))

}
