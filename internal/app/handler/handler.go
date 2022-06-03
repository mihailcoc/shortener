package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
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
	AddURL(ctx context.Context, longURL model.LongURL, shortURL model.ShortURL, user model.UserID) error
	GetURL(ctx context.Context, shortURL model.ShortURL) (model.ShortURL, error)
	// интерфейс для получения URL созданных пользователем
	GetUserURLs(ctx context.Context, user model.UserID) ([]ResponseGetURL, error)
	// интерфейс для проверки связи с DB
	Ping(ctx context.Context) error
	// интерфейс для добавления множества URL
	AddMultipleURLs(ctx context.Context, user model.UserID, urls ...RequestGetURLs) ([]ResponseGetURLs, error)
}

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

type ErrorWithDB struct {
	Err   error
	Title string
}

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

func NewHandler(repo Repository, baseURL string) *Handler {
	return &Handler{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateShortURL")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// непонятно какой пакет использовать, чтобы применить !isURL
	if len(body) == 0 {
		http.Error(w, "the body cannot be an empty", http.StatusBadRequest)
		return
	}
	userIDCtx := r.Context().Value(mw.UserIDCtxName)
	userID := "default"
	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}
	longURL := model.LongURL(body)
	shortURL := shorturl.ShorterURL(longURL)
	// добавляем URL через интерфейс для добавления URL
	err = h.repo.AddURL(r.Context(), longURL, shortURL, userID)
	if err != nil {
		var dbErr *ErrorWithDB
		// перебирает все поля ошибки dbErr и возвращает true если поле Title равно UniqConstraint
		if errors.As(err, &dbErr) && dbErr.Title == "UniqConstraint" {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			slURL := fmt.Sprintf("%s/%s", h.baseURL, shortURL)
			_, err = w.Write([]byte(slURL))
			if err != nil {
				http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	slURL := fmt.Sprintf("%s/%s", h.baseURL, shortURL)
	_, err = w.Write([]byte(slURL))
	if err != nil {
		http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
	}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	url := URL{}
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, "an unexpected error when unmarshaling JSON", http.StatusBadRequest)
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

	result := map[string]string{}

	err = h.repo.AddURL(r.Context(), url.URL, shortURL, userID)

	//body, errReadAll := io.ReadAll(r.Body)
	//if errReadAll != nil {
	//	http.Error(w, errReadAll.Error(), http.StatusInternalServerError)
	//	return
	//}
	if err != nil {
		var dbErr *ErrorWithDB
		if errors.As(err, &dbErr) && dbErr.Title == "UniqConstraint" {
			result["result"] = slURL
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			buf := bytes.NewBuffer([]byte{})
			encoder := json.NewEncoder(buf)
			encoder.SetEscapeHTML(false) // без этой опции символ '&' будет заменён на "\u0026"
			encoder.Encode(result)
			//body, err = json.Marshal(result)
			//if err != nil {
			//	http.Error(w, "an unexpected error when marshaling JSON", http.StatusInternalServerError)
			//	return
			//}
			jsonResp, _ := json.Marshal(result)
			_, err = w.Write(jsonResp)
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
	//body, err = json.Marshal(result)
	//if err != nil {
	//	http.Error(w, "an unexpected error when marshaling JSON", http.StatusInternalServerError)
	//	return
	//}

	//buf := bytes.NewBuffer([]byte{})
	//encoder := json.NewEncoder(buf)
	//encoder.SetEscapeHTML(false) // без этой опции символ '&' будет заменён на "\u0026"
	//encoder.Encode(result)
	//err = json.NewEncoder(w).Encode(result)
	jsonResp, _ := json.Marshal(result)
	_, err = w.Write(jsonResp)
	if err != nil {
		http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) RetrieveShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "the parameter is missing", http.StatusBadRequest)
		return
	}
	url, err := h.repo.GetURL(r.Context(), id)
	if err != nil {
		var dbErr *ErrorWithDB

		if errors.As(err, &dbErr) && dbErr.Title == "deleted" {
			w.WriteHeader(http.StatusGone)
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Add("Location", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	userIDCtx := r.Context().Value(mw.UserIDCtxName)
	userID := "default"
	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}
	urls, err := h.repo.GetUserURLs(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(urls) == 0 {
		http.Error(w, errors.New("no content").Error(), http.StatusNoContent)
		return
	}
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
	err := h.repo.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateBatch(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data []RequestGetURLs
	userIDCtx := r.Context().Value(mw.UserIDCtxName)
	userID := "default"
	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
