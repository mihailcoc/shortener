package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	ServerAddress string `env:"localhost:8080"`
	BaseURL       string `env:"localhost:8080"`
}

var (
	SERVER_ADDRESS = "localhost:8080"
	scheme         = "http"
	BASE_URL       = scheme + "://" + SERVER_ADDRESS
)

const port = ":8080"

func main() {
	// init router
	router := mux.NewRouter()

	//os.Setenv("Server_Address", "localhost"+port)
	//os.Setenv("Base_Url", "http:/"+os.Getenv("Server_Address")+"/")
	srv := http.Server{
		Addr:    SERVER_ADDRESS,
		Handler: router,
	}

	// router handler / endpoints
	router.HandleFunc("/{url}", handlerGet).Methods("GET")
	router.HandleFunc("/", handlerPost).Methods("POST")
	router.HandleFunc("/api/shorten", handlerPostAPI).Methods("POST")
	log.Fatal(srv.ListenAndServe())
}
