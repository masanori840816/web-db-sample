package models

import (
	"github.com/uptrace/bun"
)

type Genres struct {
	bun.BaseModel `bun:"table:app_user_role,alias:gnr"`
	ID            int64  `bun:"id,pk,autoincrement"`
	Name          string `bun:"name,notnull,type:varchar(64)"`
}
