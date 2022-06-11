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
	"net/url"

	"github.com/go-chi/chi"
	"github.com/mihailcoc/shortener/internal/app/crypt"
	"github.com/mihailcoc/shortener/internal/app/model"
	"github.com/mihailcoc/shortener/internal/app/shorturl"
	"github.com/mihailcoc/shortener/internal/app/workers"
)

var (
	urls        = make(map[string]string)
	ErrNotFound = errors.New("not found")
	ErrDelete   = errors.New("deleted")
)

// тег, значение которого нужно получить
// имя поля, о котором нужно получить информацию.
const (
	UniqConstraint = "UniqConstraint"
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
	DeleteMultipleURLs(ctx context.Context, user model.UserID, urls ...string) error
}

type Handler struct {
	repo    Repository
	baseURL string
	wp      *workers.WorkerPool
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
		wp:      wp,
	}
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateShortURL")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isURL(string(body)) {
		http.Error(w, "the body cannot be an empty", http.StatusBadRequest)
		return
	}

	userID := crypt.GetUserID(r.Context(), r)

	longURL := model.LongURL(body)
	shortURL := shorturl.ShorterURL(longURL)
	// добавляем URL через интерфейс для добавления URL
	err = h.repo.AddURL(r.Context(), longURL, shortURL, userID)
	if err != nil {
		var dbErr *ErrorWithDB
		// перебирает все поля ошибки dbErr и возвращает true если поле Title равно UniqConstraint
		if errors.As(err, &dbErr) && dbErr.Title == UniqConstraint {
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
	result := map[string]string{}

	body, errReadAll := io.ReadAll(r.Body)

	if errReadAll != nil {
		http.Error(w, errReadAll.Error(), http.StatusInternalServerError)
		return
	}

	url := URL{}
	err := json.Unmarshal(body, &url)
	//err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, "an unexpected error when unmarshaling JSON", http.StatusBadRequest)
		return
	}
	if !isURL(string(url.URL)) {
		http.Error(w, "the URL property is missing", http.StatusBadRequest)
		return
	}
	userID := crypt.GetUserID(r.Context(), r)

	shortURL := shorturl.ShorterURL(url.URL)

	slURL := fmt.Sprintf("%s/%s", h.baseURL, shortURL)

	//type result struct {
	//	result string 'json:"result"'
	//}{}

	//result := map[string]string{}

	err = h.repo.AddURL(r.Context(), url.URL, shortURL, userID)

	if err != nil {
		var dbErr *ErrorWithDB
		if errors.As(err, &dbErr) && dbErr.Title == UniqConstraint {
			result["result"] = slURL
			// result.result = slURL
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			body, err = json.Marshal(result)
			if err != nil {
				http.Error(w, "an unexpected error when marshaling JSON", http.StatusInternalServerError)
				return
			}

			_, err = w.Write(body)

			//buf := bytes.NewBuffer([]byte{})
			//encoder := json.NewEncoder(buf)
			//encoder.SetEscapeHTML(false) // без этой опции символ '&' будет заменён на "\u0026"
			//encoder.Encode(result)

			//jsonResp, _ := json.Marshal(result)
			//_, err = w.Write(jsonResp)
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

	//buf := bytes.NewBuffer([]byte{})
	//encoder := json.NewEncoder(buf)
	//encoder.SetEscapeHTML(false) // без этой опции символ '&' будет заменён на "\u0026"
	//encoder.Encode(result)

	//jsonResp, _ := json.Marshal(result)
	//_, err = w.Write(jsonResp)
	if err != nil {
		http.Error(w, "unexpected error when writing the response body", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) RetrieveShortURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	url, err := h.repo.GetURL(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDelete) {
			w.WriteHeader(http.StatusGone)
			return
		} else if errors.Is(err, ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Location", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	userID := crypt.GetUserID(r.Context(), r)

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
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateBatch(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data []RequestGetURLs

	userID := crypt.GetUserID(r.Context(), r)

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

func (h *Handler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "only DELETE requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDCtx := r.Context().Value(crypt.UserIDCtxName)

	userID := "default"

	if userIDCtx != nil {
		userID = userIDCtx.(string)
	}

	defer r.Body.Close()

	var data []string

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

	var sliceData [][]string

	for i := 10; i <= len(data); i += 10 {
		sliceData = append(sliceData, data[i-10:i])
	}

	rem := len(data) % 10
	if rem > 0 {
		sliceData = append(sliceData, data[len(data)-rem:])
	}

	for _, item := range sliceData {
		func(taskData []string) {
			h.wp.Push(func(ctx context.Context) error {
				err := h.repo.DeleteMultipleURLs(ctx, userID, taskData...)
				return err
			})
		}(item)
	}

	w.WriteHeader(http.StatusAccepted)
}
