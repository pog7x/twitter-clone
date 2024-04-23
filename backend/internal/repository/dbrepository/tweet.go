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
	database := do.MustInvoke[*database.Database](i)
	logger := do.MustInvoke[*logrus.Logger](i)

	return &TweetRepository{
		db:     database,
		logger: logger,
	}, nil
}

type CreateTweetPayload struct {
	AuthorID  uint64
	TweetData string
	Media     []database.Media
}

func (r *TweetRepository) Create(ctx context.Context, payload CreateTweetPayload) (*database.Tweet, error) {
	tweet := database.Tweet{AuthorID: payload.AuthorID, TweetData: payload.TweetData, TweetMedia: payload.Media}

	result := r.db.WithContext(ctx).Create(&tweet)
	if err := result.Error; err != nil {
		return nil, ErrInternal
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
		Preload("Likes").
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
}

func (r *TweetRepository) List(ctx context.Context, payload ListTweetPayload) ([]database.Tweet, error) {
	var tweets []database.Tweet

	result := r.db.WithContext(ctx).
		Joins("Author").
		Preload("TweetMedia").
		Preload("Likes").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}).
		Limit(payload.Limit).
		Offset(payload.Offset).
		Find(&tweets)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return tweets, nil
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

	return &tweet, nil
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
