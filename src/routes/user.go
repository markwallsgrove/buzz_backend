package routes

import (
	"context"
	"fmt"
	"net/http"

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
	results := models.Results{
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
