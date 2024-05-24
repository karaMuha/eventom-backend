package repositories

import (
	"context"
	"eventom-backend/testutils"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

type UsersRepositoryTestSuite struct {
	suite.Suite
	ctx             context.Context
	pgContainer     *testutils.PostgresContainer
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

	suite.usersRepository = NewUsersRepository(pgContainer.DB)
}

func (suite *UsersRepositoryTestSuite) AfterTest(suiteName, testName string) {
	// clear users table before every test to avoid dependencies and side effects between tests
	query := `
		DELETE FROM
			users`
	_, err := suite.pgContainer.DB.Exec(query)

	if err != nil {
		log.Fatal(err)
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
	require.Nil(suite.T(), err)

	user, err := suite.usersRepository.QueryGetUser("test@test.com")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.NotNil(suite.T(), user.ID)
	assert.Equal(suite.T(), "test@test.com", user.Email)
}
