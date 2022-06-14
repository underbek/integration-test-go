package step_2

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"

	"github.com/AndreyAndreevich/articles/integration_tests/step_2"
	"github.com/AndreyAndreevich/articles/user_service/api"
	"github.com/AndreyAndreevich/articles/user_service/billing"
	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/AndreyAndreevich/articles/user_service/migrate"
	"github.com/AndreyAndreevich/articles/user_service/server"
	"github.com/AndreyAndreevich/articles/user_service/storage"
	"github.com/AndreyAndreevich/articles/user_service/use_case"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/suite"
)

const billingAddr = "http://localhost:8085"

type TestSuite struct {
	suite.Suite
	psqlContainer *step_2.PostgreSQLContainer
	server        *httptest.Server
}

func (s *TestSuite) SetupSuite() {
	// create db container
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := step_2.NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	s.psqlContainer = psqlContainer
	//

	// run migrations
	err = migrate.Migrate(psqlContainer.GetDSN(), migrate.Migrations)
	s.Require().NoError(err)
	//

	// copy from main
	repo, err := storage.New(psqlContainer.GetDSN())
	if err != nil {
		log.Fatal(err)
	}

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
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.psqlContainer.Terminate(ctx))

	s.server.Close()

	httpmock.DeactivateAndReset()
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
	requestBody := `{"name": "test_name"}`

	res, err := s.server.Client().Post(s.server.URL+"/users", "", bytes.NewBufferString(requestBody))
	s.Require().NoError(err)

	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode)

	// check response
	response := api.CreateUserResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	s.Require().NoError(err)

	// id maybe any
	// so we will check each field separately
	s.Assert().Equal("test_name", response.Name)
	s.Assert().Equal("0", response.Balance.String())
}

func (s *TestSuite) TestGetUser() {
	// create fixtures
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("../step_5/fixtures/storage"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	//

	res, err := s.server.Client().Get(s.server.URL + "/users/1")
	s.Require().NoError(err)

	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode)

	// check response
	response := api.GetUserResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	s.Require().NoError(err)

	// so we will check each field separately
	s.Assert().Equal(1, response.ID)
	s.Assert().Equal("test_name", response.Name)
	s.Assert().Equal("0", response.Balance.String())
}

func (s *TestSuite) TestDepositBalance() {
	// create fixtures
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("../step_5/fixtures/storage"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	//

	// mock http call
	httpmock.RegisterResponder(
		http.MethodPost,
		billingAddr+"/deposit",
		httpmock.NewStringResponder(http.StatusOK, ""),
	)

	requestBody := `{"id": 1, "amount": "100"}`

	res, err := s.server.Client().Post(s.server.URL+"/users/deposit", "", bytes.NewBufferString(requestBody))
	s.Require().NoError(err)

	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode)

	// check response
	response := api.GetUserResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	s.Require().NoError(err)

	// so we will check each field separately
	s.Assert().Equal(1, response.ID)
	s.Assert().Equal("test_name", response.Name)
	s.Assert().Equal("100", response.Balance.String())
}
