package database

import (
	"twitter-clone/config"

	"github.com/samber/do"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(i *do.Injector) (Database, error) {
	cfg := do.MustInvoke[*config.Config](i)

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database") // TODO
	}

	db.AutoMigrate(allModels...)

	return Database{db}, nil
}

func (db Database) HealthCheck() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

func (db Database) Shutdown() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
