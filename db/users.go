package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"

	dto "github.com/web-db-sample/dto"
	hash "github.com/web-db-sample/hash"
	models "github.com/web-db-sample/models"
)

type Users struct {
	db *bun.DB
}

func NewUsers(database *bun.DB) *Users {
	return &Users{
		db: database,
	}
}
func (u Users) CraeteUser(ctx *context.Context, user dto.AppUserForUpdate) error {
	// Use tx instead of db to enable transactions
	tx, err := u.db.BeginTx(*ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	// Make sure the Role ID is registered
	exists, err := tx.NewSelect().Model(new(models.AppUserRoles)).
		Where("id=?", user.RoleID).Exists(*ctx)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("INVALID ROLE ID:%d", user.RoleID)
	}
	// Make sure the user name is unique
	exists, err = tx.NewSelect().Model(new(models.AppUsers)).
		Where("name=?", user.Name).Exists(*ctx)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("USER NAME IS ALREADY EXITS:%s", user.Name)
	}

	hashedPassword, err := hash.GeneratePasswordHash(user.Password)
	if err != nil {
		return err
	}
	// Insert new user
	newUser := models.NewAppUsers(user.RoleID, user.Name, hashedPassword)
	_, err = tx.NewInsert().Model(newUser).Exec(*ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (u Users) GetUser(ctx *context.Context, userId int64) (*models.AppUsers, error) {
	user := new(models.AppUsers)
	err := u.db.NewSelect().
		Model(user).
		Relation("Role").
		Where("usr.id=?", userId).
		Limit(1).
		Scan(*ctx)
	if err != nil {
		// Ignore no rows error and return nil
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func (u Users) UpdateUser(ctx *context.Context, user dto.AppUserForUpdate) error {
	tx, err := u.db.BeginTx(*ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	existedUser, err := u.GetUser(ctx, user.ID)
	if err != nil {
		return err
	}
	hashedPassword, err := hash.GeneratePasswordHash(user.Password)
	if err != nil {
		return err
	}
	existedUser.Update(user.RoleID, user.Name, hashedPassword)
	res, err := tx.NewUpdate().Model(existedUser).WherePK().Exec(*ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if r != 1 {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (u Users) GetAllUsersForView(ctx *context.Context) ([]dto.AppUserForView, error) {
	/*roles := make([]models.AppUserRoles, 0)
	err := u.db.NewSelect().
		Model(roles).
		Scan(*ctx)
	if err != nil {
		return nil, err
	}*/
	results := make([]dto.AppUserForView, 0)

	err := u.db.NewRaw(
		`SELECT usr.id AS "id", url.id AS "roleId", usr.name AS "name", url.name AS "roleName",
		usr.last_update_date AS "lastUpdateDate" FROM app_user usr
		JOIN app_user_role url ON usr.app_user_role_id = url.id
		`).Scan(*ctx, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (u Users) Signin(ctx *context.Context, value dto.SigninValues) (bool, int64, error) {
	user := new(models.AppUsers)
	err := u.db.NewSelect().
		Model(user).
		Where("name=?", value.UserName).
		Limit(1).
		Scan(*ctx)
	if err != nil {
		// ignore no rows error
		if err != sql.ErrNoRows {
			return false, -1, err
		}
	}
	result, err := hash.VerifyPassword(value.Password, user.Password)
	if err != nil {
		return false, -1, err
	}
	if result {
		return result, user.ID, nil
	}
	return result, -1, nil
}
