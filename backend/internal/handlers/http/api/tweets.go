package api

import (
	"time"

	"twitter-clone/internal/domain/pagination"
	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/domain/validation"
	"twitter-clone/internal/infrastructure/database"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/kataras/iris/v12"
)

type Tweet struct {
	ID          uint64    `json:"id"`
	Content     string    `json:"content"`
	Attachments []string  `json:"attachments"`
	Author      UserBrief `json:"author"`
	LikesCount  int       `json:"likes_count"`
	IsLiked     bool      `json:"is_liked"`
	IsOwner     bool      `json:"is_owner"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateTweet struct {
	TweetData     string   `json:"tweet_data"`
	TweetMediaIDs []uint64 `json:"tweet_media_ids"`
}

func CreateTweetHandler(ctx iris.Context, tweetRepo *dbrepository.TweetRepository) {
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	var request CreateTweet

	if !readJSON(ctx, &request) {
		return
	}

	content, err := validation.TweetContent(request.TweetData, len(request.TweetMediaIDs))
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	tweet, err := tweetRepo.Create(ctx, dbrepository.CreateTweetPayload{
		TweetData: content,
		AuthorID:  userID,
		MediaIDs:  request.TweetMediaIDs,
	})
	if err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	// Re-read the tweet so the response carries the author and the attachment
	// links, matching the shape returned by the list endpoints.
	sendTweet(ctx, tweetRepo, tweet.ID, userID)
}

func GetTweetHandler(ctx iris.Context, tweetRepo *dbrepository.TweetRepository) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}

	sendTweet(ctx, tweetRepo, id, viewerID(ctx))
}

func UpdateTweetHandler(ctx iris.Context, tweetRepo *dbrepository.TweetRepository) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}

	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	var request CreateTweet

	if !readJSON(ctx, &request) {
		return
	}

	existing, ok := ownedTweet(ctx, tweetRepo, id, userID)
	if !ok {
		return
	}

	content, err := validation.TweetContent(request.TweetData, len(existing.TweetMedia))
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	if _, err = tweetRepo.Update(ctx, id, dbrepository.UpdateTweetPayload{TweetData: content}); err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	sendTweet(ctx, tweetRepo, id, userID)
}

func DeleteTweetHandler(ctx iris.Context, tweetRepo *dbrepository.TweetRepository) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}

	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	if _, ok = ownedTweet(ctx, tweetRepo, id, userID); !ok {
		return
	}

	if err := tweetRepo.Delete(ctx, id); err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	response.SendOkResponse(ctx, nil)
}

func ListTweetHandler(ctx iris.Context, tweetRepo *dbrepository.TweetRepository) {
	page := pagination.FromContext(ctx)

	sendTweetList(ctx, tweetRepo, dbrepository.ListTweetPayload{
		Limit:  page.Limit,
		Offset: page.Offset,
	})
}

func UserTweetsHandler(ctx iris.Context, tweetRepo *dbrepository.TweetRepository) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}

	page := pagination.FromContext(ctx)

	sendTweetList(ctx, tweetRepo, dbrepository.ListTweetPayload{
		Limit:    page.Limit,
		Offset:   page.Offset,
		AuthorID: id,
	})
}

// ownedTweet loads a tweet and asserts that userID is its author. It answers
// 404 when the tweet is gone and 403 when it belongs to somebody else, so
// editing and deleting are limited to the author.
func ownedTweet(
	ctx iris.Context,
	tweetRepo *dbrepository.TweetRepository,
	tweetID, userID uint64,
) (*database.Tweet, bool) {
	tweet, err := tweetRepo.Get(ctx, dbrepository.GetTweetPayload{TweetID: tweetID})
	if err != nil {
		sendRepositoryError(ctx, err)
		return nil, false
	}

	if tweet.AuthorID != userID {
		response.SendErrorResponse(ctx, iris.StatusForbidden, "tweet belongs to another user")
		return nil, false
	}

	return tweet, true
}

func sendTweet(ctx iris.Context, tweetRepo *dbrepository.TweetRepository, tweetID, viewer uint64) {
	tweet, err := tweetRepo.Get(ctx, dbrepository.GetTweetPayload{TweetID: tweetID})
	if err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	stats, err := tweetRepo.Stats(ctx, []uint64{tweetID}, viewer)
	if err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	response.SendOkResponse(ctx, tweetDto(*tweet, stats[tweetID], viewer))
}

func sendTweetList(ctx iris.Context, tweetRepo *dbrepository.TweetRepository, payload dbrepository.ListTweetPayload) {
	tweets, err := tweetRepo.List(ctx, payload)
	if err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	viewer := viewerID(ctx)

	tweetIDs := make([]uint64, 0, len(tweets))
	for _, tweet := range tweets {
		tweetIDs = append(tweetIDs, tweet.ID)
	}

	stats, err := tweetRepo.Stats(ctx, tweetIDs, viewer)
	if err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	tweetsDto := make([]Tweet, 0, len(tweets))
	for _, tweet := range tweets {
		tweetsDto = append(tweetsDto, tweetDto(tweet, stats[tweet.ID], viewer))
	}

	response.SendOkResponse(ctx, tweetsDto)
}

func tweetDto(tweet database.Tweet, stats dbrepository.TweetStats, viewer uint64) Tweet {
	links := make([]string, 0, len(tweet.TweetMedia))
	for _, media := range tweet.TweetMedia {
		links = append(links, media.Link)
	}

	tw := Tweet{
		ID:          tweet.ID,
		Content:     tweet.TweetData,
		CreatedAt:   tweet.CreatedAt,
		Attachments: links,
		LikesCount:  stats.LikesCount,
		IsLiked:     stats.IsLiked,
		IsOwner:     viewer != 0 && viewer == tweet.AuthorID,
	}

	if tweet.Author != nil {
		tw.Author = userBriefDto(*tweet.Author)
	}

	return tw
}
