package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type want struct {
	code        int
	response    string
	contentType string
}

func Test_handlerGet(t *testing.T) {
	type args struct {
		g *gin.Context
	}
	g := handlerGet
	var tests = []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
		{
			name: "positive test #1",
			args: args{
				g * gin.Context,
			},
			want: want{
				code:        307,
				response:    `{"status":"redirect"}`,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerGet(tt.args.g)
			request := httptest.NewRequest(http.MethodGet, "/status", nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(handlerGet)
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
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func Test_handlerPost(t *testing.T) {
	type args struct {
		g *gin.Context
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			// TODO: Add test cases.
			name: "positive test #2",
			want: want{
				code:        200,
				response:    `{"status":"redirect"}`,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerPost(tt.args.g)
			request := httptest.NewRequest(http.MethodPost, "/status", nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()
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
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "", body)

	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	mKey := randomString(len(b) / 4)
	responseURL := "http://" + r.Host + r.URL.String() + mKey

	resp, body = testRequest(t, ts, "POST", "/")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, responseURL, body)

}
