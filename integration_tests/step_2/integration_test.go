package step_2

import (
	"bytes"
	"context"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/server"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
	"github.com/stretchr/testify/assert"
)

/*
=== RUN   TestCreateUser
2022/06/12 17:41:25 Starting container id: 20ccb108f948 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 17:41:25 Waiting for container id 20ccb108f948 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 17:41:25 Container is ready id: 20ccb108f948 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 17:41:25 Starting container id: 81d293f7ee95 image: postgres:11.5
2022/06/12 17:41:26 Waiting for container id 81d293f7ee95 image: postgres:11.5
2022/06/12 17:41:28 Container is ready id: 81d293f7ee95 image: postgres:11.5
Host: localhost 49770
error pq: relation "users" does not exist
--- PASS: TestCreateUser (3.55s)
PASS
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
	srv := httptest.NewServer(server.New("", h).Router)

	_, err = srv.Client().Post(srv.URL+"/users", "", bytes.NewBufferString(requestBody))
	assert.NoError(t, err)
}
