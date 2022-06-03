package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Задаём функцию соединения с базой данных
func Conn(driverName, dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn can not be missing")
	}

	if driverName == "" {
		return nil, fmt.Errorf("driver name can not be missing")
	}
	// открываем соединение с бд
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return db, err
	}
	log.Println("Connect to database")
	// возвращаем бд
	return db, nil
}

// Функция создаёт поля бд.
func SetUpDataBase(ctx context.Context, db *sql.DB) error {
	// описываем форму полей в бд
	sqlCreateDB := `CREATE TABLE IF NOT EXISTS urls (
								id serial PRIMARY KEY,
								user_id VARCHAR NOT NULL, 	
								origin_url VARCHAR NOT NULL, 
								short_url VARCHAR NOT NULL UNIQUE,
                                is_deleted BOOLEAN NOT NULL DEFAULT FALSE
					);`
	// отправляем команду на создание бд
	res, err := db.ExecContext(ctx, sqlCreateDB)
	log.Println("Create table", err, res)

	return nil
}
