package mw

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/mihailcoc/shortener/internal/app/crypt"
	"github.com/mihailcoc/shortener/internal/app/helper"
)

const CookieUserIDName = "user_id"

type ContextType string

const UserIDCtxName ContextType = "ctxUserId"

func CookieMiddleware(key []byte) func(next http.Handler) http.Handler {
	// переопределяем вывод куки как handler
	return func(next http.Handler) http.Handler {
		// переопределяем вывод как http.handler функцию
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// получаем id пользователя из request.Cookie
			cookieUserID, _ := r.Cookie(CookieUserIDName)
			// Зашифровываем куки
			encryptor, err := crypt.NewCipherBlock(key)
			// Проверяем на ошибки
			if err != nil {
				return
			}
			if cookieUserID != nil {
				// расшифровываем
				userID, err := encryptor.Decode(cookieUserID.Value)

				if err == nil {
					// run the original handler
					next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtxName, userID)))
					return
				}
			}
			// получаем userID по формату uuid
			userID, err := uuid.NewV4()
			if err != nil {
				return
			}
			// зашифровываем
			encoded := encryptor.Encode(userID.Bytes())
			// создаем куки по форме http.Cookie
			cookie := helper.CreateMyCookie(CookieUserIDName, encoded)
			// устанавливаем куки
			http.SetCookie(w, cookie)
			// run the original handler
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtxName, userID.String())))
		})
	}
}
