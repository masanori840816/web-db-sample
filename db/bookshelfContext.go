package db

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type BookshelfContext struct {
	db    *bun.DB
	Users *Users
}

func NewBookshelfContext() *BookshelfContext {
	result := BookshelfContext{}
	dsn := "postgresql://postgres:postgres@localhost:5432/book_shelf?sslmode=disable"
	result.db = bun.NewDB(
		sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn))),
		pgdialect.New(),
	)
	result.Users = NewUsers(result.db)
	return &result
}
