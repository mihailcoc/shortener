package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

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
	default:
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

	case "application/json; charset=utf-8":
		log.Printf("Получен application/json")
		key := g.Param("key")
		if url, ok := urls[key]; ok {
			g.Header("Location", url)
			g.JSON(http.StatusTemporaryRedirect, gin.H{"url": url})
			return
		}
	case "application/xml":
		key := g.Param("key")
		if url, ok := urls[key]; ok {
			g.Header("Location", url)
			g.Redirect(http.StatusTemporaryRedirect, url)
			return
		}
	default:
		log.Printf("Получен default")
		key := g.Param("key")
		if url, ok := urls[key]; ok {
			g.Header("Location", url)
			g.Redirect(http.StatusTemporaryRedirect, url)
			return
		}
	}
}
