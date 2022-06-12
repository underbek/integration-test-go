package step_2

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
	"github.com/stretchr/testify/assert"
)

/*
=== RUN   TestCreateUser
2022/06/12 16:38:32 Starting container id: 69745226ac9c image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 16:38:33 Waiting for container id 69745226ac9c image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 16:38:33 Container is ready id: 69745226ac9c image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 16:38:33 Starting container id: 9be46e6ca54f image: postgres:11.5
2022/06/12 16:38:33 Waiting for container id 9be46e6ca54f image: postgres:11.5
2022/06/12 16:38:36 Container is ready id: 9be46e6ca54f image: postgres:11.5
Host: localhost 50153
error pq: relation "users" does not exist
    integration_test.go:51:
        	Error Trace:	integration_test.go:51
        	Error:      	Not equal:
        	            	expected: 200
        	            	actual  : 500
        	Test:       	TestCreateUser
--- FAIL: TestCreateUser (4.15s)


Expected :200
Actual   :500
*/
func TestCreateUser(t *testing.T) {
	// create db container
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := NewPostgreSQLContainer(ctx)
	defer psqlContainer.Terminate(context.Background())
	assert.NoError(t, err)
	//

	// copy from main
	repo, err := storage.New(psqlContainer.GetDSN())
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
