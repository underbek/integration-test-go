package step_4_get_user_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-testfixtures/testfixtures/v3"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
=== RUN   TestGetUser
2022/06/12 17:53:01 Starting container id: 1f29d6f8e2e1 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 17:53:02 Waiting for container id 1f29d6f8e2e1 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 17:53:02 Container is ready id: 1f29d6f8e2e1 image: docker.io/testcontainers/ryuk:0.3.3
2022/06/12 17:53:02 Starting container id: f6c8771b6f97 image: postgres:11.5
2022/06/12 17:53:03 Waiting for container id f6c8771b6f97 image: postgres:11.5
2022/06/12 17:53:04 Container is ready id: f6c8771b6f97 image: postgres:11.5
Host: localhost 52765
2022/06/12 17:53:04 OK    20220612163022_create_users.sql
2022/06/12 17:53:04 goose: no migrations to run. current version: 20220612163022
error sql: no rows in result set
    integration_test.go:69:
        	Error Trace:	integration_test.go:69
        	Error:      	Not equal:
        	            	expected: 200
        	            	actual  : 500
        	Test:       	TestGetUser
    integration_test.go:74:
        	Error Trace:	integration_test.go:74
        	Error:      	Received unexpected error:
        	            	EOF
        	Test:       	TestGetUser
    integration_test.go:78:
        	Error Trace:	integration_test.go:78
        	Error:      	Not equal:
        	            	expected: 1
        	            	actual  : 0
        	Test:       	TestGetUser
    integration_test.go:79:
        	Error Trace:	integration_test.go:79
        	Error:      	Not equal:
        	            	expected: "test_name"
        	            	actual  : ""

        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test_name
        	            	+
        	Test:       	TestGetUser
--- FAIL: TestGetUser (3.46s)


Expected :test_name
Actual   :
<Click to see difference>


FAIL
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
	///

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

	// use httptest
	srv := httptest.NewServer(server.New("", h).Router)

	res, err := srv.Client().Get(srv.URL + "/users/1")
	require.NoError(t, err)

	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	// check response
	response := api.GetUserResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)

	// so we will check each field separately
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "test_name", response.Name)
	assert.Equal(t, "0", response.Balance.String())
}
