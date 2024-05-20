package dbrepository

import (
	"context"

	"twitter-clone/internal/infrastructure/database"

	"github.com/samber/do"
	"github.com/sirupsen/logrus"
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
	Link string
}

func (r *MediaRepository) Create(ctx context.Context, payload CreateMediaPayload) (*database.Media, error) {
	media := database.Media{Link: payload.Link}

	result := r.db.WithContext(ctx).Create(&media)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return &media, nil
}
