package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mihailcoc/shortener/cmd/shortener/configs"
	"github.com/mihailcoc/shortener/cmd/shortener/router"
	"github.com/mihailcoc/shortener/internal/app/handler"
	"github.com/mihailcoc/shortener/internal/app/servers"
	"github.com/mihailcoc/shortener/internal/app/storage"
	"golang.org/x/sync/errgroup"
)

func main() {
	var httpServer *servers.CustomServer
	// инициируем контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	cfg := configs.NewConfig()

	var repo handler.Repository

	// Если переменная бд задана то создаём файловый репозиторий.
	repo = storage.NewFileRepository(ctx, cfg.FileStoragePath, cfg.BaseURL)

	g, ctx := errgroup.WithContext(ctx)

	handler := router.NewRouter(repo, cfg)
	// Запускаем функцию с контекстом errgroup.WithContext(ctx)
	g.Go(func() error {
		//Создаем новый сервер
		httpServer = servers.NewServer(cfg.ServerAddress, cfg.Key, handler)
		//Запускаем новый сервер
		err := httpServer.StartServer()
		if err != nil {
			return err
		}

		log.Printf("httpServer starting at: %v", cfg.ServerAddress)

		return nil
	})

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}
	//
	log.Println("Receive shutdown signal")
	// Задаем контекст остановки сервера через 5 секунд
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Откладываем остановку сервера
	defer shutdownCancel()
	// Если сервер запущен
	if httpServer != nil {
		// Если сервер запущен то останавливаем сервер через 5 секунд
		_ = httpServer.Shutdown(shutdownCtx)
	}
	// Выводим сообщение после остановки сервера
	err := g.Wait()
	if err != nil {
		log.Printf("server returning an error: %v", err)
		os.Exit(2)
	}
}
