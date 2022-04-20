package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

//type Config struct {
//	Server_Address string `env:"localhost"`
//	Base_Url       string `env:"http://localhost:8080"`
//}

var (
	addr    = "localhost:8080"
	scheme  = "http"
	baseURL = scheme + "://" + addr
)

//const port = ":8080"

func main() {
	// init router
	router := mux.NewRouter()
	//router := http.NewRouter()

	//os.Setenv("ServerAddress", "localhost"+port)
	//os.Setenv("baseURL", "http:/"+os.Getenv("ServerAddress")+"/")
	//srv := http.Server{
	//	//Addr:    addr,
	//	Handler: router,
	//}

	// router handler / endpoints
	//router.HandleFunc("/", handlerPost).Methods("POST")
	router.HandleFunc("/", handlerPost).Methods("POST")
	router.HandleFunc("/:key", handlerGet).Methods("GET")
	router.HandleFunc("/api/shorten", handlerPostAPI).Methods("POST")
	http.ListenAndServe(":8080", nil)
}
