package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Books struct {
	bun.BaseModel  `bun:"table:book,alias:bok"`
	ID             int64      `bun:"id,pk,autoincrement"`
	Name           string     `bun:"name,notnull,type:varchar(512)"`
	GenreID        int64      `bun:"genre_id,notnull,type:bigint"`
	AuthorID       int64      `bun:"author_id,notnull,type:bigint"`
	LanguageID     int64      `bun:"language_id,notnull,type:bigint"`
	LastUpdateDate time.Time  `bun:"last_update_date,notnull,type:timestamp with time zone,default:CURRENT_TIMESTAMP"`
	Genre          *Genres    `bun:"rel:has-one,join:genre_id=id"`
	Author         *Authors   `bun:"rel:has-one,join:author_id=id"`
	Language       *Languages `bun:"rel:has-one,join:language_id=id"`
}
