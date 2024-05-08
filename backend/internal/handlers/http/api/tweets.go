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

type Tweet struct {
	ID          uint64    `json:"id"`
	Content     string    `json:"content"`
	Attachments []string  `json:"attachments"`
	Author      User      `json:"author"`
	Likes       []Like    `json:"likes"`
	CreatedAt   time.Time `json:"created_at"`
}

type Like struct {
	UserID uint64 `json:"user_id"`
	Name   string `json:"name"`
}

type CreateTweet struct {
	TweetData     string   `json:"tweet_data"`
	TweetMediaIDs []uint64 `json:"tweet_media_ids"`
}

func CreateTweetHandler(ctx iris.Context, tweetRepo *dbrepository.TweetRepository) {
	userID, err := ctx.Values().GetUint64(middlewares.UserIDKey)
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	var request CreateTweet

	err = ctx.ReadJSON(&request)
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	var medias []database.Media
	for _, m := range request.TweetMediaIDs {
		medias = append(medias, database.Media{ID: m})
	}

	u, err := tweetRepo.Create(ctx, dbrepository.CreateTweetPayload{TweetData: request.TweetData, AuthorID: userID, Media: medias})
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	response.SendOkResponse(ctx, iris.Map{"tweet_id": u.ID})
}

func UpdateTweetHandler(ctx iris.Context, tweetRepo *dbrepository.TweetRepository) {
	id, err := ctx.Params().GetUint64("id")
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	var request CreateTweet

	err = ctx.ReadJSON(&request)
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	tweet, err := tweetRepo.Update(ctx, id, dbrepository.UpdateTweetPayload{
		TweetData: request.TweetData,
	})
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	response.SendOkResponse(ctx, tweetDto(*tweet))
}

func ListTweetHandler(ctx iris.Context, tweetRepo *dbrepository.TweetRepository) {
	tweets, err := tweetRepo.List(ctx, dbrepository.ListTweetPayload{Limit: 100})
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	var tweetsDto []Tweet

	for _, tweet := range tweets {
		tweetsDto = append(tweetsDto, tweetDto(tweet))
	}

	response.SendOkResponse(ctx, tweetsDto)
}

func UserTweetsHandler(ctx iris.Context, tweetRepo *dbrepository.TweetRepository) {
	id, err := ctx.Params().GetUint64("id")
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	tweets, err := tweetRepo.List(ctx, dbrepository.ListTweetPayload{Limit: 100, AuthorID: id})
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	var tweetsDto []Tweet

	for _, tweet := range tweets {
		tweetsDto = append(tweetsDto, tweetDto(tweet))
	}

	response.SendOkResponse(ctx, tweetsDto)
}

func DeleteTweetHandler(ctx iris.Context, tweetRepo *dbrepository.TweetRepository) {
	id, err := ctx.Params().GetUint64("id")
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	err = tweetRepo.Delete(ctx, id)
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

func tweetDto(tweet database.Tweet) Tweet {
	var likes []Like
	var links []string

	for _, like := range tweet.Likes {
		likes = append(likes, Like{UserID: like.UserID, Name: like.User.Name})
	}

	for _, media := range tweet.TweetMedia {
		links = append(links, media.Link)
	}

	tw := Tweet{
		ID:          tweet.ID,
		Content:     tweet.TweetData,
		CreatedAt:   tweet.CreatedAt,
		Likes:       likes,
		Attachments: links,
	}

	if tweet.Author != nil {
		tw.Author = userDto(*tweet.Author)
	}

	return tw
}
