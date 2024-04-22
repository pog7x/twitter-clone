package dbrepository

import "context"

type DBRepository interface {
	Create(ctx context.Context) (interface{}, error)
	Get(ctx context.Context) (interface{}, error)
	List(ctx context.Context) (interface{}, error)
	Update(ctx context.Context) (interface{}, error)
	Delete(ctx context.Context) (interface{}, error)
}
