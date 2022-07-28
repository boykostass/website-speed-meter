package services

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/julienschmidt/httprouter"
)

// RepQuery - Интерфейсы для работы с репозиторием
type RepQuery interface {
	FindOne(ctx context.Context, name string) ([]SiteInfo, error)
	FindAll(ctx context.Context) ([]SiteInfo, error)
}

// Database - Интерфейсы для работы с базой данных
type Database interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

// Handler - Интерфейс для регистрации хэндлера
type Handler interface {
	Register(router *httprouter.Router)
}
