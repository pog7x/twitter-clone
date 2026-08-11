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
	ID          uint64
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
		ID:          payload.ID,
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

// Get loads a single user without its follow graph. Follower and following
// numbers come from Stats, so a profile response never carries nested users.
func (r *UserRepository) Get(ctx context.Context, payload GetUserPayload) (*database.User, error) {
	var user database.User

	result := r.db.WithContext(ctx).
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
	Username      string
}

func (r *UserRepository) List(ctx context.Context, payload ListUserPayload) ([]database.User, error) {
	var users []database.User

	query := r.db.WithContext(ctx).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}).
		Limit(payload.Limit).
		Offset(payload.Offset)

	if payload.Username != "" {
		query = query.Where(&database.User{Username: payload.Username})
	}

	result := query.Find(&users)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return users, nil
}

// UserStats holds the aggregates shown on a profile. They are computed with
// COUNT queries instead of loading the related rows.
type UserStats struct {
	FollowersCount  int
	FollowingsCount int
	TweetsCount     int
	IsFollowing     bool
}

// Stats returns the aggregates for userID. viewerID is the currently
// authenticated user and is used to resolve IsFollowing.
func (r *UserRepository) Stats(ctx context.Context, userID, viewerID uint64) (UserStats, error) {
	var stats UserStats

	followers := r.db.WithContext(ctx).Model(&database.User{ID: userID}).Association("Followers")
	if followers.Error != nil {
		return stats, ErrInternal
	}
	stats.FollowersCount = int(followers.Count())

	followings := r.db.WithContext(ctx).Model(&database.User{ID: userID}).Association("Followings")
	if followings.Error != nil {
		return stats, ErrInternal
	}
	stats.FollowingsCount = int(followings.Count())

	var tweetsCount int64
	if err := r.db.WithContext(ctx).
		Model(&database.Tweet{}).
		Where("author_id = ?", userID).
		Count(&tweetsCount).Error; err != nil {
		return stats, ErrInternal
	}
	stats.TweetsCount = int(tweetsCount)

	if viewerID != 0 && viewerID != userID {
		viewerFollowings := r.db.WithContext(ctx).
			Model(&database.User{ID: viewerID}).
			Where("users.id = ?", userID).
			Association("Followings")
		if viewerFollowings.Error != nil {
			return stats, ErrInternal
		}
		stats.IsFollowing = viewerFollowings.Count() > 0
	}

	return stats, nil
}

// ListFollowers returns the users that follow userID, newest first.
func (r *UserRepository) ListFollowers(ctx context.Context, userID uint64, limit, offset int) ([]database.User, error) {
	return r.listRelation(ctx, "Followers", userID, limit, offset)
}

// ListFollowings returns the users that userID follows, newest first.
func (r *UserRepository) ListFollowings(ctx context.Context, userID uint64, limit, offset int) ([]database.User, error) {
	return r.listRelation(ctx, "Followings", userID, limit, offset)
}

func (r *UserRepository) listRelation(ctx context.Context, relation string, userID uint64, limit, offset int) ([]database.User, error) {
	var users []database.User

	association := r.db.WithContext(ctx).
		Model(&database.User{ID: userID}).
		Order("users.id DESC").
		Limit(limit).
		Offset(offset).
		Association(relation)
	if association.Error != nil {
		return nil, ErrInternal
	}

	if err := association.Find(&users); err != nil {
		return nil, ErrInternal
	}

	return users, nil
}

type UpdateProfilePayload struct {
	Name        *string
	Website     *string
	Description *string
	Pic         *string
	PicCover    *string
}

// UpdateProfile applies only the fields present in the payload. It uses a map
// so that an explicitly empty string clears the column instead of being
// skipped as a zero value.
func (r *UserRepository) UpdateProfile(ctx context.Context, userID uint64, payload UpdateProfilePayload) (*database.User, error) {
	updates := make(map[string]any)

	if payload.Name != nil {
		updates["name"] = *payload.Name
	}
	if payload.Website != nil {
		updates["website"] = *payload.Website
	}
	if payload.Description != nil {
		updates["description"] = *payload.Description
	}
	if payload.Pic != nil {
		updates["pic"] = *payload.Pic
	}
	if payload.PicCover != nil {
		updates["pic_cover"] = *payload.PicCover
	}

	if len(updates) > 0 {
		result := r.db.WithContext(ctx).
			Model(&database.User{}).
			Where("id = ?", userID).
			Updates(updates)
		if err := result.Error; err != nil {
			return nil, ErrInternal
		}
		if result.RowsAffected == 0 {
			return nil, ErrNotFound
		}
	}

	return r.Get(ctx, GetUserPayload{UserID: userID})
}

// Follow makes followerID follow followingID. Re-following is a no-op because
// the join row insert ignores conflicts.
func (r *UserRepository) Follow(ctx context.Context, followerID, followingID uint64) error {
	if followerID == followingID {
		return ErrSelfFollow
	}

	if _, err := r.Get(ctx, GetUserPayload{UserID: followingID}); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).
		Model(&database.User{ID: followerID}).
		Association("Followings").
		Append(&database.User{ID: followingID}); err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}

		return ErrInternal
	}

	return nil
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
