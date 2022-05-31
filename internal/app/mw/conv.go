package mw

import "net/http"

type Middleware func(http.Handler) http.Handler

// задаем конвейер обработки хэндлеров
func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	// для каждого хэндлера
	for _, middleware := range middlewares {
		// присоединяем к общему пулу хэндлеров
		h = middleware(h)
	}
	// возвращаем пул хэндлеров
	return h
}
