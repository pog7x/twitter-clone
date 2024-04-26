package api

import (
	"errors"
	"time"
	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/infrastructure/database"
	"twitter-clone/internal/infrastructure/middlewares"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/kataras/iris/v12"
)

type User struct {
	ID          uint64    `json:"id"`
	Username    string    `json:"username"`
	Nickname    string    `json:"nickname"`
	Name        string    `json:"name"`
	Website     string    `json:"website"`
	Description string    `json:"description"`
	Pic         string    `json:"pic"`
	PicCover    string    `json:"pic_cover"`
	Followings  []User    `json:"followings"`
	Followers   []User    `json:"followers"`
	CreatedAt   time.Time `json:"created_at"`
}

func MeHandler(ctx iris.Context, userRepo *dbrepository.UserRepository) {
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

	response.SendOkResponse(ctx, userDto(*user))
}

func UserHandler(ctx iris.Context, userRepo *dbrepository.UserRepository) {
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

	response.SendOkResponse(ctx, userDto(*user))
}

func FollowHandler(ctx iris.Context, userRepo *dbrepository.UserRepository) {
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

	response.SendOkResponse(ctx, nil)
}

func UnfollowHandler(ctx iris.Context, userRepo *dbrepository.UserRepository) {
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

	response.SendOkResponse(ctx, nil)
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
		ID:          user.ID,
		Name:        user.Name,
		Username:    user.Username,
		Nickname:    user.Username,
		Website:     user.Website,
		Description: user.Description,
		Pic:         user.Pic,
		PicCover:    user.PicCover,
		Followings:  followings,
		Followers:   followers,
		CreatedAt:   user.CreatedAt,
	}
}
