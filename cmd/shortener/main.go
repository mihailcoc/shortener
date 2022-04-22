package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"localhost:8080"`
	Host          string `env:"localhost:8080"`
}

var (
	port          = ":8080"
	ServerAddress = "localhost" + port
	scheme        = "http"
	BaseURL       = scheme + "://" + ServerAddress
)

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("cfg.ServerAddress: %s", cfg.ServerAddress)
	//os.Setenv("port", ":8080")
	os.Setenv("SERVER_ADDRESS", ":8080")
	//os.Setenv("ServerAddress", "localhost"+os.Getenv("port"))
	os.Setenv("BaseURL", "http:/localhost"+os.Getenv("ServerAddress")+"/")

	ServerAddress := flag.String("a", "127.0.0.1:8000", "SERVER_ADDRESS - адрес для запуска HTTP-сервера")
	flag.Parse()
	//pflag.StringVarP(*ServerAddress, port)
	log.Printf("1 &ServerAddress:")
	log.Println(ServerAddress)
	//flag.StringVar(&config.ServerAddress, "a", "127.0.0.1:8080", "server address")
	//flag.Parse()
	//log.Printf("ServerAddress: %s", ServerAddress)

	//ServerAddress, exist := os.LookupEnv("ADDRESS")
	//if exist {
	//	config.ServerAddress = ServerAddress
	//}
	//log.Printf("ServerAddress: %s", ServerAddress)

	if u, f := os.LookupEnv("SERVER_ADDRESS"); f {
		*ServerAddress = u
	}
	flag.Parse()
	log.Printf("2 &SERVER_ADDRESS:")
	log.Println(ServerAddress)

	//port := flag.String("b", ":8000", "PORT - порт для запуска HTTP-сервера")
	//if uport, f := os.LookupEnv("port"); f {
	///	*port = uport
	//}
	//flag.Parse()
	//log.Printf("port:")
	//log.Println(port)
	// init router
	router := mux.NewRouter()

	srv := http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	// router handler / endpoints
	router.HandleFunc("/{url}", handlerGet).Methods("GET")
	router.HandleFunc("/", handlerPost).Methods("POST")
	router.HandleFunc("/api/shorten", handlerPostAPI).Methods("POST")
	log.Fatal(srv.ListenAndServe())
}
