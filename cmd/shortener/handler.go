package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

type Body struct {
	Url string
}

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
	mKey := randomString(len(body) / 4)

	urls[mKey] = string(body)

	response := fmt.Sprintf("%s/%s", baseURL, mKey)
	g.String(http.StatusCreated, response)
}

func handlerPostApi(g *gin.Context) {
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
	// responseApi, _ := json.Marshal(response)
	//g.String(http.StatusCreated, response)
	jsonresponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stdout, "%s", jsonresponse)
	g.JSON(200, jsonresponse) // jsonPessoal is still a []byte !!
	if err != nil {
		return
	}

}
