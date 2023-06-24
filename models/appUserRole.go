package models

import (
	"github.com/uptrace/bun"
)

type AppUserRoles struct {
	bun.BaseModel `bun:"table:app_user_role,alias:url"`
	ID            int64  `bun:"id,pk,autoincrement"`
	Name          string `bun:"name,notnull,type:varchar(64)"`
}
