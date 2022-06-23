package helper

import (
	"net/http"
)

func CreateMyCookie(name, value string) *http.Cookie {
	// переопределяем куки
	return &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
}
