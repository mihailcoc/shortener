package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

//type Config struct {
//	Server_Address string `env:"localhost"`
//	Base_Url       string `env:"http://localhost:8080"`
//}

const port = ":8080"

func main() {
	// init router
	router := mux.NewRouter()

	//os.Setenv("Server_Address", "localhost"+port)
	//os.Setenv("Base_Url", "http:/"+os.Getenv("Server_Address")+"/")
	srv := http.Server{
		//Addr:    addr,
		Handler: router,
	}

	// router handler / endpoints
	router.HandleFunc("/", handlerPost).Methods("POST")
	router.HandleFunc("/:key", handlerGet).Methods("GET")
	router.HandleFunc("/api/shorten", handlerPostAPI).Methods("POST")
	srv.ListenAndServe()
}
