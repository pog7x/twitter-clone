package dbrepository

import (
	"context"
	"errors"

	"twitter-clone/internal/infrastructure/database"

	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TweetRepository struct {
	db     *database.Database
	logger *logrus.Logger
}

func NewTweetDBRepository(i *do.Injector) (*TweetRepository, error) {
	return &TweetRepository{
		db:     do.MustInvoke[*database.Database](i),
		logger: do.MustInvoke[*logrus.Logger](i),
	}, nil
}

type CreateTweetPayload struct {
	AuthorID  uint64
	TweetData string
	MediaIDs  []uint64
}

// Create stores the tweet and attaches the uploads in one transaction. An
// upload is only attached when it belongs to the author and is not already
// used by another tweet; otherwise nothing is written at all.
func (r *TweetRepository) Create(ctx context.Context, payload CreateTweetPayload) (*database.Tweet, error) {
	tweet := database.Tweet{AuthorID: payload.AuthorID, TweetData: payload.TweetData}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&tweet).Error; err != nil {
			return ErrInternal
		}

		if len(payload.MediaIDs) == 0 {
			return nil
		}

		result := tx.Model(&database.Media{}).
			Where("id IN ? AND owner_id = ? AND tweet_id IS NULL", payload.MediaIDs, payload.AuthorID).
			Update("tweet_id", tweet.ID)
		if result.Error != nil {
			return ErrInternal
		}

		if result.RowsAffected != int64(len(payload.MediaIDs)) {
			return ErrInvalidMedia
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}

type GetTweetPayload struct {
	TweetID uint64
}

func (r *TweetRepository) Get(ctx context.Context, payload GetTweetPayload) (*database.Tweet, error) {
	var tweet database.Tweet

	result := r.db.WithContext(ctx).
		Joins("Author").
		Preload("TweetMedia").
		First(&tweet, payload.TweetID)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	return &tweet, nil
}

type ListTweetPayload struct {
	Limit, Offset int
	AuthorID      uint64
}

// List returns a page of tweets, newest first. Likes are not preloaded: their
// aggregates come from Stats so the payload stays flat.
func (r *TweetRepository) List(ctx context.Context, payload ListTweetPayload) ([]database.Tweet, error) {
	var tweets []database.Tweet

	result := r.db.WithContext(ctx).
		Joins("Author").
		Preload("TweetMedia").
		Where(&database.Tweet{AuthorID: payload.AuthorID}).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}).
		Limit(payload.Limit).
		Offset(payload.Offset).
		Find(&tweets)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return tweets, nil
}

// TweetStats holds the per-tweet aggregates rendered by the client.
type TweetStats struct {
	LikesCount int
	IsLiked    bool
}

// Stats returns like aggregates for the given tweets in two grouped queries,
// instead of loading every like row of every tweet.
func (r *TweetRepository) Stats(ctx context.Context, tweetIDs []uint64, viewerID uint64) (map[uint64]TweetStats, error) {
	stats := make(map[uint64]TweetStats, len(tweetIDs))

	if len(tweetIDs) == 0 {
		return stats, nil
	}

	var counts []struct {
		TweetID uint64
		Count   int
	}

	if err := r.db.WithContext(ctx).
		Model(&database.Like{}).
		Select("tweet_id, COUNT(*) AS count").
		Where("tweet_id IN ?", tweetIDs).
		Group("tweet_id").
		Find(&counts).Error; err != nil {
		return nil, ErrInternal
	}

	for _, c := range counts {
		stats[c.TweetID] = TweetStats{LikesCount: c.Count}
	}

	if viewerID != 0 {
		var likedIDs []uint64

		if err := r.db.WithContext(ctx).
			Model(&database.Like{}).
			Where("user_id = ? AND tweet_id IN ?", viewerID, tweetIDs).
			Distinct().
			Pluck("tweet_id", &likedIDs).Error; err != nil {
			return nil, ErrInternal
		}

		for _, id := range likedIDs {
			stat := stats[id]
			stat.IsLiked = true
			stats[id] = stat
		}
	}

	return stats, nil
}

type UpdateTweetPayload struct {
	TweetData string
	Likes     []database.Like
	Media     []database.Media
}

func (r *TweetRepository) Update(ctx context.Context, tweetID uint64, payload UpdateTweetPayload) (*database.Tweet, error) {
	tweet := database.Tweet{
		ID:         tweetID,
		TweetData:  payload.TweetData,
		TweetMedia: payload.Media,
		Likes:      payload.Likes,
	}

	result := r.db.WithContext(ctx).Updates(&tweet)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	return r.Get(ctx, GetTweetPayload{TweetID: tweetID})
}

func (r *TweetRepository) Delete(ctx context.Context, tweetID uint64) error {
	result := r.db.WithContext(ctx).Delete(&database.Tweet{ID: tweetID})
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}

	return nil
}
