package database

import (
	"twitter-clone/config"

	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(i *do.Injector) (*Database, error) {
	cfg := do.MustInvoke[*config.Config](i)
	logger := do.MustInvoke[*logrus.Logger](i)

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logger.WithError(err).Errorf("Create database connection error")
		return nil, err
	}

	db.AutoMigrate(allModels...)

	return &Database{db}, nil
}

func (db *Database) HealthCheck() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

func (db *Database) Shutdown() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
