package dbrepository

import (
	"context"
	"errors"
	"twitter-clone/internal/infrastructure/database"

	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LikeRepository struct {
	db     *database.Database
	logger *logrus.Logger
}

func NewLikeDBRepository(i *do.Injector) (*LikeRepository, error) {
	database := do.MustInvoke[*database.Database](i)
	logger := do.MustInvoke[*logrus.Logger](i)

	return &LikeRepository{
		db:     database,
		logger: logger,
	}, nil
}

type CreateLikePayload struct {
	TweetID, UserID uint64
}

func (r *LikeRepository) Create(ctx context.Context, payload CreateLikePayload) (*database.Like, error) {
	like := database.Like{TweetID: payload.TweetID, UserID: payload.UserID}

	result := r.db.WithContext(ctx).Create(&like)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return &like, nil
}

type DeleteLikePayload struct {
	TweetID, UserID uint64
}

func (r *LikeRepository) Delete(ctx context.Context, payload DeleteLikePayload) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND tweet_id = ?", payload.UserID, payload.TweetID).
		Delete(&database.Like{})
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}

	return nil
}
