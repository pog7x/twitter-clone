package main

import (
	"fmt"

	"twitter-clone/config"
	"twitter-clone/internal/handlers/http/api"
	"twitter-clone/internal/infrastructure/di"
	"twitter-clone/internal/infrastructure/jwtprovide"
	"twitter-clone/internal/infrastructure/logger"
	"twitter-clone/internal/infrastructure/middlewares"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/kataras/iris/v12"
	irislog "github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"

	"github.com/samber/do"
)

func main() {
	cfg := config.LoadConfig()
	logger, err := logger.NewLogger(cfg)
	if err != nil {
		panic(err)
	}

	dinj := di.NewInjector(logger, cfg)
	defer dinj.Shutdown()

	app := iris.New()

	app.HandleDir("/uploads", iris.Dir("./uploads"))

	app.Logger().SetLevel(cfg.LogLevel)

	app.UseRouter(requestid.New())
	app.UseRouter(recover.New())
	app.UseRouter(irislog.New())
	app.UseRouter(middlewares.CORS)

	app.PartyFunc("/login", func(login iris.Party) {
		login.Post("/", middlewares.SessionJWTLoginMiddleware(
			do.MustInvoke[*dbrepository.UserRepository](dinj),
			do.MustInvoke[*dbrepository.SessionRepository](dinj),
			do.MustInvoke[*jwtprovide.JWTProvider](dinj).Signer,
		))
	})

	apiRouter := app.Party("/api")

	apiRouter.UseRouter(
		middlewares.SessionJWTMiddleware(
			do.MustInvoke[*dbrepository.SessionRepository](dinj),
			do.MustInvoke[*jwtprovide.JWTProvider](dinj).Verifier,
		),
	)

	apiRouter.Party("/tweets").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[*dbrepository.TweetRepository](dinj))

		r.Post("/", api.CreateTweetHandler)
		r.Get("/", api.ListTweetHandler)

		r.Delete("/{id:uint64}/", api.DeleteTweetHandler)
		r.Patch("/{id:uint64}/", api.UpdateTweetHandler)
	})

	apiRouter.Party("/tweets/{id:uint64}/likes/").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[*dbrepository.LikeRepository](dinj))

		r.Post("/", api.CreateLikeTweetHandler)
		r.Delete("/", api.DeleteLikeTweetHandler)
	})

	apiRouter.Party("/users").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[*dbrepository.UserRepository](dinj))

		// pass, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
		// do.MustInvoke[*dbrepository.UserRepository](dinj).Create(context.Background(), dbrepository.CreateUserPayload{
		// 	Password:    pass,
		// 	Username:    "pog7x",
		// 	Name:        "developer",
		// 	Website:     "https://github.com/pog7x",
		// 	PicCover:    "https://ideogram.ai/api/images/direct/T91kUQhETeyPyiyqCOfwcQ.png",
		// 	Pic:         "https://avataaars.io/?avatarStyle=Circle&topType=LongHairFrida&accessoriesType=Round&facialHairType=Blank&clotheType=ShirtVNeck&clotheColor=Gray01&eyeType=Happy&eyebrowType=RaisedExcitedNatural&mouthType=Smile&skinColor=Pale",
		// 	Description: "Just a developer that interested in JavaScript.",
		// })

		r.Get("/me/", api.MeHandler)
		r.Put("/me/", api.MeUpdateHandler)
		r.Get("/{id:uint64}/", api.UserHandler)

		r.Post("/{id:uint64}/follow/", api.FollowHandler)
		r.Delete("/{id:uint64}/follow/", api.UnfollowHandler)
	})

	apiRouter.Party("/medias").ConfigureContainer(func(r *iris.APIContainer) {
		r.RegisterDependency(do.MustInvoke[*dbrepository.MediaRepository](dinj))

		r.Post("/", api.CreateMediaHandler)
	})

	app.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), iris.WithOptimizations)
}
