package database

import (
	"context"
	"errors"
	"strconv"

	"github.com/markwallsgrove/muzz_devops/src/models"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("Cannot find resource")

type Database interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, id string) (models.User, error)
	FindMatches(
		ctx context.Context,
		gender models.Gender,
		minAge int,
		maxAge int,
	) ([]models.UserProfile, error)
	Close() error
}

// NewMariaDB create a new MariaDB instance
func NewMariaDB(dsn string, logger *zap.Logger) (Database, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return &MariaDB{}, err
	}

	return &MariaDB{db, logger}, nil
}

type MariaDB struct {
	db     *gorm.DB
	Logger *zap.Logger
}

// CreateUser create a new user. The ID will be ignored as it's automatically generated &
// the email address must be unique.
func (d *MariaDB) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("id").Create(&user).Error; err != nil {
			d.Logger.Error("cannot create user", zap.Error(err))
			return err
		}

		return nil
	})

	return user, err
}

// GetUser find a user by their id
func (d *MariaDB) GetUser(ctx context.Context, id string) (models.User, error) {
	var user models.User

	// check if user id is numeric
	if _, err := strconv.Atoi(id); err != nil {
		d.Logger.Error("user id is not numeric", zap.Error(err))
		return models.User{}, err
	}

	results := d.db.Find(&user, id)
	if results.Error != nil {
		d.Logger.Error("cannot find user", zap.Error(results.Error))
		return models.User{}, results.Error
	}

	if results.RowsAffected != 1 {
		return models.User{}, ErrNotFound
	}

	return user, nil
}

// FindMatches find potenial matches
func (d *MariaDB) FindMatches(
	ctx context.Context,
	gender models.Gender,
	minAge int,
	maxAge int,
) ([]models.UserProfile, error) {
	var profiles []models.UserProfile

	results := d.db.
		Model(&models.User{}).
		Where("age >= ? AND age <= ? AND gender = ?", minAge, maxAge, gender).
		Find(&profiles)

	if results.Error != nil {
		return []models.UserProfile{}, results.Error
	}

	return profiles, nil
}

func (d *MariaDB) Close() error {
	return nil
}
