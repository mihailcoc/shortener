package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Repository interface {
	LinkBy(sl string) (string, error)
	Save(url string) (sl string)
	Load(c Config) error
	Flush(c Config) error
}

type Producer interface {
	WriteEvent(event *storage)
	Close() error
}

type Consumer interface {
	ReadEvent() (*storage, error)
	Close() error
}

type storage struct {
	Data map[string]string
	mu   sync.Mutex
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

func (s *storage) LinkBy(sl string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// По ключу достаём значение в Data map.
	link, ok := s.Data[sl]
	if !ok {
		return link, errors.New("url not found")
	}

	return link, nil
}

func (s *storage) Save(url string) (mKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Генерируем короткий url.
	mKey = string(randomString(len(url) / 4))
	// Сохраняем короткий url.
	s.Data[mKey] = url
	return
}

func NewProducer(filename string) (*producer, error) {
	// Открываем и записываем файл.
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return nil, err
	}

	return &producer{
		file:    file,
		write:   bufio.NewWriter(file),
		encoder: json.NewEncoder(file),
	}, nil
}

func NewConsumer(filename string) (*consumer, error) {
	// Открываем и читаем файл.
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0777)

	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		read:    bufio.NewReader(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (p *producer) Close() error {
	return p.file.Close()
}

func (p *consumer) Close() error {
	return p.file.Close()
}

func (s *storage) Load(c Config) error {
	if c.FileStoragePath == "" {
		return nil
	}
	// Загружаем данные из файла
	cns, err := NewConsumer(c.FileStoragePath)

	if err != nil {
		return err
	}

	cns.decoder.Decode(&s.Data)

	return nil
}

func (s *storage) Flush(c Config) error {
	if c.FileStoragePath == "" {
		return nil
	}
	// Загружаем данные из файла
	p, err := NewProducer(c.FileStoragePath)

	if err != nil {
		return err
	}

	p.encoder.Encode(&s.Data)

	return p.write.Flush()
}

func NewStorage() *storage {
	return &storage{
		Data: make(map[string]string),
	}
}
