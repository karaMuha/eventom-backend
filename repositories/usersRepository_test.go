package repositories

import (
	"context"
	"database/sql"
	"eventom-backend/testutils"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

type UsersRepositoryTestSuite struct {
	suite.Suite
	ctx             context.Context
	pgContainer     *testutils.PostgresContainer
	dbHandler       *sql.DB
	usersRepository UsersRepositoryInterface
}

func TestUsersRepositorySuite(t *testing.T) {
	suite.Run(t, &UsersRepositoryTestSuite{})
}

func (suite *UsersRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testutils.CreatePostgresContainer(suite.ctx)

	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer
	dbHandler, err := sql.Open("postgres", suite.pgContainer.ConnectionString)

	if err != nil {
		log.Fatal(err)
	}

	err = dbHandler.Ping()

	if err != nil {
		log.Fatal(err)
	}

	suite.dbHandler = dbHandler

	suite.usersRepository = NewUsersRepository(dbHandler)
}

func (suite *UsersRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	// delete content of users table before every test to avoid dependencies and side effects between tests

	query := `
		DELETE FROM
			users`
	_, err := suite.dbHandler.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

func (suite *UsersRepositoryTestSuite) TearDownSuite() {
	err := suite.pgContainer.Terminate(suite.ctx)
	if err != nil {
		log.Fatalf("Error while terminating postgres container: %s", err)
	}
}

func (suite *UsersRepositoryTestSuite) TestSignupUserSuccess() {
	err := suite.usersRepository.QuerySignupUser("test@test.com", "Test123")
	assert.Nil(suite.T(), err)
}

func (suite *UsersRepositoryTestSuite) TestSignupUserFailUniqueConstraintEmail() {
	err := suite.usersRepository.QuerySignupUser("test@test.com", "Test123")
	assert.Nil(suite.T(), err)

	err = suite.usersRepository.QuerySignupUser("test@test.com", "Test123")
	assert.NotNil(suite.T(), err)
}

func (suite *UsersRepositoryTestSuite) TestGetUserFailNoUser() {
	user, err := suite.usersRepository.QueryGetUser("test@test.com")
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), user)
}

func (suite *UsersRepositoryTestSuite) TestGetUserSuccess() {
	err := suite.usersRepository.QuerySignupUser("test@test.com", "Test123")
	assert.Nil(suite.T(), err)

	user, err := suite.usersRepository.QueryGetUser("test@test.com")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.NotNil(suite.T(), user.ID)
	assert.Equal(suite.T(), "test@test.com", user.Email)
}
