package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mihailcoc/shortener/cmd/shortener/configs"
)

var (
	reqbody = strings.NewReader("http://rqls3b.com/bnclubmjprl")
	key     = string("key")
)

func Test_handlerPost(t *testing.T) {
	c := configs.NewConfig()
	h := NewHandler(c)

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	type request struct {
		method string
		target string
		path   string
	}

	tests := []struct {
		name    string
		want    want
		request request
	}{
		{
			// TODO: Add test cases.
			name: "positive test #1",
			want: want{
				code:        201,
				response:    "http://localhost:8080/gmwjgsa",
				contentType: "text/plain; charset=utf-8",
			},
			request: request{
				method: http.MethodPost,
				target: "http://rqls3b.com/bnclubmjprl",
				path:   "/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создаём тело запроса
			reqbody = strings.NewReader("http://rqls3b.com/bnclubmjprl")
			// создаем request
			request := httptest.NewRequest(http.MethodPost, "/", reqbody)
			// создаём новый Recorder
			recorder := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(h.HandlerPost)
			// запускаем сервер
			h.ServeHTTP(http.ResponseWriter(recorder), request)
			res := recorder.Result()
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, recorder.Code)
			}
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, recorder.Body.String())
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
	c := configs.NewConfig()
	h := NewHandler(c)

	type want struct {
		code        int
		response    string
		contentType string
		string
	}
	type request struct {
		method string
		target string
		path   string
		body   string
	}

	tests := []struct {
		name    string
		want    want
		request request
	}{
		{
			// TODO: Add test cases.
			name: "positive test #2",
			want: want{
				code: 201,
				//response:    "http://127.0.0.1:8080/pgatlmo",
				response:    "{\"result\":\"http://localhost:8080/gmwjgsa\"}",
				contentType: "application/json",
			},
			request: request{
				method: http.MethodPost,
				//target: "http://rqls3b.com/bnclubmjprl",
				target: "/",
				path:   "/",
				body:   "{\"url\":\"http://rqls3b.com/bnclubmjprl\"}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создаём тело запроса
			reader := strings.NewReader(tt.request.body)
			// создаем request
			request := httptest.NewRequest(tt.request.method, tt.request.target, reader)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// задаём Content-Type
			w.Header().Set("Content-Type", "application/json")
			// задаем статус
			w.WriteHeader(http.StatusCreated)
			// определяем хендлер
			h := http.HandlerFunc(h.HandlerPostAPI)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, resBody)
			}
			key = string(resBody)
			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func Test_handlerGet(t *testing.T) {
	c := configs.NewConfig()
	h := NewHandler(c)

	type want struct {
		code        int
		response    string
		contentType string
	}
	type request struct {
		method string
		target string
		path   string
	}
	tests := []struct {
		name    string
		want    want
		request request
	}{
		// TODO: Add test cases.
		{
			name: "positive test #3",
			// args:{}
			want: want{
				code:        307,
				response:    "{\"url\":\"http://rqls3b.com/bnclubmjprl\"}",
				contentType: "text/plain; charset=utf-8",
			},
			request: request{
				method: http.MethodGet,
				target: "http://localhost:8080/gmwjgsa",
				path:   "/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// задаем header
			w.Header().Set("Location", key)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusTemporaryRedirect)
			// создаем request
			request := httptest.NewRequest(tt.request.method, tt.request.target, nil)

			// определяем хендлер
			h := http.HandlerFunc(h.HandlerGet)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()
			// закрываем тело запроса
			defer res.Body.Close()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// заголовок ответа
			if w.Header().Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, w.Header().Get("Content-Type"))
			}
		})
	}
}
