package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mihailcoc/shortener/internal/app/model"
	"github.com/mihailcoc/shortener/internal/app/mw"
	"github.com/mihailcoc/shortener/internal/app/shorturl"
)

var (
	urls = make(map[string]string)
)

// тег, значение которого нужно получить
// имя поля, о котором нужно получить информацию.
const (
	targetField = "URL"
	targetTag   = "json"
)

type Repository interface {
	// интерфейс для добавления URL
	AddURL(ctx context.Context, longURL model.LongURL, shortURL model.ShortURL, user model.UserID) error
	// интерфейс для получения URL
	GetURL(ctx context.Context, shortURL model.ShortURL) (model.ShortURL, error)
	// интерфейс для получения URL созданных пользователем
	GetUserURLs(ctx context.Context, user model.UserID) ([]ResponseGetURL, error)
	// интерфейс для проверки связи с DB
	Ping(ctx context.Context) error
	// интерфейс для добавления множества URL
	AddMultipleURLs(ctx context.Context, user model.UserID, urls ...RequestGetURLs) ([]ResponseGetURLs, error)
}

//  описываем структуру Handler в запросе на получение данных их репозитория
type Handler struct {
	repo    Repository
	baseURL string
}

//  описываем структуру JSON в запросе - {"url":"<some_url>"}
type URL struct {
	URL string `json:"url"`
}

//  описываем структуру создаваемого ответа JSON вида {"ShortURL":"<short_url>", "OriginalURL":"<original_url>"}
type ResponseGetURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

//  описываем структуру создаваемого запроса JSON вида {"CorrelationID":"<correlation_id>", "OriginalURL":"<original_url>"}
type RequestGetURLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

//  описываем структуру создаваемого ответа JSON вида {"CorrelationID":"<correlation_id>", "ShortURL":"<short_url>"}
type ResponseGetURLs struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// описываем структуру ошибки
type ErrorWithDB struct {
	Err   error
	Title string
}

// Error добавляет поддержку интерфейса error для типа ErrorWithDB.
func (err *ErrorWithDB) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

// Unwrap добавляет поддержку интерфейса error для типа ErrorWithDB.
func (err *ErrorWithDB) Unwrap() error {
	return err.Err
}

// NewErrorWithDB упаковывает ошибку err в тип ErrorWithDB c текущим временем.
func NewErrorWithDB(err error, title string) error {
	return &ErrorWithDB{
		Err:   err,
		Title: title,
	}
}

//  описываем новый handler
func NewHandler(repo Repository, baseURL string) *Handler {
	return &Handler{
		repo:    repo,
		baseURL: baseURL,
	}
}

//  описываем новый handler по созданию короткого URL.
func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	//Используем конструкцию отложенного исполнения defer, чтобы закрыть соединение и освободить ресурс
	defer r.Body.Close()
	// читаем тело запроса
	body, err := io.ReadAll(r.Body)
	// если ошибка то статус 500
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// если в теле запроса нет URL, то статус 400
	if len(body) == 0 {
		http.Error(w, "the body cannot be an empty", http.StatusBadRequest)
		return
	}
	// получаем из контекста запроса, значение userID
	userIDCtx := r.Context().Value(mw.UserIDCtxName)
	// определяем дефолное значение userID
	userID := "default"
	// если полученное из контекста запроса, значение userID не равно нулю то присваиваем переменной userID
	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}
	// присваиваем переменной longURL значение из тела запроса по форме из модели LongURL
	longURL := model.LongURL(body)
	// присваиваем переменной shortURL значение из функции ShorterURL
	shortURL := shorturl.ShorterURL(longURL)
	// добавляем URL через интерфейс для добавления URL
	err = h.repo.AddURL(r.Context(), longURL, shortURL, userID)
	// если ошибка существует
	if err != nil {
		// присваиваем переменной dbErr форму *ErrorWithDB
		var dbErr *ErrorWithDB
		// перебирает все поля ошибки dbErr и возвращает true если поле Title равно UniqConstraint
		if errors.As(err, &dbErr) && dbErr.Title == "UniqConstraint" {
			// и добавляет в хэдеру в ответе Content-Type
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			// и добавляет в хэдеру в ответе статус 409
			w.WriteHeader(http.StatusConflict)
			// выводим URL из запроса и shorlURL
			slURL := fmt.Sprintf("%s/%s", h.baseURL, shortURL)
			// записываем вывод в файл методом Write
			_, err = w.Write([]byte(slURL))
			if err != nil {
				http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
			}

			return
		}
		// если такой ошибки нет в полях ошибки то выводим статус 500 в заголовок header
		w.WriteHeader(http.StatusInternalServerError)
	}
	// и добавляет в хэдеру в ответе Content-Type
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// и добавляет в хэдеру в ответе статус 201
	w.WriteHeader(http.StatusCreated)
	// выводим URL из запроса и shorlURL
	slURL := fmt.Sprintf("%s/%s", h.baseURL, shortURL)
	// записываем вывод в файл методом Write
	_, err = w.Write([]byte(slURL))
	if err != nil {
		http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
	}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	result := map[string]string{}

	body, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		http.Error(w, errReadAll.Error(), http.StatusInternalServerError)
		return
	}

	url := URL{}

	err := json.Unmarshal(body, &url)
	if err != nil {
		http.Error(w, "an unexpected error when unmarshaling JSON", http.StatusInternalServerError)
		return
	}

	if url.URL == "" {
		http.Error(w, "the URL property is missing", http.StatusBadRequest)
		return
	}

	userIDCtx := r.Context().Value(mw.UserIDCtxName)

	userID := "default"

	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}

	shortURL := shorturl.ShorterURL(url.URL)

	slURL := fmt.Sprintf("%s/%s", h.baseURL, shortURL)

	err = h.repo.AddURL(r.Context(), url.URL, shortURL, userID)
	if err != nil {
		var dbErr *ErrorWithDB
		if errors.As(err, &dbErr) && dbErr.Title == "UniqConstraint" {
			result["result"] = slURL

			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusConflict)

			body, err = json.Marshal(result)
			if err != nil {
				http.Error(w, "an unexpected error when marshaling JSON", http.StatusInternalServerError)
				return
			}

			_, err = w.Write(body)
			if err != nil {
				http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
				return
			}

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result["result"] = slURL

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	body, err = json.Marshal(result)
	if err != nil {
		http.Error(w, "an unexpected error when marshaling JSON", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
		return
	}
}

// описываем хэндлер получающий все URL созданные пользователем
func (h *Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	// получаем из контекста запроса, значение userID
	userIDCtx := r.Context().Value(mw.UserIDCtxName)
	// определяем дефолное значение userID
	userID := "default"
	// если полученное из контекста запроса, значение userID не равно нулю то присваиваем переменной userID
	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}
	// получаем URL пользователя через интерфейс для получения URL
	urls, err := h.repo.GetUserURLs(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		http.Error(w, errors.New("no content").Error(), http.StatusNoContent)
		return
	}
	// переводим URL в формат JSON
	body, err := json.Marshal(urls)

	if err == nil {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		_, err = w.Write(body)
		if err == nil {
			return
		}
	}
}

func (h *Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	// проверяем связь с бд через интерфейс для проверки связи с DB
	err := h.repo.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateBatch(w http.ResponseWriter, r *http.Request) {
	//Используем конструкцию отложенного исполнения defer, чтобы закрыть соединение и освободить ресурс
	defer r.Body.Close()
	// определяем переменную data
	var data []RequestGetURLs
	// получаем из контекста запроса, значение userID
	userIDCtx := r.Context().Value(mw.UserIDCtxName)
	// определяем дефолное значение userID
	userID := "default"
	// если полученное из контекста запроса, значение userID не равно нулю то присваиваем переменной userID
	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}
	// считываем JSON из тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// переводим из формата JSON
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// получаем часть URL пользователя через интерфейс для получения части URL
	urls, err := h.repo.AddMultipleURLs(r.Context(), userID, data...)
	if err != nil {
		log.Println("err.Error(): ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err = json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
		return
	}
}
