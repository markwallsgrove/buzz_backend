package database

import (
	"context"
	"errors"

	"github.com/markwallsgrove/muzz_devops/src/models/domain"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNotFound = errors.New("Cannot find resource")

type Database interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id int) (domain.User, error)
	FindMatches(
		ctx context.Context,
		user *domain.User,
		gender []domain.Gender,
		minAge int,
		maxAge int,
	) ([]domain.UserProfile, error)
	GetSwipe(
		ctx context.Context,
		firstUserId int,
		secondUserId int,
	) (domain.Swipe, error)
	Swipe(
		ctx context.Context,
		firstUserId int,
		secondUserId int,
	) (domain.Swipe, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (*domain.User, error)
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
func (d *MariaDB) CreateUser(ctx context.Context, user *domain.User) error {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("id").Create(&user).Error; err != nil {
			d.Logger.Error("cannot create user", zap.Error(err))
			return err
		}

		return nil
	})

	return err
}

// GetUser find a user by their id
func (d *MariaDB) GetUser(ctx context.Context, id int) (domain.User, error) {
	var user domain.User

	results := d.db.Find(&user, "ID = ?", id)
	if results.Error != nil {
		d.Logger.Error("cannot find user", zap.Error(results.Error))
		return domain.User{}, results.Error
	}

	if results.RowsAffected != 1 {
		return domain.User{}, ErrNotFound
	}

	return user, nil
}

// FindMatches find potenial matches
func (d *MariaDB) FindMatches(
	ctx context.Context,
	user *domain.User,
	genders []domain.Gender,
	minAge int,
	maxAge int,
) ([]domain.UserProfile, error) {
	var users []domain.UserDistance

	results := d.db.Raw(
		"SELECT *, ST_DISTANCE(u.location, POINT(?, ?)) as distance FROM dating.users u WHERE u.id != ? AND u.id NOT IN (? UNION ?) AND u.gender IN ? AND u.age >= ? AND u.age <= ? ORDER by distance ASC",
		user.Location.X,
		user.Location.Y,
		user.ID,
		d.db.
			Table("dating.swipes s1").
			Select("s1.second_user_id as id").
			Where("s1.first_user_id = ? AND s1.first_user_swiped = TRUE", user.ID),
		d.db.
			Table("dating.swipes s2").
			Select("s2.first_user_id as id").
			Where("s2.second_user_id = ? AND s2.second_user_swiped = TRUE", user.ID),
		genders,
		minAge,
		maxAge,
	).Scan(&users)

	if results.Error != nil {
		return []domain.UserProfile{}, results.Error
	}

	profiles := make([]domain.UserProfile, len(users))
	for i, user := range users {
		profiles[i] = domain.UserProfile{
			ID:             user.ID,
			Name:           user.Name,
			Gender:         user.Gender.String(),
			Age:            user.Age,
			DistanceFromMe: user.Distance,
		}
	}

	return profiles, nil
}

// Swipe create a new swipe which details which user (if not both) swipped
func (d *MariaDB) Swipe(
	ctx context.Context,
	firstUserId int,
	secondUserId int,
) (domain.Swipe, error) {
	// Ensure the records have a predictable index by ordering the user ids.
	// If we do not order the IDs between the columns we might get duplicates.
	var updatesOnConflict map[string]interface{}
	var swipe domain.Swipe

	if firstUserId < secondUserId {
		updatesOnConflict = map[string]interface{}{
			"first_user_swiped": true,
		}
		swipe = domain.Swipe{
			FirstUserID:      firstUserId,
			SecondUserID:     secondUserId,
			FirstUserSwiped:  true,
			SecondUserSwiped: false,
		}
	} else {
		updatesOnConflict = map[string]interface{}{
			"second_user_swiped": true,
		}
		swipe = domain.Swipe{
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
		return domain.Swipe{}, result.Error
	}

	return d.GetSwipe(ctx, firstUserId, secondUserId)
}

// GetSwipe retrieve a swipe by the two user ids
func (d *MariaDB) GetSwipe(
	ctx context.Context,
	firstUserId int,
	secondUserId int,
) (domain.Swipe, error) {
	// We must use the ids in the correct order else we might not
	// be able to find the record.
	fuid := firstUserId
	suid := secondUserId

	if fuid > suid {
		t := suid
		suid = fuid
		fuid = t
	}

	var swipe domain.Swipe
	result := d.db.
		Where("first_user_id = ? AND second_user_id = ?", fuid, suid).
		Find(&swipe)

	if result.Error != nil && result.RowsAffected == 0 {
		d.Logger.Error("cannot retrieve swipe", zap.Error(result.Error))
		return domain.Swipe{}, result.Error
	}

	return swipe, nil
}

// GetUserByEmail Load a user's details by their email address
func (d *MariaDB) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user *domain.User
	if err := d.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (d *MariaDB) Close() error {
	return nil
}
