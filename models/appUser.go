package models

import (
	"time"

	"github.com/uptrace/bun"
)

type AppUsers struct {
	bun.BaseModel  `bun:"table:app_user,alias:usr"`
	ID             int64         `bun:"id,pk,autoincrement"`
	RoleID         int64         `bun:"app_user_role_id,notnull,type:bigint"`
	Name           string        `bun:"name,notnull,type:varchar(64)"`
	Password       string        `bun:"password,notnull,type:text"`
	LastUpdateDate time.Time     `bun:"last_update_date,notnull,type:timestamp with time zone,default:CURRENT_TIMESTAMP"`
	Role           *AppUserRoles `bun:"rel:has-one,join:app_user_role_id=id"`
}

func NewAppUsers(roleID int64, name string, hashedPassword string) *AppUsers {
	return &AppUsers{
		RoleID:   roleID,
		Name:     name,
		Password: hashedPassword,
	}
}
func (u *AppUsers) Update(roleID int64, name string, hashedPassword string) {
	u.RoleID = roleID
	u.Name = name
	u.Password = hashedPassword
	u.LastUpdateDate = time.Now()
}
