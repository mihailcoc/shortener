package servers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mihailcoc/shortener/cmd/shortener/configs"
	"github.com/mihailcoc/shortener/internal/app/compressor"
	"github.com/mihailcoc/shortener/internal/app/handler"
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
	router.HandleFunc("/api/user/urls", h.GetUserURLs).Methods("GET")

	log.Printf("Сервер запущен")
	//log.Printf("FileStoragePath перед запуском сервера %s", h.config.FileStoragePath)
	log.Fatal(http.ListenAndServe(s.addr, compressor.GzipHandle(router)))

}
