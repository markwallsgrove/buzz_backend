package database

import (
	"context"

	"github.com/markwallsgrove/muzz_devops/src/models"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	Close() error
}

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

func (d *MariaDB) Close() error {
	return nil
}
