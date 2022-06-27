package step_5_add_testfixtures

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	step2 "github.com/AndreyAndreevich/articles/integration_tests/step_2_1_improved_psql_container"
	"github.com/AndreyAndreevich/articles/user_service/api"
	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/migrate"
	"github.com/AndreyAndreevich/articles/user_service/server"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
=== RUN   TestGetUser
2022/06/12 17:52:06 Starting container id: b43bd5af71f3 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 17:52:07 Waiting for container id b43bd5af71f3 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 17:52:07 Container is ready id: b43bd5af71f3 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 17:52:07 Starting container id: 2d7713d56850 image: postgres:11.5
2022/06/12 17:52:07 Waiting for container id 2d7713d56850 image: postgres:11.5
2022/06/12 17:52:09 Container is ready id: 2d7713d56850 image: postgres:11.5
Host: localhost 52521
2022/06/12 17:52:09 OK    20220612163022_create_users.sql
2022/06/12 17:52:09 goose: no migrations to run. current version: 20220612163022
--- PASS: TestGetUser (3.47s)
PASS
*/
func TestGetUser(t *testing.T) {
	// create db container
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := step2.NewPostgreSQLContainer(ctx)
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

	// use httptest
	srv := httptest.NewServer(server.New("", h).Router)

	// create fixtures
	db, err := sql.Open("postgres", psqlContainer.GetDSN())
	require.NoError(t, err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("fixtures/storage"),
	)
	require.NoError(t, err)
	require.NoError(t, fixtures.Load())
	//

	res, err := srv.Client().Get(srv.URL + "/users/1")
	require.NoError(t, err)

	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	// check response
	response := api.GetUserResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "test_name", response.Name)
	assert.Equal(t, "0", response.Balance.String())
}
