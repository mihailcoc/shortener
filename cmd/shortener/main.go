package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

var m = make(map[string]string)

const (
	addr    = "localhost:8080"
	scheme  = "http"
	baseURL = scheme + "://" + addr
)

func main() {
	server := gin.Default()
	server.GET("/:key", handlerGet)
	server.POST("/", handlerPost)
	server.Run(addr)
}

// если методом GET
func handlerGet(c *gin.Context) {
	// извлекаем фрагмент id из URL запроса GET /{id}
	key := c.Param("key")
	if url, ok := m[key]; ok {
		c.Redirect(http.StatusTemporaryRedirect, url)
		return
	} else {
		c.String(http.StatusNotFound, "url not found")
		return
	}
}

// если методом POST
func handlerPost(c *gin.Context) {
	b, err := io.ReadAll(c.Request.Body)
	// обрабатываем ошибку
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	// Генерируем ключ
	mKey := randomString(len(b) / 4)
	// По ключу проверяем наличие в map.
	_, keyIsInMap := m[mKey]
	if !keyIsInMap {
		fmt.Println("key not in map")
	}

	m[mKey] = string(b)
	// fmt.Println("m", m)

	responseURL := fmt.Sprintf("%s/%s", baseURL, mKey)
	//fmt.Println("responseURL", responseURL)
	c.String(http.StatusCreated, responseURL)
	// fmt.Println("c.Writer.Header()", c.Writer.Header())
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(97, 122))
	}
	return string(bytes)
}
