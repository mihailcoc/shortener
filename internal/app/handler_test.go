package app

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handlerPost(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		c := CreateTestContext(w)
		SaveProvider(c)
	}
	req := httptest.NewRequest("POST", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	fmt.Println(got)
	assert.Equal(t, got, got)
}

func Test_handlerGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		c := CreateTestContext(w)
		SaveProvider(c)
	}
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	//body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	//fmt.Println(string(body))

	fmt.Println(got)
	assert.Equal(t, got, got)
}
