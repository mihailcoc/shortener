package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

var m = make(map[string]string)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	// m := make(map[string]string)

	switch r.Method {
	// если методом POST
	case "POST":
		b, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		//fmt.Println(string(b))
		//fmt.Println("r.Body", string(b))
		//u, err := url.Parse(string(b))
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(u)
		//fmt.Println("Path:", u.Path)
		//fmt.Println("len(Path):", len(u.Path))
		//fmt.Println("randomString", randomString(len(u.Path)/4))
		// Генерируем ключ
		mKey := randomString(len(b) / 4)
		// По ключу проверяем наличие в map.
		if intid, ok := strconv.Atoi(m[mKey]); ok != nil {
			fmt.Println("Значение в map уже задано:", strconv.Itoa(intid))
		}
		// По ключу помещаем значение localhost map.
		m[mKey] = string(b)
		// Генерируем ответ
		responseURL := "http://" + r.Host + r.URL.String() + mKey
		//fmt.Println(responseURL)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		//w.Header().Set("Location", responseURL)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(responseURL))
		//fmt.Fprint(w)
	// если методом GET
	case "GET":
		// извлекаем фрагмент id из URL запроса GET /{id}
		//qq := r.URL.Path
		q := strings.TrimPrefix(r.URL.Path, "/")
		// fmt.Println("q", q)
		if q == "" {
			http.Error(w, "The query parameter is missing", http.StatusBadRequest)
			return
		}
		// достаем из map оригинальный URL
		//origURL, exists := m[q]
		//if exists {
		//	http.Error(w, "The query parameter is missing", http.StatusBadRequest)
		//	return
		//}
		origURL := m[q]
		fmt.Println("65 origURL ", origURL)
		// устанавливаем в заголовке оригинальный URL
		w.Header().Set("Location", origURL)
		// устанавливаем статус-код 307
		w.WriteHeader(http.StatusTemporaryRedirect)
		// отдаем редирект на собственный url и код 307
		fmt.Fprint(w)
	default:
		http.Error(w, "Only GET and POST requests are allowed!", http.StatusMethodNotAllowed)

	}
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(97, 122))
	}
	return string(bytes)
}

func main() {
	http.HandleFunc("/", viewHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
