package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var m = make(map[string]string)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		// если методом POST
		b, err := io.ReadAll(c.Request.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(c.Writer, err.Error(), 500)
			return
		}
		// Генерируем ключ
		mKey := randomString(len(b) / 4)
		// По ключу проверяем наличие в map.
		if intid, ok := strconv.Atoi(m[mKey]); ok != nil {
			fmt.Println("Значение в map уже задано:", strconv.Itoa(intid))
		}
		// По ключу помещаем значение localhost map.
		m[mKey] = string(b)
		// Генерируем ответ
		responseURL := "http://" + c.Request.Host + c.Request.URL.String() + mKey
		//fmt.Println(responseURL)
		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		//w.Header().Set("Location", responseURL)
		c.Writer.WriteHeader(http.StatusCreated)
		c.Writer.Write([]byte(responseURL))
		//fmt.Fprint(w)
	})
	// если методом GET
	r.GET("/", func(c *gin.Context) {
		// извлекаем фрагмент id из URL запроса GET /{id}
		q := strings.TrimPrefix(c.Request.URL.Path, "/")
		// fmt.Println("q", q)
		if q == "" {
			http.Error(c.Writer, "The query parameter is missing", http.StatusBadRequest)
			return
		}
		// достаем из map оригинальный URL
		//origURL, exists := m[q]
		//if exists {
		//	http.Error(w, "The query parameter is missing", http.StatusBadRequest)
		//	return
		//}
		origURL := m[q]
		fmt.Println("origURL ", origURL)
		// устанавливаем в заголовке оригинальный URL
		c.Writer.Header().Set("Location", origURL)
		// устанавливаем статус-код 307
		c.Writer.WriteHeader(http.StatusTemporaryRedirect)
		// отдаем редирект на собственный url и код 307
		// fmt.Fprint(c.Writer)
	})
	return r
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

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
