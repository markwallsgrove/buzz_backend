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
	"github.com/markwallsgrove/muzz_devops/src/models/domain"
	"github.com/markwallsgrove/muzz_devops/src/models/httpDomain"
	"github.com/markwallsgrove/muzz_devops/src/models/security"
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
	password, err := security.CreateRandomPassword()
	if err != nil {
		u.Logger.Error("cannot create password", zap.Error(err))
		return err
	}

	hash, err := security.CreatePasswordHash(password)
	if err != nil {
		u.Logger.Error("cannot create hash", zap.Error(err))
		return err
	}

	var gender domain.Gender
	var firstName string
	if randomdata.Boolean() {
		firstName = randomdata.FirstName(randomdata.Male)
		gender = domain.Male
	} else {
		firstName = randomdata.FirstName(randomdata.Female)
		gender = domain.Female
	}

	lastName := randomdata.LastName()
	age := randomdata.Number(13, 100)

	user := &domain.User{
		Email:        fmt.Sprintf("%s.%s@gmail.com", firstName, lastName),
		PasswordHash: hash,
		Name:         fmt.Sprintf("%s %s", firstName, lastName),
		Gender:       gender,
		Age:          age,
	}

	err = u.Database.CreateUser(u.Ctx, user)
	if err != nil {
		return err
	}

	// return the password back to the user rather than the hash
	results := httpDomain.UserResult{
		Result: httpDomain.UserToHTTPDomain(user, password),
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
// Status codes:
//   - 404 user was not found
//   - 500 internal server error
//   - 200 results
//
// Matching two users is based on their age difference and the opposite sex.
// All matches will be within five years of the current user's age, unless
// the value is below 13 or above 100.
func (u *UserController) Profiles(c echo.Context) error {
	uid, ok := c.Get("userId").(int)
	if !ok {
		u.Logger.Error("cannot find user id within context")
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	var genders []domain.Gender
	gender := domain.StringToGender(c.QueryParam("gender"))
	if gender == domain.UnknownGender {
		genders = domain.Genders
	} else {
		genders = []domain.Gender{gender}
	}

	minAge := models.StringToInt(c.QueryParam("minAge"), 0)
	maxAge := models.StringToInt(c.QueryParam("maxAge"), 200)

	userProfiles, err := u.Database.FindMatches(u.Ctx, uid, genders, minAge, maxAge)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, httpDomain.UserProfileResults{Results: userProfiles})
}

// Create a swipe between two users. The swipe will contain who swiped.
//
// Query Params:
//   - targetUser who the user has swiped
//
// Status codes:
//   - 500 internal server error
//   - 200 success
func (u *UserController) Swipe(c echo.Context) error {
	currentUserId, ok := c.Get("userId").(int)
	if !ok {
		u.Logger.Error("cannot find user id within context")
		return c.String(http.StatusBadRequest, "internal server error")
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

	swipeResult := &httpDomain.SwipeResult{
		Matched: false,
	}

	if swipe.FirstUserSwiped && swipe.SecondUserSwiped {
		swipeResult.ID = swipe.ID
		swipeResult.Matched = true
	}

	return c.JSON(http.StatusOK, httpDomain.SwipeResults{
		Results: *swipeResult,
	})
}
