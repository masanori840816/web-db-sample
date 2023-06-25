package main

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	db "github.com/web-db-sample/db"
	dto "github.com/web-db-sample/dto"
)

type templateHandler struct {
	once      sync.Once
	filename  string
	templ     *template.Template
	serverUrl string
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// "sync.Once" executes only one time.
	t.once.Do(func() {
		// "Must()" wraps "ParseFiles()" results, so I can put it into "templateHandler.templ" directly
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, t.serverUrl)
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
	//err = dbCtx.Users.CraeteUser(&ctx, sampleUser)
	/*
		users, err := dbCtx.Users.GetAllUsersForView(&ctx)
		if err != nil {
			log.Panicln(err.Error())
		}
		j, _ := json.Marshal(users)
		log.Println(string(j))*/
	u, err := dbCtx.Users.GetUser(&ctx, 1)
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Println(u)
	http.Handle("/js/", http.FileServer(http.Dir("templates")))
	http.Handle("/", &templateHandler{filename: "signin.html", serverUrl: "http://localhost:8081"})
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
