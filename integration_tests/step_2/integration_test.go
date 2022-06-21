package step_2

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/server"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
	"github.com/stretchr/testify/require"
)

/*
=== RUN   TestCreateUser
2022/06/19 19:38:23 Starting container id: cefb7fcdbde3 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/19 19:38:23 Waiting for container id cefb7fcdbde3 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/19 19:38:23 Container is ready id: cefb7fcdbde3 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/19 19:38:23 Starting container id: 115811151b36 image: postgres:11.5
2022/06/19 19:38:24 Waiting for container id 115811151b36 image: postgres:11.5
2022/06/19 19:38:26 Container is ready id: 115811151b36 image: postgres:11.5
Host: localhost 55297
error pq: relation "users" does not exist
    integration_test.go:59:
        	Error Trace:	integration_test.go:59
        	Error:      	Not equal:
        	            	expected: 200
        	            	actual  : 500
        	Test:       	TestCreateUser
--- FAIL: TestCreateUser (3.27s)


Expected :200
Actual   :500
<Click to see difference>


FAIL
*/
func TestCreateUser(t *testing.T) {
	// create db container
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := NewPostgreSQLContainer(ctx)
	require.NoError(t, err)
	defer psqlContainer.Terminate(context.Background())
	//

	// copy from main
	repo, err := storage.New(psqlContainer.GetDSN())
	require.NoError(t, err)
	useCase := use_case.New(repo, nil)
	h := handler.New(useCase)
	///

	requestBody := `{"name": "test_name"}`

	// use httptest
	srv := httptest.NewServer(server.New("", h).Router)

	res, err := srv.Client().Post(srv.URL+"/users", "", bytes.NewBufferString(requestBody))
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
}
