package routes_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/markwallsgrove/muzz_devops/src/models"
	"github.com/markwallsgrove/muzz_devops/src/routes"
	"github.com/markwallsgrove/muzz_devops/src/routes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

//go:generate mockery --dir ../database --filename mocks.go --name Database

func TestUsersRoute(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user/create", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db := &mocks.Database{}
	db.On(
		"CreateUser",
		mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("models.User"),
	).Return(nil)

	users := routes.UserController{
		db,
		zap.NewNop(),
		context.Background(),
	}

	assert.NoError(t, users.CreateUser(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check the user data that was sent to the database
	assert.True(t, db.AssertNumberOfCalls(t, "CreateUser", 1))
	assert.Len(t, db.Calls, 1)
	call := db.Calls[0]
	assert.Equal(t, "CreateUser", call.Method)

	user := call.Arguments.Get(1).(models.User)
	assert.GreaterOrEqual(t, user.Age, 13)
	assert.LessOrEqual(t, user.Age, 100)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Name)
	assert.NotEmpty(t, user.Password)
	assert.GreaterOrEqual(t, user.Gender, models.UnknownGender)

	// Check the user data that was encoded in the HTTP response body
	var results models.Results
	assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &results))

	assert.Equal(t, user.Age, results.Result.Age)
	assert.Equal(t, user.Email, results.Result.Email)
	assert.Equal(t, user.Name, results.Result.Name)
	assert.Equal(t, user.Gender, results.Result.Gender)

	// should not return the hash which was sent to the database
	assert.NotEqual(t, user.Password, results.Result.Password)
	assert.NotEmpty(t, results.Result.Password)
}
