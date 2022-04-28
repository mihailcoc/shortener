package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

var (
	urls = make(map[string]string)
)

// тег, значение которого нужно получить
// имя поля, о котором нужно получить информацию.
const (
	targetField = "URL"
	targetTag   = "json"
)

//  описываем структуру Handler в запросе на получение переменных окружения
type Handler struct {
	storage Repository
	config  Config
}

//  описываем структуру JSON в запросе - {"url":"<some_url>"}
type jsonURLBody struct {
	URL string `json:"url"`
}

//  описываем структуру создаваемого JSON вида {"result":"<shorten_url>"}
type ResultURL struct {
	Result string `json:"result"`
}

func NewHandler(c Config) *Handler {
	h := &Handler{
		storage: NewStorage(),
		config:  c,
	}

	if err := h.storage.Load(h.config); err != nil {
		panic(err)
	}

	return h
}

func (h *Handler) handlerPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен post text/plain")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("Получено тело запроса: %s", body)

	origin := string(body)

	short := string(h.storage.Save(origin))

	defer h.storage.Flush(h.config)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	response := fmt.Sprintf("%s/%s", h.config.BaseURL, short)
	w.Write([]byte(response))
	defer r.Body.Close()
}

func (h *Handler) handlerPostAPI(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен post application/json")
	jsonURL, err := io.ReadAll(r.Body) // считываем JSON из тела запроса
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Получено тело запроса: %s", jsonURL)

	// создаеём экземпляр структуры для заполнения из JSON
	jsonBody := jsonURLBody{}

	// парсим JSON и записываем результат в экземпляр структуры
	err = json.Unmarshal([]byte(jsonURL), &jsonBody)
	if err != nil {
		log.Printf("Распарсили JSON: %s", err)
	}
	log.Printf("Распарсили JSON string(jsonBody.URL): %s", string(jsonBody.URL))

	// получаем Go-описание типа
	objType := reflect.ValueOf(jsonBody).Type()

	// ищем поле по имени URL
	objField, ok := objType.FieldByName(targetField)
	if !ok {
		panic(fmt.Errorf("field (%s): not found", targetField))
	}

	// получаем метаинформацию о поле
	fieldTag := objField.Tag
	// ищем тег по имени
	tagValue, ok := fieldTag.Lookup(targetTag)
	if !ok {
		panic(fmt.Errorf("tag (%s) for field (%s): not found", targetTag, targetField))
	}

	fmt.Printf("Значение тега (%s) поля (%s): %s\n", targetTag, targetField, tagValue)

	fmt.Printf("Распарсили JSON tagValue: %s jsonBody.URL: %s string(jsonBody.URL): %s", tagValue, jsonBody.URL, string(jsonBody.URL))

	// По ключу помещаем значение localhost map.
	sl := h.storage.Save(jsonBody.URL)

	defer h.storage.Flush(h.config)

	shortURL := fmt.Sprintf("%s/%s", h.config.BaseURL, string(sl))

	log.Printf("Получен shortURL: %s", shortURL)

	// создаем экземпляр структуры и вставляем в него короткий URL для отправки в JSON
	resultURL := ResultURL{
		Result: shortURL,
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// изготавливаем JSON
	shortJSONURL, err := json.Marshal(resultURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Получен shortJSONURL: %s", shortJSONURL)
	//w.Write(shortJSONURL)
	w.Write([]byte(shortJSONURL))
	defer r.Body.Close()
}

func (h *Handler) handlerGet(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		log.Printf("Получен get application/json")
		jsonURL := strings.TrimPrefix(r.URL.Path, "/")
		log.Printf("Получен key %s", jsonURL)
		// создаеём экземпляр структуры для заполнения из JSON
		jsonBody := jsonURLBody{}

		// парсим JSON и записываем результат в экземпляр структуры
		if err := json.Unmarshal([]byte(jsonURL), &jsonBody); err != nil {
			log.Printf("Распарсили JSON: %s", err)
		}
		if url, ok := urls[jsonURL]; ok {
			log.Printf("Отдаем url %s", url)
			// создаем экземпляр структуры и вставляем в него короткий URL для отправки в JSON
			resultURL := ResultURL{
				Result: url,
			}
			// изготавливаем JSON
			longJSONURL, err := json.Marshal(resultURL)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Получен shortJSONURL: %s", longJSONURL)
			// Respond with JSON
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Location", string(longJSONURL))
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Write([]byte(longJSONURL))
			defer r.Body.Close()
		}
	default:
		log.Printf("%s %q", r.Method, html.EscapeString(r.URL.Path))
		log.Printf("Получен get text/html")
		vars := mux.Vars(r)
		key, ok := vars["url"]
		if !ok {
			fmt.Println("url is missing in parameters")
		}
		fmt.Println(`url := `, key)
		keykey := strings.TrimPrefix(r.URL.Path, "/")
		log.Printf("Получен key %s", keykey)
		url, err := h.storage.LinkBy(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		//if url, ok := urls[key]; ok {
		//	log.Printf("Отдаем url %s", url)
		//	w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//	w.Header().Set("Location", url)
		//	w.WriteHeader(http.StatusTemporaryRedirect)
		//	defer r.Body.Close()
		//}
	}

}
