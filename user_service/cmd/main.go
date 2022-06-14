package main

import (
	"log"
	"net/http"

	"github.com/AndreyAndreevich/articles/user_service/billing"
	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/server"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
)

const (
	addr        = ":8080"
	dbDsn       = "host=localhost port=5432 user=user password=user dbname=user sslmode=disable"
	billingAddr = "localhost:8081"
)

func main() {
	repo, err := storage.New(dbDsn)
	if err != nil {
		log.Fatal(err)
	}

	billingClient := billing.New(http.DefaultClient, billingAddr)
	useCase := use_case.New(repo, billingClient)
	h := handler.New(useCase)
	srv := server.New(addr, h)
	log.Fatal(srv.Serve())
}
