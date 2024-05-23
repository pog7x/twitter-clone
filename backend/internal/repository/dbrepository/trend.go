package dbrepository

import (
	"context"

	"twitter-clone/internal/infrastructure/database"

	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

type TrendRepository struct {
	db     *database.Database
	logger *logrus.Logger
}

func NewTrendDBRepository(i *do.Injector) (*TrendRepository, error) {
	return &TrendRepository{
		db:     do.MustInvoke[*database.Database](i),
		logger: do.MustInvoke[*logrus.Logger](i),
	}, nil
}

type CreateTrendPayload struct {
	ID         uint64
	Name       string
	TweetCount uint64
}

func (r *TrendRepository) Create(ctx context.Context, payload CreateTrendPayload) (*database.Trend, error) {
	trend := database.Trend{
		ID:         payload.ID,
		Name:       payload.Name,
		TweetCount: payload.TweetCount,
	}

	result := r.db.WithContext(ctx).Create(&trend)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return &trend, nil
}

type ListTrendPayload struct {
	Limit, Offset int
	Name          string
}

func (r *TrendRepository) List(ctx context.Context, payload ListTrendPayload) ([]database.Trend, error) {
	var trends []database.Trend

	query := r.db.WithContext(ctx).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}).
		Limit(payload.Limit).
		Offset(payload.Offset)

	if payload.Name != "" {
		query = query.Where(&database.Trend{Name: payload.Name})
	}

	result := query.Find(&trends)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return trends, nil
}
