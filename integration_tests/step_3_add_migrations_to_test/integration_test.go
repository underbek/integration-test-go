package step_3_add_migrations_to_test

import (
	"bytes"
	"context"
	"encoding/json"
	step_2_1 "github.com/AndreyAndreevich/articles/integration_tests/step_2_1_improved_psql_container"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AndreyAndreevich/articles/user_service/api"
	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/migrate"
	"github.com/AndreyAndreevich/articles/user_service/server"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
=== RUN   TestCreateUser
2022/06/12 16:44:46 Starting container id: c94311c9ea05 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 16:44:46 Waiting for container id c94311c9ea05 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 16:44:46 Container is ready id: c94311c9ea05 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 16:44:46 Starting container id: fe9c994b3e7b image: postgres:11.5
2022/06/12 16:44:47 Waiting for container id fe9c994b3e7b image: postgres:11.5
2022/06/12 16:44:49 Container is ready id: fe9c994b3e7b image: postgres:11.5
Host: localhost 51709
2022/06/12 16:44:49 OK    20220612163022_create_users.sql
2022/06/12 16:44:49 goose: no migrations to run. current version: 20220612163022
--- PASS: TestCreateUser (3.51s)
PASS
*/
func TestCreateUser(t *testing.T) {
	// create db container
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := step_2_1.NewPostgreSQLContainer(ctx)
	defer psqlContainer.Terminate(context.Background())
	require.NoError(t, err)
	//

	// run migrations
	err = migrate.Migrate(psqlContainer.GetDSN(), migrate.Migrations)
	require.NoError(t, err)
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

	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	// check response
	response := api.CreateUserResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)

	// id maybe any
	// so we will check each field separately
	assert.Equal(t, "test_name", response.Name)
	assert.Equal(t, "0", response.Balance.String())
}
