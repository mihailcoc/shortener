package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

var urls = make(map[string]string)

func handlerGet(g *gin.Context) {
	key := g.Param("key")
	if url, ok := urls[key]; ok {
		g.Redirect(http.StatusTemporaryRedirect, url)
		return
	}
}

func handlerPost(g *gin.Context) {
	body, err := io.ReadAll(g.Request.Body)
	if err != nil {
		g.String(http.StatusBadRequest, "bad request")
		return
	}
	// По ключу помещаем значение localhost map.
	mKey := service.randomString(len(body) / 4)

	urls[mKey] = string(body)

	response := fmt.Sprintf("%s/%s", main.baseURL, mKey)
	g.String(http.StatusCreated, response)
}
