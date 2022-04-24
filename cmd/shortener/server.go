package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mihailcoc/shortener/configs"
)

type server struct {
	addr   string
	config configs.Config
}

func New(addr string, config configs.Config) *server {
	return &server{
		addr:   addr,
		config: config,
	}
}

func (s *server) Start() {
	h := handlers.New(s.config)

	router := mux.NewRouter()

	//router.Use(mux.middleware.Logger)
	//router.Use(mux.middleware.Recoverer)

	router.Route("/", func(r mux.Router) {
		router.Get("/{id}", h.Get)
		router.Get("/", h.Get)
		router.Post("/", h.Save)
		router.Post("/api/shorten", h.SaveJSON)
	})

	log.Fatal(http.ListenAndServe(s.addr, handlers.GzipHandle(router)))
}
