package database

import (
	"database/sql"
	"log"
	"time"

	"bitbucket.org/brasilio/pandora/cerberus/config"
	_ "github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Pool struct {
	Write *gorm.DB
	Read  *gorm.DB
}

func NewPool(config *config.Config) (*Pool, error) {
	write, err := connect(config.WriteDSN, config)
	if err != nil {
		return nil, err
	}

	read, err := connect(config.ReadDSN, config)
	if err != nil {
		return nil, err
	}

	return &Pool{
		Write: write,
		Read:  read.Scopes(notDeleted, orderByCreated),
	}, nil
}

func connect(dsn string, config *config.Config) (*gorm.DB, error) {
	logger := logger.New(log.Default(), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	return gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{Logger: logger})
}

func notDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}

func orderByCreated(db *gorm.DB) *gorm.DB {
	return db.Order("created_at")
}
