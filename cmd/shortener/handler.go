package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Body struct {
	URL string `json:"str_value"`
}

var (
	urls = make(map[string]string)
)

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
	mKey := randomString(len(body) / 4)

	urls[mKey] = string(body)

	response := fmt.Sprintf("%s/%s", baseURL, mKey)
	g.String(http.StatusCreated, response)
}

func handlerPostAPI(g *gin.Context) {
	// var v Body
	//if err := json.NewDecoder(g.Request.Body).Decode(&v); err != nil {
	// http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	body, err := io.ReadAll(g.Request.Body)
	if err != nil {
		g.String(http.StatusBadRequest, "bad request")
		return
	}
	value := Body{}
	if err := json.Unmarshal([]byte(body), &value); err != nil {
		panic(err)
	}
	// По ключу помещаем значение localhost map.
	mKey := randomString(len(body) / 4)

	urls[mKey] = string(body)
	response := fmt.Sprintf("%s/%s", baseURL, mKey)
	responsebyte, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	g.String(http.StatusCreated, string(responsebyte))
}
