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

func NewBookshelfContext(connectionString string) *BookshelfContext {
	result := BookshelfContext{}
	result.db = bun.NewDB(
		sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionString))),
		pgdialect.New(),
	)
	result.Users = NewUsers(result.db)
	return &result
}
