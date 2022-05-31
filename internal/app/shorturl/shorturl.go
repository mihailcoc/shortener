package shorturl

import (
	"crypto/sha1"
	"encoding/base64"
	"strings"
)

// Определяем функцию для получения кодировки
// Функция возвращает массив в виде строки в кодировке base64.
func ShorterURL(longURL string) string {
	// получаем splitURL после разделения longURL
	splitURL := strings.Split(longURL, "://")
	//Вычислим хеш-функцию SHA1
	hasher := sha1.New()
	//если длина URL меньше 2
	if len(splitURL) < 2 {
		// прочитаем байтовый слайс от начального longURL
		hasher.Write([]byte(longURL))
	} else {
		// прочитаем байтовый слайс от splitURL
		hasher.Write([]byte(splitURL[1]))
	}
	// получаем urlHash закодировав в строку
	urlHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return urlHash
}
