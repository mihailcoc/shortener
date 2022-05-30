ackage compressor

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Функция архивации
func GzipHandle(next http.Handler) http.Handler {
	// переопределяем вывод функции как хэндлер
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// с помощью пакета gzip читаем тело запроса
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Println("decompress error", err)
				next.ServeHTTP(w, r)
				return
			}
			//откладываем закрытие ридера
			defer reader.Close()
			r.Body = reader
		}
		// если gzip не поддерживается, передаём управление
		// дальше без изменений
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			log.Println("compress error", err)
			next.ServeHTTP(w, r)
			return
		}
		//откладываем закрытие обработчика
		defer gz.Close()
		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}
