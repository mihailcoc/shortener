package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	w            = httptest.NewRecorder()
	resp, engine = gin.CreateTestContext(w)
	reqbody      = strings.NewReader("http://rqls3b.com/bnclubmjprl")
	key          = string("key")
)

func Test_handlerPost(t *testing.T) {
	type args struct {
		g *gin.Context
	}
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			// TODO: Add test cases.
			name: "positive test #1",
			want: want{
				code:        201,
				response:    `http://localhost:8080/gmwjgsa`,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// handlerPost(tt.args.g)

			req := httptest.NewRequest(http.MethodPost, "/", reqbody)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем handler
			engine.POST("/", handlerPost)
			// запускаем сервер
			engine.ServeHTTP(http.ResponseWriter(w), req)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}
			key = string(resBody)
			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func Test_handlerPostAPI(t *testing.T) {
	type args struct {
		g *gin.Context
	}
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			// TODO: Add test cases.
			name: "positive test #3",
			want: want{
				code:        201,
				response:    `http://localhost:8080/gmwjgsa`,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создаём новый Request
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", reqbody)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем handler
			engine.POST("/api/shorten", handlerPostAPI)
			// запускаем сервер
			engine.ServeHTTP(http.ResponseWriter(w), req)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}
			key = string(resBody)
			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}	}
}

func Test_handlerGet(t *testing.T) {
	type args struct {
		g *gin.Context
	}
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		// TODO: Add test cases.
		{
			name: "positive test #2",
			// args:{}
			want: want{
				code:        307,
				response:    ``,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//"key" = key
			// создаём новый Recorder
			w := httptest.NewRecorder()
			w.Header().Set("Location", key)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusTemporaryRedirect)
			req := httptest.NewRequest(http.MethodGet, "/:key", nil)
			// определяем handler
			engine.GET("/:key", handlerGet)
			// запускаем сервер
			engine.ServeHTTP(http.ResponseWriter(w), req)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

			// заголовок ответа
			if w.Header().Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, w.Header().Get("Content-Type"))
			}
		})
	}
}

func Test_handlerGetAPI(t *testing.T) {
	type args struct {
		g *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerGetAPI(tt.args.g)
		})
	}
}
