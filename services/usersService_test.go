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
	usersService UsersServiceInterface
}

func TestUsersServiceSuite(t *testing.T) {
	suite.Run(t, &UsersServiceTestSuite{})
}

func (suite *UsersServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	if testutils.TestContainer == nil {
		pgContainer, err := testutils.CreatePostgresContainer(suite.ctx)

		if err != nil {
			log.Fatal(err)
		}

		testutils.TestContainer = pgContainer
		testutils.TestContainer.Container.IsRunning()
	}

	usersRepository := repositories.NewUsersRepository(testutils.TestContainer.DB)
	suite.usersService = NewUsersService(usersRepository)
}

func (suite *UsersServiceTestSuite) BeforeTest(suiteName, testName string) {
	// clear users table before every test to avoid dependencies and side effects between tests
	query := `
		DELETE FROM
			users`
	_, err := testutils.TestContainer.DB.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

func (suite *UsersServiceTestSuite) TestSignupFailUniqueEmail() {
	user := &models.User{
		Email:    "test@test.com",
		Password: "Test123",
	}

	err := suite.usersService.SignupUser(user)
	assert.Nil(suite.T(), err)

	userTwo := &models.User{
		Email:    "test@test.com",
		Password: "Test123",
	}

	err = suite.usersService.SignupUser(userTwo)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), 409, err.Status)
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

func (suite *UsersServiceTestSuite) TestLoginUserNotFound() {
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

	token, err := suite.usersService.LoginUser(wrongEmail)
	assert.NotNil(suite.T(), err)
	assert.Empty(suite.T(), token)
}

func (suite *UsersServiceTestSuite) TestLogindWrongPassword() {
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

	token, err := suite.usersService.LoginUser(wrongPw)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), 401, err.Status)
	assert.Empty(suite.T(), token)
}
