package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	db "github.com/web-db-sample/db"
	dto "github.com/web-db-sample/dto"
)

const secretKey = "kwTL6Nnm.4gbTPBCU_6kveHEZg"

func Signin(w http.ResponseWriter, r *http.Request, dbCtx *db.BookshelfContext) (bool, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return false, err
	}
	signinValue := &dto.SigninValues{}
	err = json.Unmarshal(body, &signinValue)
	if err != nil {
		return false, err
	}
	ctx := context.Background()
	result, userID, err := dbCtx.Users.Signin(&ctx, *signinValue)
	if err != nil {
		return false, err
	}
	if !result {
		return false, nil
	}
	token, err := generateToken(userID)
	if err != nil {
		return false, err
	}
	expiration := time.Now()
	expiration = expiration.AddDate(0, 0, 1)
	cookie := http.Cookie{Name: "AuthSample", Value: token, Expires: expiration, HttpOnly: true}
	http.SetCookie(w, &cookie)
	return result, nil
}
func Signout(w http.ResponseWriter) {
	cookie := http.Cookie{Name: "AuthSample", Value: "", Expires: time.Unix(0, 0), HttpOnly: true}
	http.SetCookie(w, &cookie)
}
func VerifyToken(w http.ResponseWriter, r *http.Request) (bool, int64) {
	tokenString := getToken(r)
	if len(tokenString) <= 0 {
		return false, -1
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("UNEXPECTED SIGNING METHOD: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return false, -1
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, int64(claims["userid"].(float64))
	}
	return false, -1
}
func generateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"userid": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Add sign
	return token.SignedString([]byte(secretKey))
}
func getToken(r *http.Request) string {
	for _, c := range r.Cookies() {
		if c.Name != "AuthSample" {
			continue
		}
		// remove "AuthSample=" from the cookie value
		result := c.String()
		// If the client has a valid cookie, it can retrieve the value.
		return strings.Replace(result, "AuthSample=", "", 1)
	}
	return ""
}
