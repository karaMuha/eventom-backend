package services

import (
	"context"
	"eventom-backend/models"
	"eventom-backend/repositories"
	"eventom-backend/testutils"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

type UsersServiceTestSuite struct {
	suite.Suite
	ctx          context.Context
	pgContainer  *testutils.PostgresContainer
	usersService UsersServiceInterface
}

func TestUsersServiceSuite(t *testing.T) {
	suite.Run(t, &UsersServiceTestSuite{})
}

func (suite *UsersServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testutils.CreatePostgresContainer(suite.ctx)

	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer

	usersRepository := repositories.NewUsersRepository(suite.pgContainer.DB)
	suite.usersService = NewUsersService(usersRepository)
}

func (suite *UsersServiceTestSuite) BeforeTest(suiteName, testName string) {
	// clear users table before every test to avoid dependencies and side effects between tests
	query := `
		DELETE FROM
			users`
	_, err := suite.pgContainer.DB.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

func (suite *UsersServiceTestSuite) TearDownSuite() {
	err := suite.pgContainer.Terminate(suite.ctx)
	if err != nil {
		log.Fatalf("Error while terminating postgres container: %s", err)
	}
}

func (suite *UsersServiceTestSuite) TestSignupUserSuccess() {
	user := &models.User{
		Email:    "test@test.com",
		Password: "Test123",
	}

	err := suite.usersService.SignupUser(user)
	assert.Nil(suite.T(), err)

	// check if password was hashed
	createdUser, err := suite.usersService.GetUser(user.Email)
	assert.Nil(suite.T(), err)
	assert.NotEqual(suite.T(), user.Password, createdUser.Password)
	assert.NotNil(suite.T(), createdUser.ID)
}

func (suite *UsersServiceTestSuite) TestValidatePasswordUserNotFound() {
	user := &models.User{
		Email:    "test@test.com",
		Password: "Test123",
	}

	wrongEmail := &models.User{
		Email:    "wrongTest@test.com",
		Password: "Wrong123",
	}

	err := suite.usersService.SignupUser(user)
	assert.Nil(suite.T(), err)

	valid, err := suite.usersService.ValidateCredentials(wrongEmail)
	assert.Nil(suite.T(), err)
	assert.False(suite.T(), valid)
}

func (suite *UsersServiceTestSuite) TestValidatePasswordWrongPassword() {
	user := &models.User{
		Email:    "test@test.com",
		Password: "Test123",
	}

	wrongPw := &models.User{
		Email:    "test@test.com",
		Password: "Wrong123",
	}

	err := suite.usersService.SignupUser(user)
	assert.Nil(suite.T(), err)

	valid, err := suite.usersService.ValidateCredentials(wrongPw)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), 401, err.Status)
	assert.False(suite.T(), valid)
}

func (suite *UsersServiceTestSuite) TestValidatePasswordSuccess() {
	user := &models.User{
		Email:    "test@test.com",
		Password: "Test123",
	}

	err := suite.usersService.SignupUser(user)
	assert.Nil(suite.T(), err)

	valid, err := suite.usersService.ValidateCredentials(user)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), valid)
}
