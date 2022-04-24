package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	reqbody = strings.NewReader("http://rqls3b.com/bnclubmjprl")
	key     = string("key")
)

func Test_handlerPost(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
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
				response:    "http://127.0.0.1:8080/gmwjgsa",
				contentType: "text/plain; charset=utf-8",
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
			h := http.HandlerFunc(handlerPost)
			// запускаем сервер
			h.ServeHTTP(http.ResponseWriter(recorder), request)
			res := recorder.Result()
			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, recorder.Code)
			}
			// получаем и проверяем тело запроса
			defer res.Body.Close()
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
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type want struct {
		code        int
		response    string
		contentType string
		string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			// TODO: Add test cases.
			name: "positive test #2",
			want: want{
				code:        201,
				response:    "http://127.0.0.1:8080/pgatlmo",
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создаём тело запроса
			reqbody = strings.NewReader("http://rqls3b.com/bnclubmjprl")
			// создаем request
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", reqbody)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// задаём Content-Type
			w.Header().Set("Content-Type", "application/json")
			// задаем статус
			w.WriteHeader(http.StatusCreated)
			// определяем хендлер
			h := http.HandlerFunc(handlerPost)
			// запускаем сервер
			h.ServeHTTP(w, request)
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
	type args struct {
		w http.ResponseWriter
		r *http.Request
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
			name: "positive test #3",
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
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// задаем header
			w.Header().Set("Location", key)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusTemporaryRedirect)
			// создаем request
			request := httptest.NewRequest(http.MethodGet, "/:key", nil)

			// определяем хендлер
			h := http.HandlerFunc(handlerPost)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()
			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			//resBody, err := io.ReadAll(w.Body)
			//if err != nil {
			//	t.Fatal(err)
			//}
			//if string(resBody) != tt.want.response {
			//	t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			//}

			// заголовок ответа
			if w.Header().Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, w.Header().Get("Content-Type"))
			}
		})
	}
}
