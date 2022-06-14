package step_1

import (
	"bytes"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/server"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
	"github.com/stretchr/testify/assert"
)

const (
	dbDsn = "host=localhost port=5432 user=user password=user dbname=user sslmode=disable"
)

/*
=== RUN   TestCreateUser
2022/06/12 16:27:07 dial tcp [::1]:5432: connect: connection refused
*/
func TestCreateUser(t *testing.T) {

	// copy from main
	repo, err := storage.New(dbDsn)
	if err != nil {
		log.Fatal(err)
	}
	useCase := use_case.New(repo, nil)
	h := handler.New(useCase)
	///

	requestBody := `{"name": "test_name"}`

	// use httptest
	srv := httptest.NewServer(server.New("", h).Router)

	_, err = srv.Client().Post(srv.URL+"/users", "", bytes.NewBufferString(requestBody))
	assert.NoError(t, err)
}
