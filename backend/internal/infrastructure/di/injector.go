package di

import (
	"twitter-clone/config"
	"twitter-clone/internal/infrastructure/authcookie"
	"twitter-clone/internal/infrastructure/database"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

func NewInjector(logger *logrus.Logger, cfg *config.Config) *do.Injector {
	injector := do.NewWithOpts(&do.InjectorOpts{Logf: logger.Debugf})

	do.ProvideValue(injector, cfg)
	do.ProvideValue(injector, logger)

	do.Provide(injector, database.NewDatabase)
	do.Provide(injector, authcookie.NewCookieCodec)

	do.Provide(injector, dbrepository.NewUserDBRepository)
	do.Provide(injector, dbrepository.NewTweetDBRepository)
	do.Provide(injector, dbrepository.NewMediaDBRepository)
	do.Provide(injector, dbrepository.NewLikeDBRepository)
	do.Provide(injector, dbrepository.NewSessionDBRepository)

	injector.HealthCheck()

	return injector
}
