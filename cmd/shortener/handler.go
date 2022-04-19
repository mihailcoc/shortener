package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Body struct {
	URL string `json:"url"`
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
		g.JSON(http.StatusCreated, gin.H{"result": strings.TrimSpace(response)})
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
		g.JSON(http.StatusCreated, gin.H{"result": strings.TrimSpace(response)})
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
		if url, ok := urls[key]; ok {
			log.Printf("Отдаем url %s", url)
			url = strings.TrimSpace(url)
			log.Printf("Отдаем url после strings.TrimSpace %s", url)
			g.Header("Location", url)
			g.Redirect(http.StatusTemporaryRedirect, url)
			return
		}
	}
}
