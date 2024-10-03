package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"eventom-backend/models"
	"eventom-backend/repositories"
	"eventom-backend/services"
	"eventom-backend/testutils"
	"eventom-backend/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UsersControllerTestSuite struct {
	suite.Suite
	ctx    context.Context
	router *http.ServeMux
}

func TestUsersControllerSuite(t *testing.T) {
	suite.Run(t, &UsersControllerTestSuite{})
}

func (suite *UsersControllerTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	if testutils.TestContainer == nil {
		pgContainer, err := testutils.CreatePostgresContainer(suite.ctx)

		if err != nil {
			log.Fatal(err)
		}

		testutils.TestContainer = pgContainer
	}

	logger := utils.NewLogger(os.Stdout)

	usersRepository := repositories.NewUsersRepository(testutils.TestContainer.DB)
	usersService := services.NewUsersService(usersRepository)
	usersController := NewUsersController(usersService, logger)

	router := http.NewServeMux()
	router.HandleFunc("POST /signup", usersController.HandleSignupUser)
	router.HandleFunc("POST /login", usersController.HandleLoginUser)

	suite.router = router
}

func (suite *UsersControllerTestSuite) BeforeTest(suiteName, testName string) {
	// clear users table before every test to avoid dependencies and side effects between tests
	query := `
		DELETE FROM
			users`
	_, err := testutils.TestContainer.DB.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

func (suite *UsersControllerTestSuite) TestSignupFailNoEmail() {
	user := &models.User{
		Email:    "",
		Password: "Test123",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/signup", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	assert.Equal(suite.T(), 400, recorder.Result().StatusCode)
}

func (suite *UsersControllerTestSuite) TestSignupFailNoPassword() {
	user := &models.User{
		Email:    "test@test.com",
		Password: "",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/signup", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	assert.Equal(suite.T(), 400, recorder.Result().StatusCode)
}

func (suite *UsersControllerTestSuite) TestSignupSuccess() {
	user := &models.User{
		Email:    "test@test.com",
		Password: "Test123",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/signup", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	assert.Equal(suite.T(), 200, recorder.Result().StatusCode)
}

func (suite *UsersControllerTestSuite) TestSignupFailUniqueEmail() {
	user := &models.User{
		Email:    "test@test.com",
		Password: "Test123",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/signup", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	// make a second request with same user, this request should fail due to unique email constraint
	userTwo := &models.User{
		Email:    "test@test.com",
		Password: "Test123",
	}

	userTwoBytes, err := json.Marshal(userTwo)
	if err != nil {
		log.Fatal(err)
	}

	bodyTwo := bytes.NewReader(userTwoBytes)
	requestFail := httptest.NewRequest("POST", "/signup", bodyTwo)
	recorderFail := httptest.NewRecorder()

	suite.router.ServeHTTP(recorderFail, requestFail)

	assert.Equal(suite.T(), 409, recorderFail.Result().StatusCode)
}

func (suite *UsersControllerTestSuite) TestLoginFailNoUser() {
	user := &models.User{
		Email:    "test@test.com",
		Password: "Test123",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/login", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	assert.Equal(suite.T(), 401, recorder.Result().StatusCode)
}

func (suite *UsersControllerTestSuite) TestLoginFailWrongPassword() {
	user := &models.User{
		Email:    "test@test.com",
		Password: "Test123",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	signupRequest := httptest.NewRequest("POST", "/signup", body)
	signupRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(signupRecorder, signupRequest)

	userWrongPw := &models.User{
		Email:    "test@test.com",
		Password: "WrongPw",
	}

	userWrongPwBytes, err := json.Marshal(userWrongPw)
	if err != nil {
		log.Fatal(err)
	}

	bodyWrongPw := bytes.NewReader(userWrongPwBytes)
	loginRequest := httptest.NewRequest("POST", "/login", bodyWrongPw)
	loginRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(loginRecorder, loginRequest)

	assert.Equal(suite.T(), 401, loginRecorder.Result().StatusCode)
}
