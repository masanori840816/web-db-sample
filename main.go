package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	auth "github.com/web-db-sample/auth"
	db "github.com/web-db-sample/db"
	dto "github.com/web-db-sample/dto"
)

const baseURL = "http://localhost:8081/"

type templateHandler struct {
	once  sync.Once
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
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
	dbCtx := db.NewBookshelfContext()
	ctx := context.Background()

	sampleUser := dto.AppUserForUpdate{
		RoleID:   1,
		Name:     "Masa",
		Password: "Sample",
	}
	err := sampleUser.Validate()
	if err != nil {
		log.Panicln(err.Error())
	}
	u, err := dbCtx.Users.GetUser(&ctx, 1)
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Println(u)
	http.Handle("/js/", http.FileServer(http.Dir("templates")))
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {

	})
	http.Handle("/", &templateHandler{})
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
func handleSigninPageRequest(w http.ResponseWriter, r *http.Request, t *templateHandler) {
	t.once.Do(func() {
		// "Must()" wraps "ParseFiles()" results, so I can put it into "templateHandler.templ" directly
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", "signin.html")))
	})
	t.templ.Execute(w, baseURL)
}
func handleMainPageRequest(w http.ResponseWriter, r *http.Request, t *templateHandler) {
	authenticated, usedID := auth.VerifyToken(w, r)
	if !authenticated {
		http.Redirect(w, r, fmt.Sprintf("%spages/signin", baseURL), http.StatusFound)
		return
	}
	log.Println(authenticated)
	log.Printf("Sign in user ID: %d", usedID)
	t.once.Do(func() {
		// "Must()" wraps "ParseFiles()" results, so I can put it into "templateHandler.templ" directly
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))
	})
	t.templ.Execute(w, baseURL)
}
