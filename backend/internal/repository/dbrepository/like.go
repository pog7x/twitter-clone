package dbrepository

import (
	"context"

	"twitter-clone/internal/infrastructure/database"

	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

type LikeRepository struct {
	db     *database.Database
	logger *logrus.Logger
}

func NewLikeDBRepository(i *do.Injector) (*LikeRepository, error) {
	return &LikeRepository{
		db:     do.MustInvoke[*database.Database](i),
		logger: do.MustInvoke[*logrus.Logger](i),
	}, nil
}

type CreateLikePayload struct {
	TweetID, UserID uint64
}

// Create is idempotent: liking an already liked tweet hits the unique
// (user_id, tweet_id) index and is silently ignored instead of adding a
// duplicate row that would inflate the counter.
func (r *LikeRepository) Create(ctx context.Context, payload CreateLikePayload) error {
	like := database.Like{TweetID: payload.TweetID, UserID: payload.UserID}

	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "tweet_id"}},
			DoNothing: true,
		}).
		Create(&like)
	if err := result.Error; err != nil {
		return ErrInternal
	}

	return nil
}

type DeleteLikePayload struct {
	TweetID, UserID uint64
}

// Delete is idempotent as well: removing a like that is not there is a no-op.
func (r *LikeRepository) Delete(ctx context.Context, payload DeleteLikePayload) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND tweet_id = ?", payload.UserID, payload.TweetID).
		Delete(&database.Like{})
	if err := result.Error; err != nil {
		return ErrInternal
	}

	return nil
}
