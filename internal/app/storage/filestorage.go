package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/mihailcoc/shortener/internal/app/model"
)

type row struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
	User     string `json:"user"`
}

type Producer interface {
	// интерфейс для закрытия записи файла
	Close() error
}

type Consumer interface {
	// интерфейс для закрытия чтения файла
	Close() error
}

type producer struct {
	file    *os.File
	write   *bufio.Writer
	encoder *json.Encoder
}

type consumer struct {
	file    *os.File
	read    *bufio.Reader
	decoder *json.Decoder
}

type Repository struct {
	urls     model.ShortURLs
	filePath string
	baseURL  string
	usersURL map[model.UserID][]model.ShortURL
	mtx      sync.Mutex
}

func NewProducer(filename string) (*producer, error) {
	// Открываем и записываем файл.
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return nil, err
	}
	// переопределяем вывод
	return &producer{
		file:    file,
		write:   bufio.NewWriter(file),
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) Close() error {
	// закрытие файла
	return p.file.Close()
}

func NewConsumer(filename string) (*consumer, error) {
	// Открываем и читаем файл.
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0777)

	if err != nil {
		return nil, err
	}
	// переопределяем вывод
	return &consumer{
		file:    file,
		read:    bufio.NewReader(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (p *consumer) Close() error {
	// закрытие файла
	return p.file.Close()
}

// функция для записи ряда в файл с помощью NewProducer
func (repo *Repository) writeRow(longURL, shortURL, filePath, userID string) error {
	p, err := NewProducer(filePath)
	if err != nil {
		return err
	}
	// сохраняем данные в JSON
	data, err := json.Marshal(&row{
		LongURL:  longURL,
		ShortURL: shortURL,
		User:     userID,
	})
	if err != nil {
		return err
	}
	// записываем событие в буфер
	if _, err := p.write.Write(data); err != nil {
		return err
	}
	// добавляем перенос строки
	if err := p.write.WriteByte('\n'); err != nil {
		return err
	}
	// записываем буфер в файл
	return p.write.Flush()
}

// функция для добавления URL в репозиторий
func (repo *Repository) AddURL(ctx context.Context, longURL, shortURL string, userID model.UserID) error {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	repo.urls[shortURL] = longURL
	err := repo.writeRow(longURL, shortURL, repo.filePath, userID)
	if err != nil {
		return errors.New("unexpected error when writing row")
	}

	repo.usersURL[userID] = append(repo.usersURL[userID], shortURL)

	return nil
}

// функция для получения URL
func (repo *Repository) GetURL(ctx context.Context, sl model.ShortURL) (model.ShortURL, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	sl, ok := repo.urls[sl]
	if !ok {
		return "", errors.New("url not found")
	}

	return sl, nil
}
