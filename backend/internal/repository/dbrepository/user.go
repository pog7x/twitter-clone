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

type UserRepository struct {
	db     *database.Database
	logger *logrus.Logger
}

func NewUserDBRepository(i *do.Injector) (*UserRepository, error) {
	return &UserRepository{
		db:     do.MustInvoke[*database.Database](i),
		logger: do.MustInvoke[*logrus.Logger](i),
	}, nil
}

type CreateUserPayload struct {
	Name        string
	Password    []byte
	Username    string
	Website     string
	Pic         string
	PicCover    string
	Description string
}

func (r *UserRepository) Create(ctx context.Context, payload CreateUserPayload) (*database.User, error) {
	user := database.User{
		Name:        payload.Name,
		Password:    payload.Password,
		Username:    payload.Username,
		Website:     payload.Website,
		Pic:         payload.Pic,
		Description: payload.Description,
		PicCover:    payload.PicCover,
	}

	result := r.db.WithContext(ctx).Create(&user)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return &user, nil
}

type GetUserPayload struct {
	UserID   uint64
	Username string
}

func (r *UserRepository) Get(ctx context.Context, payload GetUserPayload) (*database.User, error) {
	var user database.User

	result := r.db.WithContext(ctx).
		Preload("Followings").
		Preload("Followers").
		Where(&database.User{ID: payload.UserID, Username: payload.Username}).
		First(&user)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	return &user, nil
}

type ListUserPayload struct {
	Limit, Offset int
}

func (r *UserRepository) List(ctx context.Context, payload ListUserPayload) ([]database.User, error) {
	var users []database.User

	result := r.db.WithContext(ctx).
		Preload("Followings").
		Preload("Followers").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}).
		Limit(payload.Limit).
		Offset(payload.Offset).
		Find(&users)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return users, nil
}

type UpdateUserPayload struct {
	Name        string
	Website     string
	Description string
	Followings  []*database.User
	Followers   []*database.User
}

func (r *UserRepository) Update(ctx context.Context, userID uint64, payload UpdateUserPayload) (*database.User, error) {
	user := database.User{
		ID:          userID,
		Name:        payload.Name,
		Website:     payload.Website,
		Description: payload.Description,
		Followings:  payload.Followings,
		Followers:   payload.Followers,
	}

	result := r.db.WithContext(ctx).Updates(&user)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	return &user, nil
}

func (r *UserRepository) Unfollow(ctx context.Context, followerID, followingID uint64) error {
	if err := r.db.WithContext(ctx).
		Model(&database.User{ID: followerID}).
		Association("Followings").
		Delete(&database.User{ID: followingID}); err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}

		return ErrInternal
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, userID uint64) error {
	result := r.db.WithContext(ctx).Delete(&database.User{ID: userID})
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}

	return nil
}
