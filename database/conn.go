package database

import (
	"errors"
	"marketplace-system/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrInvalidArgs = errors.New("nil arg or invalid argument")

func GetConnection(database config.Database) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(database.ConnectionString), &gorm.Config{})
	if err != nil {
		return
	}

	dbx, err := db.DB()
	dbx.SetConnMaxIdleTime(database.ConnMaxLifetime)
	dbx.SetConnMaxLifetime(database.ConnMaxLifetime)
	dbx.SetMaxOpenConns(database.MaxOpenConns)
	dbx.SetMaxIdleConns(database.MaxIdleConns)
	return
}
