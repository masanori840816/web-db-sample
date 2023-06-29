package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"time"

	auth "github.com/web-db-sample/auth"
	db "github.com/web-db-sample/db"
	dto "github.com/web-db-sample/dto"
)

type templateHandler struct {
	once     sync.Once
	templ    *template.Template
	settings *AppSettings
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		handleMainPageRequest(w, r, t)
	case "/pages/signin":
		handleSigninPageRequest(w, r, t)
	default:
		w.WriteHeader(404)
	}
}

func main() {
	settings, err := NewAppSettings()
	if err != nil {
		log.Panicln(err.Error())
	}
	dbCtx := db.NewBookshelfContext(settings.ConnectionString)
	http.Handle("/js/", http.FileServer(http.Dir("templates")))
	http.Handle("/img/", http.FileServer(http.Dir("templates")))
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		result := &dto.SigninResult{}
		signinResult, err := auth.Signin(w, r, dbCtx)
		if err != nil {
			log.Println(err.Error())
			result.Succeeded = false
			result.ErrorMessage = "Sever Error"
			resultJSON, _ := json.Marshal(result)
			w.Write(resultJSON)
			return
		}
		if signinResult {
			result.Succeeded = true
			result.ErrorMessage = ""
			result.NextURL = settings.BaseURL
		} else {
			result.Succeeded = false
			result.ErrorMessage = "Invalid user name or password"
		}
		resultJSON, _ := json.Marshal(result)
		w.Write(resultJSON)
	})
	http.HandleFunc("/signout", func(w http.ResponseWriter, r *http.Request) {
		auth.Signout(w)
		result := &dto.SigninResult{
			Succeeded: true,
			NextURL:   fmt.Sprintf("%spages/signin", settings.BaseURL),
		}
		resultJson, _ := json.Marshal(result)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resultJson)
	})
	http.Handle("/", &templateHandler{settings: &settings})
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", settings.Host, settings.Port), nil))
}
func handleSigninPageRequest(w http.ResponseWriter, r *http.Request, t *templateHandler) {
	t.once.Do(func() {
		// "Must()" wraps "ParseFiles()" results, so I can put it into "templateHandler.templ" directly
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", "signin.html")))
	})
	t.templ.Execute(w, t.settings.BaseURL)
}
func handleMainPageRequest(w http.ResponseWriter, r *http.Request, t *templateHandler) {
	authenticated, usedID := auth.VerifyToken(w, r)
	if !authenticated {
		expiration := time.Now()
		expiration = expiration.Add(time.Hour)
		cookie := http.Cookie{Name: "RedirectURL", Value: t.settings.BaseURL, Expires: expiration, HttpOnly: true}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, fmt.Sprintf("%spages/signin", t.settings.BaseURL), http.StatusFound)
		return
	}
	log.Printf("Sign in user ID: %d", usedID)
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))
	})
	t.templ.Execute(w, t.settings.BaseURL)
}
