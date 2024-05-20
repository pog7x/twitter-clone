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

func (r *TrendRepository) List(ctx context.Context) ([]database.Trend, error) {
	var trends []database.Trend

	result := r.db.WithContext(ctx).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}).
		Find(&trends)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return trends, nil
}
