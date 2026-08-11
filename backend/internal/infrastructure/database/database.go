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

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL))
	if err != nil {
		logger.WithError(err).Errorf("Create database connection error")
		return nil, err
	}

	if err = dedupeLikes(db); err != nil {
		logger.WithError(err).Errorf("Deduplicating likes before migration error")
		return nil, err
	}

	if err = db.AutoMigrate(allModels...); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

// dedupeLikes clears rows that would prevent the unique (user_id, tweet_id)
// index from being created on a database written by an older build: previously
// likes were soft-deleted and could be inserted twice.
func dedupeLikes(db *gorm.DB) error {
	migrator := db.Migrator()

	if !migrator.HasTable(&Like{}) {
		return nil
	}

	if migrator.HasColumn(&Like{}, "deleted_at") {
		if err := db.Exec(`DELETE FROM likes WHERE deleted_at IS NOT NULL`).Error; err != nil {
			return err
		}
	}

	return db.Exec(
		`DELETE FROM likes a USING likes b
		 WHERE a.id > b.id AND a.user_id = b.user_id AND a.tweet_id = b.tweet_id`,
	).Error
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
