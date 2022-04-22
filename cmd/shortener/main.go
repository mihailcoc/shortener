package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Config struct {
	ServerAddress int    `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"localhost:8080"`
}

var (
	//port          = ":8080"
	ServerAddress = "localhost" + ":8080"
	scheme        = "http"
	BaseURL       = scheme + "://" + ServerAddress
)

func main() {
	os.Setenv("SERVER_ADDRESS", ":8080")
	os.Getenv("SERVER_ADDRESS")
	log.Printf("Getenv ServerAddress")
	log.Println(ServerAddress)
	log.Printf("ServerAddress")
	log.Println(ServerAddress)
	// 1 вариант
	//var cfg Config
	//err := env.Parse(&cfg)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("cfg.ServerAddress: %s", cfg.ServerAddress)
	//*ServerAddress := cfg.ServerAddress
	// 1 вариант

	//os.Setenv("port", ":8080")
	//os.Setenv("SERVER_ADDRESS", ":8080")
	//os.Setenv("ServerAddress", "localhost"+os.Getenv("port"))
	//os.Setenv("BaseURL", "http:/localhost"+os.Getenv("ServerAddress")+"/")

	//3 вариант
	//pflag.StringVarP(*ServerAddress, port)
	//3 вариант

	//2 вариант
	ServerAddress := flag.String(ServerAddress, "127.0.0.1:8000", "SERVER_ADDRESS - адрес для запуска HTTP-сервера")
	flag.Parse()

	if u, f := os.LookupEnv("SERVER_ADDRESS"); f {
		*ServerAddress = u
	}
	flag.Parse()
	log.Printf("*ServerAddress перед сервером")
	log.Println(*ServerAddress)
	//2 вариант

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
