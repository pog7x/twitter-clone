package di

import (
	"twitter-clone/config"
	"twitter-clone/internal/infrastructure/authcookie"
	"twitter-clone/internal/infrastructure/database"
	"twitter-clone/internal/infrastructure/jwtprovide"
	"twitter-clone/internal/repository/dbrepository"
	"twitter-clone/internal/services/authservice"

	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

func NewInjector(logger *logrus.Logger, cfg *config.Config) (*do.Injector, func(*do.Injector, *logrus.Logger)) {
	injector := do.NewWithOpts(&do.InjectorOpts{Logf: logger.Debugf})

	do.ProvideValue(injector, cfg)
	do.ProvideValue(injector, logger)

	// infrastructure
	do.Provide(injector, database.NewDatabase)
	do.Provide(injector, jwtprovide.NewJWTProvider)
	do.Provide(injector, authcookie.NewCookieCodec)

	// repositories
	do.Provide(injector, dbrepository.NewUserDBRepository)
	do.Provide(injector, dbrepository.NewTweetDBRepository)
	do.Provide(injector, dbrepository.NewMediaDBRepository)
	do.Provide(injector, dbrepository.NewLikeDBRepository)
	do.Provide(injector, dbrepository.NewSessionDBRepository)
	do.Provide(injector, dbrepository.NewTrendDBRepository)

	// services
	do.Provide(injector, authservice.NewAuthService)

	injector.HealthCheck()

	return injector, func(i *do.Injector, l *logrus.Logger) {
		if err := i.Shutdown(); err != nil {
			l.WithError(err).Error("Dependencies injector shutdown error.")
		}
	}
}
