package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	addr   string
	config Config
}

func NewServer(addr string, config Config) *server {
	return &server{
		addr:   addr,
		config: config,
	}
}

func (s *server) Start() {
	h := NewHandler(s.config)

	router := mux.NewRouter()

	router.HandleFunc("/{url}", h.handlerGet).Methods("GET")
	router.HandleFunc("/", h.handlerPost).Methods("POST")
	router.HandleFunc("/api/shorten", h.handlerPostAPI).Methods("POST")

	log.Fatal(http.ListenAndServe(s.addr, handlers.GzipHandle(router)))
}
