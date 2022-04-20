package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	addr    = "localhost:8080"
	scheme  = "http"
	baseURL = scheme + "://" + addr
)

//  описываем структуру JSON в запросе - {"url":"<some_url>"}
type jsonURLBody struct {
	URL string `json:"url"`
}

//  описываем структуру создаваемого JSON вида {"result":"<shorten_url>"}
type ResultURL struct {
	Result string `json:"result"`
}

func handlerPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	//log.Printf("Получено тело запроса: %s", body)
	// По ключу помещаем значение localhost map.
	mKey := randomString(len(body) / 4)

	urls[mKey] = string(body)

	response := fmt.Sprintf("%s/%s", baseURL, mKey)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
	defer r.Body.Close()
}

func handlerPostAPI(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен post application/json")
	jsonURL, err := io.ReadAll(r.Body) // считываем JSON из тела запроса
	if err != nil {
		panic(err)
	}
	log.Printf("Получено тело запроса: %s", jsonURL)

	// создаеём экземпляр структуры для заполнения из JSON
	jsonBody := jsonURLBody{}

	// парсим JSON и записываем результат в экземпляр структуры
	err = json.Unmarshal(jsonURL, &jsonBody)
	if err != nil {
		panic(err)
	}
	// По ключу помещаем значение localhost map.
	mKey := randomString(len(jsonURL) / 4)
	log.Printf("Получен mKey: %s", mKey)

	urls[mKey] = string(jsonURL)
	shortURL := fmt.Sprintf("%s/%s", baseURL, mKey)
	log.Printf("Получен shortURL: %s", shortURL)

	// создаем экземпляр структуры и вставляем в него короткий URL для отправки в JSON
	resultURL := ResultURL{
		Result: shortURL,
	}
	// изготавливаем JSON
	shortJSONURL, err := json.Marshal(resultURL)
	if err != nil {
		panic(err)
	}
	log.Printf("Получен shortJSONURL: %s", shortJSONURL)
	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortJSONURL))
	defer r.Body.Close()
}

func handlerGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен get default")
	//key := r.Param("key")
	key := strings.TrimPrefix(r.URL.Path, "/")
	log.Printf("Получен key %s", key)
	if url, ok := urls[key]; ok {
		log.Printf("Отдаем url %s", url)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Location", url)
		//	w.Header.Set("Location", url)
		//g.Header("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
		//g.Redirect(http.StatusTemporaryRedirect, url)
		//return
		defer r.Body.Close()
	}
}
