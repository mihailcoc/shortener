package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/mihailcoc/shortener/internal/app/handler"
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

func (repo *Repository) GetURL(ctx context.Context, sl model.ShortURL) (model.ShortURL, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	sl, ok := repo.urls[sl]
	if !ok {
		return "", errors.New("url not found")
	}

	return sl, nil
}

func (repo *Repository) GetUserURLs(ctx context.Context, userID model.UserID) ([]handler.ResponseGetURL, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	var result []handler.ResponseGetURL

	shortLinks := repo.usersURL[userID]

	for _, v := range shortLinks {
		result = append(result, handler.ResponseGetURL{
			ShortURL:    fmt.Sprintf("%s/%s", repo.baseURL, v),
			OriginalURL: repo.urls[v],
		})
	}

	return result, nil
}

func (repo *Repository) Ping(ctx context.Context) error {
	return errors.New("not supported with filebase repository")
}

func (repo *Repository) readRow(reader *bufio.Scanner) (bool, error) {
	if !reader.Scan() {
		return false, reader.Err()
	}
	data := reader.Bytes()
	row := &row{}
	err := json.Unmarshal(data, row)

	if err != nil {
		return false, err
	}
	repo.urls[row.ShortURL] = row.LongURL
	repo.usersURL[row.User] = append(repo.usersURL[row.User], row.ShortURL)

	return true, nil
}

func FileRepository(ctx context.Context, filePath string, baseURL string) *Repository {
	repo := Repository{
		urls:     model.ShortURLs{},
		filePath: filePath,
		baseURL:  baseURL,
		usersURL: map[model.UserID][]model.ShortURL{},
	}
	cns, err := NewConsumer(filePath)
	if err != nil {
		log.Printf("Error with reading file: %v\n", err)
	}
	defer cns.Close()
	reader := bufio.NewScanner(cns.file)
	for {
		ok, err := repo.readRow(reader)

		if err != nil {
			log.Printf("Error while parsing file: %v\n", err)
		}

		if !ok {
			break
		}
	}

	return &repo
}

func NewFileRepository(ctx context.Context, filePath string, baseURL string) handler.Repository {
	return handler.Repository(FileRepository(ctx, filePath, baseURL))
}

func (repo *Repository) AddMultipleURLs(ctx context.Context, user model.UserID, urls ...handler.RequestGetURLs) ([]handler.ResponseGetURLs, error) {
	return nil, nil
}
