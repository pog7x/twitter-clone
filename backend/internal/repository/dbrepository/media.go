package dbrepository

import (
	"context"
	"errors"

	"twitter-clone/internal/infrastructure/database"

	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MediaRepository struct {
	db     *database.Database
	logger *logrus.Logger
}

func NewMediaDBRepository(i *do.Injector) (*MediaRepository, error) {
	return &MediaRepository{
		db:     do.MustInvoke[*database.Database](i),
		logger: do.MustInvoke[*logrus.Logger](i),
	}, nil
}

type CreateMediaPayload struct {
	Link    string
	OwnerID uint64
}

func (r *MediaRepository) Create(ctx context.Context, payload CreateMediaPayload) (*database.Media, error) {
	media := database.Media{Link: payload.Link, OwnerID: payload.OwnerID}

	result := r.db.WithContext(ctx).Create(&media)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return &media, nil
}

// GetOwned loads an upload only if it belongs to ownerID, so a client cannot
// reference a file uploaded by somebody else.
func (r *MediaRepository) GetOwned(ctx context.Context, mediaID, ownerID uint64) (*database.Media, error) {
	var media database.Media

	result := r.db.WithContext(ctx).
		Where("id = ? AND owner_id = ?", mediaID, ownerID).
		First(&media)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidMedia
		}
		return nil, ErrInternal
	}

	return &media, nil
}
