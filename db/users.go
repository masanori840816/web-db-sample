package db

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/big"

	"crypto/rand"
	"crypto/sha512"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/pbkdf2"

	dto "github.com/web-db-sample/dto"
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
	salt, err := generateRandomSalt(128 / 8)
	if err != nil {
		return err
	}
	// Get base 64 encoded Hasu value to save the password
	key := pbkdf2.Key([]byte(user.Password), salt, 100_000, 256/8, sha512.New)
	hashedPassword := base64.StdEncoding.EncodeToString(key)
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

// Generate a salt value
func generateRandomSalt(length int) ([]byte, error) {
	results := make([]byte, length)
	for i := 0; i < length; i++ {
		salt, err := rand.Int(rand.Reader, big.NewInt(255))
		if err != nil {
			return nil, err
		}
		results[i] = byte(salt.Int64())
	}
	return results, nil
}
