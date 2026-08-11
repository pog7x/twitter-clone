package api

import (
	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/kataras/iris/v12"
)

// TweetLikes is returned by the like endpoints so the client can update a
// single card in place instead of refetching the whole feed.
type TweetLikes struct {
	TweetID    uint64 `json:"tweet_id"`
	LikesCount int    `json:"likes_count"`
	IsLiked    bool   `json:"is_liked"`
}

func CreateLikeTweetHandler(
	ctx iris.Context,
	likeRepo *dbrepository.LikeRepository,
	tweetRepo *dbrepository.TweetRepository,
) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}

	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	// Checked up front so a missing tweet answers 404 instead of tripping the
	// foreign key and surfacing as a 500.
	if _, err := tweetRepo.Get(ctx, dbrepository.GetTweetPayload{TweetID: id}); err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	if err := likeRepo.Create(ctx, dbrepository.CreateLikePayload{TweetID: id, UserID: userID}); err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	sendTweetLikes(ctx, tweetRepo, id, userID)
}

func DeleteLikeTweetHandler(
	ctx iris.Context,
	likeRepo *dbrepository.LikeRepository,
	tweetRepo *dbrepository.TweetRepository,
) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}

	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	if err := likeRepo.Delete(ctx, dbrepository.DeleteLikePayload{TweetID: id, UserID: userID}); err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	sendTweetLikes(ctx, tweetRepo, id, userID)
}

func sendTweetLikes(ctx iris.Context, tweetRepo *dbrepository.TweetRepository, tweetID, viewer uint64) {
	stats, err := tweetRepo.Stats(ctx, []uint64{tweetID}, viewer)
	if err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	response.SendOkResponse(ctx, TweetLikes{
		TweetID:    tweetID,
		LikesCount: stats[tweetID].LikesCount,
		IsLiked:    stats[tweetID].IsLiked,
	})
}
