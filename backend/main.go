package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"golang.org/x/crypto/bcrypt"

	"twitter-clone/config"
	"twitter-clone/internal/handlers/http/api"
	"twitter-clone/internal/infrastructure/di"
	"twitter-clone/internal/infrastructure/jwtprovide"
	"twitter-clone/internal/infrastructure/logger"
	"twitter-clone/internal/infrastructure/middlewares"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/kataras/iris/v12"
	irisLog "github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"

	"github.com/samber/do"
)

func main() {
	cfg := config.LoadConfig()
	log, err := logger.NewLogger(cfg)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(cfg.UploadsDirPath); os.IsNotExist(err) {
		if err = os.Mkdir(cfg.UploadsDirPath, fs.ModePerm); err != nil {
			panic(err)
		}
	}

	dInj := di.NewInjector(log, cfg)
	defer dInj.Shutdown()

	app := iris.New()

	app.HandleDir(cfg.UploadsDirPath, iris.Dir(cfg.UploadsDirPath))

	app.Logger().SetLevel(cfg.LogLevel)

	app.UseRouter(requestid.New())
	app.UseRouter(recover.New())
	app.UseRouter(irisLog.New())
	app.UseRouter(middlewares.CORS)

	app.PartyFunc("/login", func(login iris.Party) {
		login.Post("/", middlewares.SessionJWTLoginMiddleware(
			do.MustInvoke[*dbrepository.UserRepository](dInj),
			do.MustInvoke[*dbrepository.SessionRepository](dInj),
			do.MustInvoke[*jwtprovide.JWTProvider](dInj).Signer,
		))
	})

	apiRouter := app.Party("/api")

	apiRouter.UseRouter(
		middlewares.SessionJWTMiddleware(
			do.MustInvoke[*dbrepository.SessionRepository](dInj),
			do.MustInvoke[*jwtprovide.JWTProvider](dInj).Verifier,
		),
	)

	apiRouter.Party("/tweets").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[*dbrepository.TweetRepository](dInj))

		r.Post("/", api.CreateTweetHandler)
		r.Get("/", api.ListTweetHandler)

		r.Delete("/{id:uint64}/", api.DeleteTweetHandler)
		r.Patch("/{id:uint64}/", api.UpdateTweetHandler)

		r.Get("/user/{id:uint64}/", api.UserTweetsHandler)
	})

	apiRouter.Party("/tweets/{id:uint64}/likes/").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[*dbrepository.LikeRepository](dInj))

		r.Post("/", api.CreateLikeTweetHandler)
		r.Delete("/", api.DeleteLikeTweetHandler)
	})

	apiRouter.Party("/users").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[*dbrepository.UserRepository](dInj))

		pass, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
		do.MustInvoke[*dbrepository.UserRepository](dInj).Create(context.Background(), dbrepository.CreateUserPayload{
			Password:    pass,
			Username:    "pog7x",
			Name:        "developer",
			Website:     "https://github.com/pog7x",
			PicCover:    "https://ideogram.ai/api/images/direct/T91kUQhETeyPyiyqCOfwcQ.png",
			Pic:         "https://avataaars.io/?avatarStyle=Circle&topType=LongHairFrida&accessoriesType=Round&facialHairType=Blank&clotheType=ShirtVNeck&clotheColor=Gray01&eyeType=Happy&eyebrowType=RaisedExcitedNatural&mouthType=Smile&skinColor=Pale",
			Description: "Just a developer that interested in JavaScript.",
		})

		r.Get("/me/", api.MeHandler)
		r.Put("/me/", api.MeUpdateHandler)
		r.Get("/{id:uint64}/", api.UserHandler)

		r.Post("/{id:uint64}/follow/", api.FollowHandler)
		r.Delete("/{id:uint64}/follow/", api.UnfollowHandler)
	})

	apiRouter.Party("/medias").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[*dbrepository.MediaRepository](dInj))

		r.Post("/", api.CreateMediaHandler)
	})

	apiRouter.Party("/trends").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[*dbrepository.TrendRepository](dInj))

		r.Get("/", api.ListTrendHandler)
	})

	app.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), iris.WithOptimizations)
}
