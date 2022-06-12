package main

import (
	"log"

	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/server"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
)

const (
	addr  = ":8080"
	dbDsn = "host=localhost port=5432 user=user password=user dbname=user sslmode=disable"
)

func main() {
	repo, err := storage.New(dbDsn)
	if err != nil {
		log.Fatal(err)
	}
	useCase := use_case.New(repo)
	h := handler.New(useCase)
	srv := server.New(addr, h)
	log.Fatal(srv.Serve())
}
