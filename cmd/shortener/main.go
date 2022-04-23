package main

import (
	"log"
	"net/http"
	"os"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL       string `env:"http://localhost:8080"`
}

var (
	ServerAddress = "127.0.0.1:8080"
	scheme        = "http://"
	BaseURL       = scheme + ServerAddress
)

func main() {
	// 1 вариант
	log.Printf("ServerAddress вначале:%s", ServerAddress)
	if u, f := os.LookupEnv("SERVER_ADDRESS"); f {
		ServerAddress = u
	}
	log.Printf("ServerAddress после LookupEnv:%s", ServerAddress)

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("cfg.ServerAddress: %s", cfg.ServerAddress)
	ServerAddress := cfg.ServerAddress
	log.Printf("ServerAddress после env.Parse:%s", ServerAddress)
	// 1 вариант

	//os.Setenv("SERVER_ADDRESS", ":8080")

	//2 вариант
	//os.Setenv("SERVER_ADDRESS", ":8080")
	// попробовать
	//os.Getenv("ServerAddress")
	//log.Printf("ServerAddress")
	//log.Println(ServerAddress)
	//if u, f := os.LookupEnv("SERVER_ADDRESS"); f {
	//	ServerAddress = u
	//}
	//log.Printf("ServerAddress после LookupEnv")
	//log.Println(ServerAddress)
	//ServerAddress := flag.String(ServerAddress, ":8000", "SERVER_ADDRESS - адрес для запуска HTTP-сервера")
	//flag.Parse()
	//log.Printf("*ServerAddress после flag.String")
	//log.Println(*ServerAddress)
	//log.Printf("*ServerAddress перед сервером")
	//log.Println(*ServerAddress)
	//2 вариант

	//3 вариант
	//os.Setenv("SERVER_ADDRESS", "127.0.0.1:8080")
	//os.Getenv("SERVER_ADDRESS")

	//if u, f := os.LookupEnv("SERVER_ADDRESS"); f {
	//	ServerAddress = u
	//}
	//log.Printf("ServerAddress после LookupEnv")
	//log.Println(ServerAddress)
	//pflag.StringVarP(&ServerAddress, "SERVER_ADDRESS", "s", "127.0.0.1:8000", "SERVER_ADDRESS - адрес для запуска HTTP-сервера")
	//pflag.Parse()
	//log.Printf("ServerAddress после pflag.Parse")
	//log.Println(ServerAddress)
	//log.Printf("ServerAddress перед сервером")
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
