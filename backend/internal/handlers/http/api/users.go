package api

import (
	"errors"
	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/infrastructure/database"
	"twitter-clone/internal/infrastructure/middlewares"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/kataras/iris/v12"
)

type User struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Followings []User `json:"following"`
	Followers  []User `json:"followers"`
}

type UserResp struct {
	Result bool `json:"result"`
	User   User `json:"user"`
}

func MeHandler(ctx iris.Context, userRepo dbrepository.UserRepository) {
	userID, err := ctx.Values().GetUint64(middlewares.UserIDKey)
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	user, err := userRepo.Get(ctx, dbrepository.GetUserPayload{UserID: userID})
	if err != nil {
		if errors.Is(err, dbrepository.ErrNotFound) {
			response.SendErrorResponse(ctx, iris.StatusNotFound, err.Error())
			return
		}
		response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(UserResp{Result: true, User: userDto(*user)})
}

func UserHandler(ctx iris.Context, userRepo dbrepository.UserRepository) {
	id, err := ctx.Params().GetUint64("id")
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	user, err := userRepo.Get(ctx, dbrepository.GetUserPayload{UserID: id})
	if err != nil {
		if errors.Is(err, dbrepository.ErrNotFound) {
			response.SendErrorResponse(ctx, iris.StatusNotFound, err.Error())
			return
		}
		response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(UserResp{Result: true, User: userDto(*user)})
}

func FollowHandler(ctx iris.Context, userRepo dbrepository.UserRepository) {
	id, err := ctx.Params().GetUint64("id")
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	userID, err := ctx.Values().GetUint64(middlewares.UserIDKey)
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	_, err = userRepo.Update(ctx, id, dbrepository.UpdateUserPayload{
		Followers: []*database.User{{ID: userID}},
	})
	if err != nil {
		if errors.Is(err, dbrepository.ErrNotFound) {
			response.SendErrorResponse(ctx, iris.StatusNotFound, err.Error())
			return
		}
		response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(iris.Map{"result": true})
}

func UnfollowHandler(ctx iris.Context, userRepo dbrepository.UserRepository) {
	id, err := ctx.Params().GetUint64("id")
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	userID, err := ctx.Values().GetUint64(middlewares.UserIDKey)
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	err = userRepo.Unfollow(ctx, userID, id)
	if err != nil {
		if errors.Is(err, dbrepository.ErrNotFound) {
			response.SendErrorResponse(ctx, iris.StatusNotFound, err.Error())
			return
		}
		response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(iris.Map{"result": true})
}

func userDto(user database.User) User {
	followings, followers := []User{}, []User{}
	for _, following := range user.Followings {
		followings = append(followings, userDto(*following))
	}

	for _, follower := range user.Followers {
		followers = append(followers, userDto(*follower))
	}

	return User{
		ID:         user.ID,
		Name:       user.Name,
		Followings: followings,
		Followers:  followers,
	}
}
