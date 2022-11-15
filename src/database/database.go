package database

import (
	"context"
	"errors"
	"strconv"

	"github.com/markwallsgrove/muzz_devops/src/models"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	GetSwipe(
		ctx context.Context,
		firstUserId int,
		secondUserId int,
	) (models.Swipe, error)
	Swipe(
		ctx context.Context,
		firstUserId int,
		secondUserId int,
	) (models.Swipe, error)
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

// Swipe create a new swipe which details which user (if not both) swipped
func (d *MariaDB) Swipe(
	ctx context.Context,
	firstUserId int,
	secondUserId int,
) (models.Swipe, error) {
	// Ensure the records have a predictable index by ordering the user ids.
	// If we do not order the IDs between the columns we might get duplicates.
	var updatesOnConflict map[string]interface{}
	var swipe models.Swipe

	if firstUserId < secondUserId {
		updatesOnConflict = map[string]interface{}{
			"first_user_swiped": true,
		}
		swipe = models.Swipe{
			FirstUserID:      firstUserId,
			SecondUserID:     secondUserId,
			FirstUserSwiped:  true,
			SecondUserSwiped: false,
		}
	} else {
		updatesOnConflict = map[string]interface{}{
			"second_user_swiped": true,
		}
		swipe = models.Swipe{
			FirstUserID:      secondUserId,
			SecondUserID:     firstUserId,
			FirstUserSwiped:  false,
			SecondUserSwiped: true,
		}
	}

	// If the update fails because a swipe already exists the record will be
	// updated to show the current user swiped.
	//
	// If a record does not already exist it will be created.
	result := d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(updatesOnConflict),
	}).Create(&swipe)

	if result.Error != nil {
		d.Logger.Error("cannot update/create swipe", zap.Error(result.Error))
		return models.Swipe{}, result.Error
	}

	return d.GetSwipe(ctx, firstUserId, secondUserId)
}

// GetSwipe retrieve a swipe by the two user ids
func (d *MariaDB) GetSwipe(
	ctx context.Context,
	firstUserId int,
	secondUserId int,
) (models.Swipe, error) {
	// We must use the ids in the correct order else we might not
	// be able to find the record.
	fuid := firstUserId
	suid := secondUserId

	if fuid > suid {
		t := suid
		suid = fuid
		fuid = t
	}

	var swipe models.Swipe
	result := d.db.
		Where("first_user_id = ? AND second_user_id = ?", fuid, suid).
		Find(&swipe)

	if result.Error != nil && result.RowsAffected == 0 {
		d.Logger.Error("cannot retrieve swipe", zap.Error(result.Error))
		return models.Swipe{}, result.Error
	}

	return swipe, nil
}

func (d *MariaDB) Close() error {
	return nil
}
