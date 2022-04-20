package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	SERVER_ADDRESS string `env:"localhost:8080"`
	BASE_URL       string `env:"localhost:8080"`
}

var (
	addr    = "localhost:8080"
	scheme  = "http"
	baseURL = scheme + "://" + addr
)

const port = ":8080"

func main() {
	// init router
	router := mux.NewRouter()

	//os.Setenv("ServerAddress", "localhost"+port)
	//os.Setenv("baseURL", "http:/"+os.Getenv("ServerAddress")+"/")
	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}

	// router handler / endpoints
	router.HandleFunc("/", handlerPost).Methods("POST")
	router.HandleFunc("/:key", handlerGet).Methods("GET")
	router.HandleFunc("/api/shorten", handlerPostAPI).Methods("POST")
	log.Fatal(srv.ListenAndServe())
}
