package api

import (
	"time"

	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/infrastructure/database"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/kataras/iris/v12"
)

type Trend struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	TweetsCount uint64    `json:"tweets_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func ListTrendHandler(ctx iris.Context, trendRepo *dbrepository.TrendRepository) {
	trends, err := trendRepo.List(ctx, dbrepository.ListTrendPayload{
		Limit: 100,
	})
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	var trendsDto []Trend

	for _, trend := range trends {
		trendsDto = append(trendsDto, trendDto(trend))
	}

	response.SendOkResponse(ctx, trendsDto)
}

func trendDto(trend database.Trend) Trend {
	return Trend{
		ID:          trend.ID,
		Name:        trend.Name,
		TweetsCount: trend.TweetCount,
		CreatedAt:   trend.CreatedAt,
	}
}
