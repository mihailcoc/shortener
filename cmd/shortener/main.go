package main

import (
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
	baseURL       = scheme + "://" + ServerAddress
)

const port = ":8080"

func main() {
	// init router
	router := mux.NewRouter()

	os.Setenv("ServerAddress", "localhost"+port)
	os.Setenv("baseURL", "http:/"+os.Getenv("ServerAddress")+"/")
	srv := http.Server{
		Addr:    baseURL,
		Handler: router,
	}

	// router handler / endpoints
	router.HandleFunc("/", handlerPost).Methods("POST")
	router.HandleFunc("/:key", handlerGet).Methods("GET")
	router.HandleFunc("/api/shorten", handlerPostAPI).Methods("POST")
	log.Fatal(srv.ListenAndServe())
}
