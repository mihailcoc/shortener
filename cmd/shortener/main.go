package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Config struct {
	ServerAddress string `env:"localhost:8080"`
	BaseURL       string `env:"localhost:8080"`
}

var (
	ServerAddress = "localhost:8080"
	scheme        = "http"
	BaseURL       = scheme + "://" + ServerAddress
)

var port = ":8080"

func main() {

	os.Setenv("ServerAddress", "localhost"+port)
	os.Setenv("BaseURL", "http:/"+os.Getenv("ServerAddress")+"/")
	ServerAddress := flag.String("a", "127.0.0.1:8000", "SERVER_ADDRESS - адрес для запуска HTTP-сервера")
	if u, f := os.LookupEnv("ServerAddress"); f {
		*ServerAddress = u
	}
	// init router
	router := mux.NewRouter()

	srv := http.Server{
		Addr:    *ServerAddress,
		Handler: router,
	}

	// router handler / endpoints
	router.HandleFunc("/{url}", handlerGet).Methods("GET")
	router.HandleFunc("/", handlerPost).Methods("POST")
	router.HandleFunc("/api/shorten", handlerPostAPI).Methods("POST")
	log.Fatal(srv.ListenAndServe())
}
