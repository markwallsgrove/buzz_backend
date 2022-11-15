package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	randomdata "github.com/Pallinder/go-randomdata"
	echo "github.com/labstack/echo/v4"
	"github.com/markwallsgrove/muzz_devops/src/database"
	"github.com/markwallsgrove/muzz_devops/src/models"
	"go.uber.org/zap"
)

// NewUserController helper function to create a new user controller
func NewUserController(ctx context.Context, db database.Database, logger *zap.Logger) *UserController {
	return &UserController{
		db,
		logger,
		ctx,
	}
}

// UserController grouping of routes to handle user requests
type UserController struct {
	database.Database
	Logger *zap.Logger
	Ctx    context.Context
}

// CreateUser create a new random user & return the model
// as a JSON blob within the body.
func (u *UserController) CreateUser(c echo.Context) error {
	password, err := models.CreateRandomPassword()
	if err != nil {
		u.Logger.Error("cannot create password", zap.Error(err))
		return err
	}

	hash, err := models.CreatePasswordHash(password)
	if err != nil {
		u.Logger.Error("cannot create hash", zap.Error(err))
		return err
	}

	var gender models.Gender
	var firstName string
	if randomdata.Boolean() {
		firstName = randomdata.FirstName(randomdata.Male)
		gender = models.Male
	} else {
		firstName = randomdata.FirstName(randomdata.Female)
		gender = models.Female
	}

	lastName := randomdata.LastName()
	age := randomdata.Number(13, 100)

	user := models.User{
		Email:    fmt.Sprintf("%s.%s@gmail.com", firstName, lastName),
		Password: hash,
		Name:     fmt.Sprintf("%s %s", firstName, lastName),
		Gender:   gender,
		Age:      age,
	}

	user, err = u.Database.CreateUser(u.Ctx, user)
	if err != nil {
		return err
	}

	// return the password back to the user rather than the hash
	user.Password = password
	results := models.Result{
		Result: &user,
	}

	err = c.JSON(http.StatusOK, results)
	if err == nil {
		return nil
	}

	u.Logger.Error("cannot marshall new user into response body", zap.Error(err))
	if err := c.String(http.StatusInternalServerError, "internal error"); err != nil {
		u.Logger.Error("cannot set create user status", zap.Error(err))
	}

	return err
}

// Profiles fetch profiles that are matched to the current user.
//
// Query Params:
//   - userId (must be numeric)
//
// Status codes:
//   - 404 user was not found
//   - 500 internal server error
//   - 200 results
//
// Matching two users is based on their age difference and the opposite sex.
// All matches will be within five years of the current user's age, unless
// the value is below 13 or above 100.
func (u *UserController) Profiles(c echo.Context) error {
	userId := c.QueryParam("userId")
	user, err := u.Database.GetUser(u.Ctx, userId)

	if err == database.ErrNotFound {
		return c.String(http.StatusNotFound, "user not found")
	}
	if err != nil {
		u.Logger.Error("cannot find user for profiles", zap.Error(err))
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	gender := models.Female
	if user.Gender == models.Female {
		gender = models.Male
	}

	minAge := user.Age - 5
	if minAge < 13 {
		minAge = 13
	}

	maxAge := user.Age + 5
	if maxAge > 100 {
		maxAge = 100
	}

	users, err := u.Database.FindMatches(u.Ctx, models.Gender(gender), minAge, maxAge)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, models.Results{Results: users})
}

// Create a swipe between two users. The swipe will contain who swiped.
//
// Query Params:
//   - currentUser the user who is swiping
//   - targetUser who the user has swiped
//
// Status codes:
//   - 500 internal server error
//   - 200 success
func (u *UserController) Swipe(c echo.Context) error {
	currentUserId, err := strconv.Atoi(c.QueryParam("currentUser"))
	if err != nil {
		return c.String(http.StatusBadRequest, "current user is not numeric")
	}

	targetUserId, err := strconv.Atoi(c.QueryParam("targetUser"))
	if err != nil {
		return c.String(http.StatusBadRequest, "target user is not numeric")
	}

	swipe, err := u.Database.Swipe(u.Ctx, currentUserId, targetUserId)
	if err != nil {
		u.Logger.Error("cannot create swipe", zap.Error(err))
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	swipeResult := &models.SwipeResult{
		Matched: false,
	}

	if swipe.FirstUserSwiped && swipe.SecondUserSwiped {
		swipeResult.ID = swipe.ID
		swipeResult.Matched = true
	}

	return c.JSON(http.StatusOK, models.SwipeResults{
		Results: *swipeResult,
	})
}
