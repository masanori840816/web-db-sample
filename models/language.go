package models

import (
	"github.com/uptrace/bun"
)

type Languages struct {
	bun.BaseModel `bun:"table:language,alias:lng"`
	ID            int64  `bun:"id,pk,type:integer"`
	Name          string `bun:"name,notnull,type:varchar(64)"`
}
