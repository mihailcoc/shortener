package main

import (
	"log"
	"net/http"
	"os"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"localhost:8080"`
}

var (
	//port          = ":8080"
	ServerAddress = "localhost" + ":8080"
	scheme        = "http"
	BaseURL       = scheme + "://" + ServerAddress
)

func main() {
	// 1 вариант
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		os.Setenv("SERVER_ADDRESS", ":8080")
		ServerAddress = os.Getenv("ServerAddress")
	}
	{
		log.Printf("cfg.ServerAddress: %s", cfg.ServerAddress)
		ServerAddress := cfg.ServerAddress
		log.Println(ServerAddress)
	}

	// 1 вариант

	//os.Setenv("port", ":8080")
	//os.Setenv("SERVER_ADDRESS", ":8080")
	//os.Setenv("ServerAddress", "localhost"+os.Getenv("port"))
	//os.Setenv("BaseURL", "http:/localhost"+os.Getenv("ServerAddress")+"/")

	//2 вариант
	//os.Setenv("SERVER_ADDRESS", ":8080")
	//log.Printf("Getenv ServerAddress")
	//log.Println(ServerAddress)
	//ServerAddress := flag.String("a", "127.0.0.1:8000", "SERVER_ADDRESS - адрес для запуска HTTP-сервера")
	//flag.Parse()
	//log.Println(*ServerAddress)
	//if u, f := os.LookupEnv("SERVER_ADDRESS"); f {
	//	*ServerAddress = u
	//}
	//log.Printf("*ServerAddress перед сервером")
	//log.Println(*ServerAddress)
	//2 вариант

	//3 вариант
	//os.Setenv("SERVER_ADDRESS", ":8080")
	//ServerAddress := os.Getenv("SERVER_ADDRESS")
	//pflag.StringVarP(&ServerAddress, "SERVER_ADDRESS", "s", "127.0.0.1:8000", "SERVER_ADDRESS - адрес для запуска HTTP-сервера")
	//pflag.Parse()
	//log.Println(ServerAddress)
	//3 вариант

	router := mux.NewRouter()

	srv := http.Server{
		Addr:    ServerAddress,
		Handler: router,
	}

	// router handler / endpoints
	router.HandleFunc("/{url}", handlerGet).Methods("GET")
	router.HandleFunc("/", handlerPost).Methods("POST")
	router.HandleFunc("/api/shorten", handlerPostAPI).Methods("POST")
	log.Fatal(srv.ListenAndServe())
}
