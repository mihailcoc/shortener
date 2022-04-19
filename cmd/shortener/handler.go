package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//  описываем структуру JSON в запросе - {"url":"<some_url>"}
type jsonURLBody struct {
	URL string `json:"url"`
}

//  описываем структуру создаваемого JSON вида {"result":"<shorten_url>"}
type ResultURL struct {
	Result string `json:"result"`
}

func handlerPost(g *gin.Context) {
	body, err := io.ReadAll(g.Request.Body)
	if err != nil {
		panic(err)
	}
	//log.Printf("Получено тело запроса: %s", body)
	// По ключу помещаем значение localhost map.
	mKey := randomString(len(body) / 4)

	urls[mKey] = string(body)

	response := fmt.Sprintf("%s/%s", baseURL, mKey)
	g.String(http.StatusCreated, response)
}

func handlerPostAPI(g *gin.Context) {

	switch g.Request.Header.Get("Content-Type") {
	case "application/json":
		log.Printf("Получен post application/json")
		jsonURL, err := io.ReadAll(g.Request.Body) // считываем JSON из тела запроса
		if err != nil {
			panic(err)
		}
		log.Printf("Получено тело запроса: %s", jsonURL)

		// создаеём экземпляр структуры для заполнения из JSON
		jsonBody := jsonURLBody{}

		// парсим JSON и записываем результат в экземпляр структуры
		err = json.Unmarshal(jsonURL, &jsonBody)

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
		g.JSON(http.StatusCreated, shortURL)
		//w := g.NewRecorder()
		//w := http.NewRequest()
		//w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(http.StatusCreated)
		//w.Write(shortJSONURL)
	case "application/xml":
		body, err := io.ReadAll(g.Request.Body)
		if err != nil {
			panic(err)
		}
		log.Printf("Получено тело запроса: %s", body)

		// По ключу помещаем значение localhost map.
		mKey := randomString(len(body) / 4)

		urls[mKey] = string(body)
		response := fmt.Sprintf("%s/%s", baseURL, mKey)
		// Respond with XML
		g.XML(http.StatusCreated, gin.H{"result": response})
	case "text/plain":
		log.Printf("Получен post text/plain")
		body, err := io.ReadAll(g.Request.Body)
		if err != nil {
			panic(err)
		}
		log.Printf("Получено тело запроса: %s", body)

		// По ключу помещаем значение localhost map.
		mKey := randomString(len(body) / 4)

		urls[mKey] = string(body)
		response := fmt.Sprintf("%s/%s", baseURL, mKey)
		// Respond with JSON
		g.JSON(http.StatusCreated, gin.H{"result": response})
	default:
		log.Printf("Получен post default")
		body, err := io.ReadAll(g.Request.Body)
		if err != nil {
			panic(err)
		}
		//log.Printf("Получено тело запроса: %s", body)

		// По ключу помещаем значение localhost map.
		mKey := randomString(len(body) / 4)

		urls[mKey] = string(body)
		response := fmt.Sprintf("%s/%s", baseURL, mKey)
		// Respond with JSON
		g.JSON(http.StatusCreated, gin.H{"result": response})
	}
}

//func handlerGet(g *gin.Context) {
//	key := g.Param("key")
//	if url, ok := urls[key]; ok {
//		g.Redirect(http.StatusTemporaryRedirect, url)
//		return
//	}
//}

func handlerGet(g *gin.Context) {
	switch g.Request.Header.Get("Content-Type") {

	case "application/json":
		log.Printf("Получен get application/json")
		key := g.Param("key")
		if url, ok := urls[key]; ok {
			g.Header("Location", url)
			g.JSON(http.StatusTemporaryRedirect, nil)
			return
		}
	case "application/xml":
		log.Printf("Получен get application/xml")
		key := g.Param("key")
		if url, ok := urls[key]; ok {
			g.Header("Location", url)
			g.XML(http.StatusTemporaryRedirect, nil)
			return
		}
	case "text/plain":
		log.Printf("Получен get text/plain")
		key := g.Param("key")
		if url, ok := urls[key]; ok {
			g.Header("Location", url)
			g.Redirect(http.StatusTemporaryRedirect, url)
			return
		}
	case "application/x-yaml":
		log.Printf("Получен get application/x-yaml")
		key := g.Param("key")
		if url, ok := urls[key]; ok {
			g.Header("Location", url)
			g.YAML(http.StatusTemporaryRedirect, nil)
			return
		}
	default:
		log.Printf("Получен get default")
		key := g.Param("key")
		log.Printf("Получен key %s", key)
		if url, ok := urls[key]; ok {
			log.Printf("Отдаем url %s", url)

			//w := http.NewRequest()
			g.Header("Content-Type", "text/html; charset=utf-8")
			//w.Header.Set("Location", url)
			g.Header("Location", url)
			//w.WriteHeader(http.StatusTemporaryRedirect)
			g.Redirect(http.StatusTemporaryRedirect, url)
			return
		}
	}
}
