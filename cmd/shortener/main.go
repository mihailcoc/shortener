package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	// если методом POST
	case "POST":
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		b, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		//length := len(b) / 2
		//UnicStr = b[:length]
		responseURL := "https://" + r.Host + r.URL.String() + string(b)
		w.Write([]byte(responseURL))
		fmt.Fprint(w)
	// если методом GET
	case "GET":
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
		// устанавливаем в заголовке оригинальный URL
		w.Header().Set("Location", "id")
		// устанавливаем статус-код 307
		w.WriteHeader(http.StatusTemporaryRedirect)
		fmt.Fprint(w)
	default:
		http.Error(w, "Only GET and POST requests are allowed!", http.StatusMethodNotAllowed)

	}
}

func main() {
	http.HandleFunc("/", viewHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
