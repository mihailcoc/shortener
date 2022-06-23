package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mihailcoc/shortener/cmd/shortener/configs"
	"github.com/mihailcoc/shortener/internal/app/handler"
	"github.com/mihailcoc/shortener/internal/app/workers"
)

func NewRouter(repo handler.Repository, cfg configs.Config, wp *workers.WorkerPool) *chi.Mux {
	h := handler.NewHandler(repo, cfg.BaseURL, wp)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/", func(r chi.Router) {
		router.Post("/", h.CreateShortURL)
		router.Get("/{id}", h.RetrieveShortURL)
		router.Post("/api/shorten", h.ShortenURL)
		router.Get("/api/user/urls", h.GetUserURLs)
		router.Get("/ping", h.PingDB)
		router.Post("/api/shorten/batch", h.CreateBatch)
		router.Delete("/api/user/urls", h.DeleteBatch)
	})

	return router
}
