package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	"github.com/mihailcoc/shortener/internal/app/handler"
	"github.com/mihailcoc/shortener/internal/app/model"
	"github.com/mihailcoc/shortener/internal/app/shorturl"
)

// Определяем структуру базы данных
type PostgresDatabase struct {
	conn    *sql.DB
	baseURL string
}

// Определяем структуру данных URL
type GetURLData struct {
	OriginalURL string
	IsDeleted   bool
}

// Переопределяем репозиторий базы данных полученную от db как структуру базу данных PostgresDatabase
func DatabaseRepository(baseURL string, db *sql.DB) *PostgresDatabase {
	return &PostgresDatabase{
		conn:    db,
		baseURL: baseURL,
	}
}

// Переопределяем репозиторий базы данных полученную от db как репозиторий Repository
func NewDatabaseRepository(baseURL string, db *sql.DB) handler.Repository {
	return handler.Repository(DatabaseRepository(baseURL, db))
}

// Функция для добавления URL в базу данных
func (db *PostgresDatabase) AddURL(ctx context.Context, longURL model.LongURL, shortURL model.ShortURL, user model.UserID) error {
	// добавляем ряд данных состоящий из user_id origin_url и short_urlв базу данных
	sqlAddRow := `INSERT INTO urls (user_id, origin_url, short_url)
				  VALUES ($1, $2, $3)`
	// выполняем запрос на добавление данных в бд, который не возвращает записи
	_, err := db.conn.ExecContext(ctx, sqlAddRow, user, longURL, shortURL)

	var pgErr *pq.Error
	// Проходимся по всем полям ошибки pgErr
	if errors.As(err, &pgErr) {
		// Если ошибка добавления данных в бд равна 23505
		if pgErr.Code == pgerrcode.UniqueViolation {
			// возвращаем хендлер NewErrorWithDB с Title UniqConstraint
			return handler.NewErrorWithDB(err, handler.UniqConstraint)
		}
	}

	return err
}

func (db *PostgresDatabase) GetURL(ctx context.Context, shortURL model.ShortURL) (model.ShortURL, error) {
	// Пишем запрос к базе данных выбрать URL где указателем будет short_url.
	sqlGetURLRow := `SELECT origin_url, is_deleted FROM urls WHERE short_url=$1 LIMIT 1`
	// Делаем запрос к базе данных результатом будет одна запись
	// Контекст позволяет ограничить по времени или прервать слишком долгие или уже не нужные операции с базой данных.
	row := db.conn.QueryRowContext(ctx, sqlGetURLRow, shortURL)

	result := GetURLData{}
	//Вызываем метод Scan(), который ассоциирует результат с переменной.
	err := row.Scan(&result.OriginalURL, &result.IsDeleted)
	if err != nil {
		return "", sql.ErrNoRows
	}
	// если результат - пустая строка, возвращаем хендлер NewErrorWithDB с Title "Not found".
	if result.OriginalURL == "" {
		return "", handler.NewErrorWithDB(errors.New("not found"), "Not found")
	}
	if result.IsDeleted {
		return "", handler.NewErrorWithDB(errors.New("deleted"), "deleted")
	}

	return result.OriginalURL, nil
}

func (db *PostgresDatabase) isOwner(ctx context.Context, url string, user string) (bool, error) {
	sqlGetURLRow := `SELECT user_id FROM urls WHERE short_url=$1 FETCH FIRST ROW ONLY;`
	query := db.conn.QueryRowContext(ctx, sqlGetURLRow, url)
	result := ""

	err := query.Scan(&result)
	if err != nil {
		return false, err
	}

	return result == user, nil
}

func (db *PostgresDatabase) GetUserURLs(ctx context.Context, user model.UserID) ([]handler.ResponseGetURL, error) {
	var result []handler.ResponseGetURL
	// Пишем запрос к базе данных выбрать URL где указателем будет user_id.
	sqlGetUserURL := `SELECT origin_url, short_url FROM urls WHERE user_id=$1;`
	// Делаем запрос к базе данных результатом будет одна запись
	// Контекст позволяет ограничить по времени или прервать слишком долгие или уже не нужные операции с базой данных.
	rows, err := db.conn.QueryContext(ctx, sqlGetUserURL, user)
	if err != nil {
		return result, err
	}
	if rows.Err() != nil {
		return result, rows.Err()
	}
	// обязательно закрываем бд перед возвратом функции
	defer rows.Close()
	// пробегаем по всем записям
	for rows.Next() {
		// Задаём переменной u структуру ResponseGetURL
		var u handler.ResponseGetURL
		err = rows.Scan(&u.OriginalURL, &u.ShortURL)
		if err != nil {
			return result, err
		}
		u.ShortURL = db.baseURL + u.ShortURL
		result = append(result, u)
	}
	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *PostgresDatabase) Ping(ctx context.Context) error {
	err := db.conn.PingContext(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (db *PostgresDatabase) AddMultipleURLs(ctx context.Context, user model.UserID, urls ...handler.RequestGetURLs) ([]handler.ResponseGetURLs, error) {
	var result []handler.ResponseGetURLs
	// Создаем соединение с бд
	// Get a Tx for making transaction requests.
	tx, err := db.conn.Begin()
	if err != nil {
		return nil, err
	}
	// Формируем контекст запроса к бд.
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO urls (user_id, origin_url, short_url) VALUES ($1, $2, $3)`)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails.
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)
	// Defer a close in case anything fails.
	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)
	// Проходимся по всем url
	for _, u := range urls {
		shortURL := shorturl.ShorterURL(u.OriginalURL)
		// База принимает запрос в теле которого множество URL для сокращения
		if _, err = stmt.ExecContext(ctx, user, u.OriginalURL, shortURL); err != nil {
			return nil, err
		}
		// База возвращает ответ в виде хендлера в котором данные в формате
		//[ { "correlation_id": "<строковый идентификатор из объекта запроса>", "short_url": "<результирующий сокращённый URL>" }, ... ]
		result = append(result, handler.ResponseGetURLs{
			CorrelationID: u.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", db.baseURL, shortURL),
		})
	}

	if err != nil {
		return nil, err
	}
	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *PostgresDatabase) DeleteMultipleURLs(ctx context.Context, user model.UserID, urls ...string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `UPDATE urls SET is_deleted=true WHERE short_url=$1 AND user_id = $2;`)
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	var urlsToDelete []string

	for _, url := range urls {
		isOwner, err := db.isOwner(ctx, url, user)

		if err == nil && isOwner {
			urlsToDelete = append(urlsToDelete, url)
		}
	}

	for _, url := range urlsToDelete {
		if _, err = stmt.ExecContext(ctx, url); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
