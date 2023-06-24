package db

import "github.com/uptrace/bun"

type Books struct {
	db *bun.DB
}

func NewBooks(database *bun.DB) Books {
	return Books{
		db: database,
	}
}
