package dto

import (
	"time"
)

type AppUserForView struct {
	ID             int64     `bun:"id" json:"id"`
	RoleID         int64     `bun:"roleId" json:"roleId"`
	RoleName       string    `bun:"roleName" json:"roleName"`
	Name           string    `bun:"name" json:"name"`
	LastUpdateDate time.Time `bun:"lastUpdateDate" json:"lastUpdateDate"`
}
