package api

import (
	"errors"

	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/infrastructure/auth"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/kataras/iris/v12"
)

func CreateLikeTweetHandler(ctx iris.Context, likeRepo *dbrepository.LikeRepository) {
	id, err := ctx.Params().GetUint64("id")
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	userID, err := ctx.Values().GetUint64(auth.UserIDKey)
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	_, err = likeRepo.Create(ctx, dbrepository.CreateLikePayload{TweetID: id, UserID: userID})
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	response.SendOkResponse(ctx, nil)
}

func DeleteLikeTweetHandler(ctx iris.Context, likeRepo *dbrepository.LikeRepository) {
	id, err := ctx.Params().GetUint64("id")
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	userID, err := ctx.Values().GetUint64(auth.UserIDKey)
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	err = likeRepo.Delete(ctx, dbrepository.DeleteLikePayload{TweetID: id, UserID: userID})
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
