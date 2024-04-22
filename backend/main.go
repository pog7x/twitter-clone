package main

import (
	"encoding/base64"
	"fmt"

	"twitter-clone/config"
	"twitter-clone/internal/handlers/http/api"
	"twitter-clone/internal/infrastructure/database"
	"twitter-clone/internal/infrastructure/logger"
	"twitter-clone/internal/infrastructure/middlewares"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12"
	irislog "github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"

	"github.com/samber/do"
)

func main() {
	cfg := config.LoadConfig()
	logger, _ := logger.NewLogger(cfg)

	hashKey, err := base64.StdEncoding.DecodeString(cfg.EncodedSessionHashKey)
	if err != nil {
		panic(err)
	}

	blockKey, err := base64.StdEncoding.DecodeString(cfg.EncodedSessionBlockKey)
	if err != nil {
		panic(err)
	}

	sc := securecookie.New(hashKey, blockKey)

	injector := do.NewWithOpts(&do.InjectorOpts{Logf: logger.Debugf})

	do.ProvideValue(injector, cfg)
	do.ProvideValue(injector, logger)

	do.Provide(injector, database.NewDatabase)

	do.Provide(injector, dbrepository.NewUserDBRepository)
	do.Provide(injector, dbrepository.NewTweetDBRepository)
	do.Provide(injector, dbrepository.NewMediaDBRepository)
	do.Provide(injector, dbrepository.NewLikeDBRepository)
	do.Provide(injector, dbrepository.NewSessionDBRepository)

	injector.HealthCheck()

	app := iris.New()

	app.Logger().SetLevel(cfg.LogLevel)

	app.UseRouter(requestid.New())
	app.UseRouter(recover.New())
	app.UseRouter(irislog.New())

	app.UseRouter(middlewares.CORS)

	app.PartyFunc("/login", func(login iris.Party) {
		login.Post("/", middlewares.LoginMiddleware(
			do.MustInvoke[dbrepository.UserRepository](injector),
			do.MustInvoke[dbrepository.SessionRepository](injector),
			sc,
		))
	})

	apiRouter := app.Party("/api")

	apiRouter.UseRouter(
		middlewares.SessionSecureCookieMiddleware(
			do.MustInvoke[dbrepository.SessionRepository](injector), sc,
		),
	)

	apiRouter.Party("/tweets").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[dbrepository.TweetRepository](injector))

		r.Post("/", api.CreateTweetHandler)
		r.Get("/", api.ListTweetHandler)

		r.Delete("/{id:uint64}/", api.DeleteTweetHandler)
	})

	apiRouter.Party("/tweets/{id:uint64}/likes/").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[dbrepository.LikeRepository](injector))

		r.Post("/", api.CreateLikeTweetHandler)
		r.Delete("/", api.DeleteLikeTweetHandler)
	})

	apiRouter.Party("/users").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[dbrepository.UserRepository](injector))

		// pass, _ := bcrypt.GenerateFromPassword([]byte("sosu_1"), bcrypt.DefaultCost)
		// do.MustInvoke[dbrepository.UserRepository](injector).Create(context.Background(), dbrepository.CreateUserPayload{
		// 	Name:     "huesos",
		// 	Password: pass,
		// 	Username: "hueta",
		// })

		r.Get("/me/", api.MeHandler)
		r.Get("/{id:uint64}/", api.UserHandler)

		r.Post("/{id:uint64}/follow/", api.FollowHandler)
		r.Delete("/{id:uint64}/follow/", api.UnfollowHandler)
	})

	apiRouter.Party("/medias").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[dbrepository.MediaRepository](injector))

		r.Post("/", api.CreateMediaHandler)
	})

	app.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), iris.WithOptimizations)

	injector.Shutdown()
}
