package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"main/cmd/shortener"
	"net/http"
)

func handlerGet(g *gin.Context) {
	key := g.Param("key")
	if url, ok := main.urls[key]; ok {
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
	mKey := main.randomString(len(body) / 4)

	main.urls[mKey] = string(body)

	response := fmt.Sprintf("%s/%s", main.baseURL, mKey)
	g.String(http.StatusCreated, response)
}
