package main

import (
	"avito-tech/internal/app/core"
	pullrequest "avito-tech/internal/app/pull_request"
	"avito-tech/internal/app/routing"
	"avito-tech/internal/app/team"
	"avito-tech/internal/app/user"
	"avito-tech/internal/db"
	"context"
	"fmt"
	"net/http"
)

const port = ":8080"

func main() {
	ctx := context.Background()

	db, err := db.CreateDB(ctx)

	if err != nil {
		fmt.Println("Failed to create DB")
		return
	}

	user := user.NewUser(user.NewUserRepo(db))
	team := team.NewTeam(team.NewTeamRepo(db))
	pull_request := pullrequest.NewPullRequest(pullrequest.NewRepo(db))

	service := core.NewService(team, user, pull_request)

	server := routing.NewServer(service)

	router := routing.NewRouter(server)

	if err := http.ListenAndServe(port, router); err != nil {
		fmt.Println("Failed to Run server")
	}
}
