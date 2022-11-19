package routes_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/markwallsgrove/muzz_devops/src/models/domain"
	"github.com/markwallsgrove/muzz_devops/src/models/httpDomain"
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
		mock.AnythingOfType("*domain.User"),
	).Return(
		nil,
	)

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

	user := call.Arguments.Get(1).(*domain.User)
	assert.GreaterOrEqual(t, user.Age, 13)
	assert.LessOrEqual(t, user.Age, 100)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Name)
	assert.NotEmpty(t, user.PasswordHash)
	assert.GreaterOrEqual(t, user.Gender, domain.UnknownGender)

	// Check the user data that was encoded in the HTTP response body
	var results httpDomain.UserResult
	assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &results))

	result := results.Result
	assert.Equal(t, user.Age, result.Age)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, "Female", result.Gender)

	// should not return the hash which was sent to the database
	assert.NotEmpty(t, result.Password)
	assert.NotEmpty(t, result.Password)
}

func TestUsersProfiles(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/profiles?userId=9", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db := &mocks.Database{}

	user := domain.User{
		ID:     9,
		Email:  "foo.bar@example.com",
		Name:   "foo bar",
		Gender: domain.Male,
		Age:    20,
	}

	profiles := []domain.UserProfile{
		{
			ID:     1,
			Name:   "Ann Thompson",
			Gender: domain.Female,
			Age:    17,
		},
	}

	db.On(
		"GetUser",
		mock.AnythingOfType("*context.emptyCtx"),
		9,
	).Return(user, nil)

	db.On(
		"FindMatches",
		mock.AnythingOfType("*context.emptyCtx"),
		9,
		domain.Gender(domain.Female),
		15,
		25,
	).Return(profiles, nil)

	users := routes.UserController{
		db,
		zap.NewNop(),
		context.Background(),
	}

	assert.NoError(t, users.Profiles(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// check how the GetUser call to the database was called
	call := db.Calls[0]
	assert.Equal(t, "GetUser", call.Method)
	assert.Equal(t, 9, call.Arguments.Get(1))

	// check how the FindMatches call to the database was called
	call = db.Calls[1]
	assert.Equal(t, "FindMatches", call.Method)
	assert.Equal(t, 9, call.Arguments.Get(1))
	assert.Equal(t, domain.Gender(domain.Female), call.Arguments.Get(2))
	assert.Equal(t, 15, call.Arguments.Get(3))
	assert.Equal(t, 25, call.Arguments.Get(4))
}

func TestSwipe(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/swipe?currentUser=9&targetUser=10", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db := &mocks.Database{}

	db.On("Swipe", mock.AnythingOfType("*context.emptyCtx"), 9, 10).Return(
		domain.Swipe{
			ID:               11,
			FirstUserID:      9,
			SecondUserID:     10,
			FirstUserSwiped:  true,
			SecondUserSwiped: false,
		},
		nil,
	)

	users := routes.UserController{
		db,
		zap.NewNop(),
		context.Background(),
	}

	assert.NoError(t, users.Swipe(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var result httpDomain.SwipeResults
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
	swipe := result.Results

	assert.False(t, swipe.Matched)
}

func TestSwipeMatched(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/swipe?currentUser=9&targetUser=10", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db := &mocks.Database{}

	db.On("Swipe", mock.AnythingOfType("*context.emptyCtx"), 9, 10).Return(
		domain.Swipe{
			ID:               11,
			FirstUserID:      9,
			SecondUserID:     10,
			FirstUserSwiped:  true,
			SecondUserSwiped: true,
		},
		nil,
	)

	users := routes.UserController{
		db,
		zap.NewNop(),
		context.Background(),
	}

	assert.NoError(t, users.Swipe(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var result httpDomain.SwipeResults
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
	swipe := result.Results

	assert.Equal(t, 11, swipe.ID)
	assert.True(t, swipe.Matched)
}
