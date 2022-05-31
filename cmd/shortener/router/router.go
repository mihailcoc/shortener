package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mihailcoc/shortener/cmd/shortener/configs"
	"github.com/mihailcoc/shortener/internal/app/handler"
)

// определяем новый роутер для запуска нового сервера
func NewRouter(repo handler.Repository, cfg configs.Config) *chi.Mux {
	h := handler.NewHandler(repo, cfg.BaseURL)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// определяем связку хендлеров и адресов
	router.Route("/", func(r chi.Router) {
		router.Post("/", h.CreateShortURL)
		router.Post("/api/shorten", h.ShortenURL)
		router.Get("/api/user/urls", h.GetUserURLs)
		router.Get("/ping", h.PingDB)
		router.Post("/api/shorten/batch", h.CreateBatch)
	})

	return router
}
