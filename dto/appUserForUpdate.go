package dto

import (
	"fmt"
)

type AppUserForUpdate struct {
	ID       int64  `json:"id"`
	RoleID   int64  `json:"role_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u AppUserForUpdate) Validate() error {
	if u.RoleID <= 0 {
		return fmt.Errorf("ROLE IS REQUIRED")
	}
	if len(u.Name) <= 0 {
		return fmt.Errorf("NAME IS REQUIRED")
	}
	if len(u.Password) <= 0 {
		return fmt.Errorf("PASSWORD IS REQUIRED")
	}
	return nil
}
