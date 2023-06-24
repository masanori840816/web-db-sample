package main

import (
	"context"
	"log"

	db "github.com/web-db-sample/db"
	dto "github.com/web-db-sample/dto"
)

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
}
