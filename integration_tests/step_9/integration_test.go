package step_9

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	step2 "github.com/AndreyAndreevich/articles/integration_tests/step_2_1_improved_psql_container"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AndreyAndreevich/articles/user_service/api"
	"github.com/AndreyAndreevich/articles/user_service/billing"
	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/migrate"
	"github.com/AndreyAndreevich/articles/user_service/server"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

const billingAddr = "http://localhost:8085"

type TestSuite struct {
	suite.Suite
	psqlContainer *step2.PostgreSQLContainer
	server        *httptest.Server
	loader        *FixtureLoader
}

func (s *TestSuite) SetupSuite() {
	// create db container
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := step2.NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	s.psqlContainer = psqlContainer
	//

	// run migrations
	err = migrate.Migrate(psqlContainer.GetDSN(), migrate.Migrations)
	s.Require().NoError(err)
	//

	// copy from main
	repo, err := storage.New(psqlContainer.GetDSN())
	s.Require().NoError(err)

	//mock client
	mockClient := &http.Client{}
	httpmock.ActivateNonDefault(mockClient)

	//added billing client
	billingClient := billing.New(mockClient, billingAddr)
	useCase := use_case.New(repo, billingClient)
	h := handler.New(useCase)
	///

	// use httptest
	s.server = httptest.NewServer(server.New("", h).Router)
	//

	// create fixture loader
	s.loader = NewFixtureLoader(s.T(), Fixtures)
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.psqlContainer.Terminate(ctx))

	s.server.Close()

	httpmock.DeactivateAndReset()
}

// create fixtures before each test
func (s *TestSuite) SetupTest() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("../step_5/fixtures/storage"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
}

/*
--- PASS: TestSuite_Run (4.07s)
=== RUN   TestSuite_Run/TestCreateUser
    --- PASS: TestSuite_Run/TestCreateUser (0.01s)
=== RUN   TestSuite_Run/TestDepositBalance
    --- PASS: TestSuite_Run/TestDepositBalance (0.06s)
=== RUN   TestSuite_Run/TestGetUser
    --- PASS: TestSuite_Run/TestGetUser (0.07s)
PASS
*/
func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestCreateUser() {
	requestBody := s.loader.LoadString("fixtures/api/create_user_request.json")

	res, err := s.server.Client().Post(s.server.URL+"/users", "", bytes.NewBufferString(requestBody))
	s.Require().NoError(err)

	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode)

	// check response
	response := api.CreateUserResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	s.Require().NoError(err)

	expected := s.loader.LoadTemplate("fixtures/api/create_user_response.json.temp", map[string]interface{}{
		"id": response.ID,
	})

	JSONEq(s.T(), expected, response)
}

func (s *TestSuite) TestGetUser() {
	res, err := s.server.Client().Get(s.server.URL + "/users/1")
	s.Require().NoError(err)

	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode)

	// check response
	expected := s.loader.LoadString("fixtures/api/get_user_response.json")

	JSONEq(s.T(), expected, res.Body)
}

func (s *TestSuite) TestDepositBalance() {
	// mock http call
	httpmock.RegisterResponder(
		http.MethodPost,
		billingAddr+"/deposit",
		httpmock.NewStringResponder(http.StatusOK, ""),
	)
	//

	requestBody := s.loader.LoadString("fixtures/api/deposit_user_request.json")

	res, err := s.server.Client().Post(s.server.URL+"/users/deposit", "", bytes.NewBufferString(requestBody))
	s.Require().NoError(err)

	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode)

	// check response
	expected := s.loader.LoadString("fixtures/api/deposit_user_response.json")

	JSONEq(s.T(), expected, res.Body)
}
