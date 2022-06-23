package servers

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mihailcoc/shortener/internal/app/compressor"
	"github.com/mihailcoc/shortener/internal/app/mw"
)

type CustomServer struct {
	addr    string
	key     []byte
	handler *chi.Mux
	s       *http.Server
}

func NewServer(addr string, key []byte, handler *chi.Mux) *CustomServer {
	srv := &http.Server{
		Addr:    addr,
		Handler: mw.Conveyor(handler, compressor.GzipHandle, mw.CookieMiddleware(key)),
	}

	return &CustomServer{
		addr:    addr,
		key:     key,
		handler: handler,
		s:       srv,
	}
}

func (s *CustomServer) StartServer() error {
	log.Printf("Сервер запущен")
	err := s.s.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *CustomServer) Shutdown(ctx context.Context) error {
	log.Printf("Сервер остановлен")
	err := s.s.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
