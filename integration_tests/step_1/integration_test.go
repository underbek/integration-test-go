package step_1

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
	"gotest.tools/assert"
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
	useCase := use_case.New(repo)
	h := handler.New(useCase)
	///

	requestBody := `{"name": "test_name"}`

	// use httptest
	request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()
	h.CreateUser(w, request)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
