package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

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
	value := Body{}
	if err := json.NewDecoder(g.Request.Body).Decode(&value); err != nil {
		http.Error(http.ResponseWriter(httptest.NewRecorder()), err.Error(), http.StatusBadRequest)
		return
		//log.Fatal(err)
	}
	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		g.String(http.StatusBadRequest, "bad request")
		return
		//log.Fatal(err)
	}

	// По ключу помещаем значение localhost map.
	mKey := randomString(len(body) / 4)

	urls[mKey] = string(body)
	response := fmt.Sprintf("%s/%s", baseURL, mKey)
	//fmt.Sprintf(response)
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	//encoder.Encode(v)
	g.String(http.StatusCreated, string(response))
	encoder.Encode(g)
	//v := struct {
	//	Url string
	//}{
	//	Url: "http://mysite.com?id=1234&param=2",
	//}
	//buf := bytes.NewBuffer([]byte{})
	//encoder := json.NewEncoder(buf)
	//encoder.SetEscapeHTML(false)
	//encoder.Encode(v)
	fmt.Println(buf)
}

func handlerGet(g *gin.Context) {
	key := g.Param("key")
	if url, ok := urls[key]; ok {
		g.Redirect(http.StatusTemporaryRedirect, url)
		return
	}
}

func handlerGetAPI(g *gin.Context) {
	key := g.Param("key")
	if url, ok := urls[key]; ok {
		g.Redirect(http.StatusTemporaryRedirect, url)
		return
	}
}
