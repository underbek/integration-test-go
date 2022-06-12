package step_2

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AndreyAndreevich/articles/integration_tests/step_2"
	"github.com/AndreyAndreevich/articles/user_service/api"
	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/migrate"
	"github.com/AndreyAndreevich/articles/user_service/server"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
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

	psqlContainer, err := step_2.NewPostgreSQLContainer(ctx)
	defer psqlContainer.Terminate(context.Background())
	assert.NoError(t, err)
	//

	// run migrations
	err = migrate.Migrate(psqlContainer.GetDSN(), migrate.Migrations)
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

	// create fixtures
	db, err := sql.Open("postgres", psqlContainer.GetDSN())
	assert.NoError(t, err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("fixtures/storage"),
	)
	assert.NoError(t, err)
	assert.NoError(t, fixtures.Load())
	//

	// use httptest
	srv := httptest.NewServer(server.New("", h).Router)

	res, err := srv.Client().Get(srv.URL + "/users/1")
	assert.NoError(t, err)

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	// check response
	response := api.GetUserResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)

	// so we will check each field separately
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "test_name", response.Name)
	assert.Equal(t, "0", response.Balance.String())
}
