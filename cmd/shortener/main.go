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
	port          = ":8080"
	ServerAddress = "localhost" + port
	scheme        = "http"
	BaseURL       = scheme + "://" + ServerAddress
)

func main() {
	os.Setenv("port", ":8080")
	os.Setenv("ServerAddress", "localhost"+os.Getenv("port"))
	os.Setenv("BaseURL", "http:/"+os.Getenv("ServerAddress")+"/")
	ServerAddress := flag.String("a", "127.0.0.1:8000", "SERVER_ADDRESS - адрес для запуска HTTP-сервера")
	if u, f := os.LookupEnv("ServerAddress"); f {
		*ServerAddress = u
	}
	flag.Parse()
	port := flag.String("b", ":8000", "PORT - порт для запуска HTTP-сервера")
	if uport, f := os.LookupEnv("port"); f {
		*port = uport
	}
	flag.Parse()
	//log.Printf("Распарсили flag: %s", flag.Parse())
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
