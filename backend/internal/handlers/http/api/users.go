package api

import (
	"context"
	"time"

	"twitter-clone/internal/domain/pagination"
	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/domain/validation"
	"twitter-clone/internal/infrastructure/database"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/kataras/iris/v12"
)

// UserBrief is the shape embedded into other resources, e.g. a tweet author or
// a row in a followers list. It never carries nested users.
type UserBrief struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Name      string    `json:"name"`
	Pic       string    `json:"pic"`
	PicCover  string    `json:"pic_cover"`
	CreatedAt time.Time `json:"created_at"`
}

// User is the full profile resource: the brief fields plus the aggregates and
// the viewer-relative flags the client renders.
type User struct {
	UserBrief
	Website         string `json:"website"`
	Description     string `json:"description"`
	FollowersCount  int    `json:"followers_count"`
	FollowingsCount int    `json:"followings_count"`
	TweetsCount     int    `json:"tweets_count"`
	IsFollowing     bool   `json:"is_following"`
	IsMe            bool   `json:"is_me"`
}

func MeHandler(ctx iris.Context, userRepo *dbrepository.UserRepository) {
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	sendUserProfile(ctx, userRepo, userID, userID)
}

// UpdateMe uses pointers so that a field left out of the request body keeps its
// current value, while an explicitly empty one clears the column. Images are
// referenced by the media id returned from the upload endpoint, never by a
// client-supplied path.
type UpdateMe struct {
	Name            *string `json:"name"`
	Website         *string `json:"website"`
	Description     *string `json:"description"`
	PicMediaID      *uint64 `json:"pic_media_id"`
	PicCoverMediaID *uint64 `json:"pic_cover_media_id"`
}

func MeUpdateHandler(
	ctx iris.Context,
	userRepo *dbrepository.UserRepository,
	mediaRepo *dbrepository.MediaRepository,
) {
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	var request UpdateMe

	if !readJSON(ctx, &request) {
		return
	}

	payload, err := buildUpdateProfilePayload(request)
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	// Resolve avatar and cover to stored links, checking that the uploads
	// belong to the caller.
	if request.PicMediaID != nil {
		link, resolveErr := resolveOwnedMediaLink(ctx, mediaRepo, *request.PicMediaID, userID)
		if resolveErr != nil {
			sendRepositoryError(ctx, resolveErr)
			return
		}
		payload.Pic = &link
	}

	if request.PicCoverMediaID != nil {
		link, resolveErr := resolveOwnedMediaLink(ctx, mediaRepo, *request.PicCoverMediaID, userID)
		if resolveErr != nil {
			sendRepositoryError(ctx, resolveErr)
			return
		}
		payload.PicCover = &link
	}

	if _, err = userRepo.UpdateProfile(ctx, userID, payload); err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	sendUserProfile(ctx, userRepo, userID, userID)
}

func resolveOwnedMediaLink(
	ctx iris.Context,
	mediaRepo *dbrepository.MediaRepository,
	mediaID, ownerID uint64,
) (string, error) {
	media, err := mediaRepo.GetOwned(ctx, mediaID, ownerID)
	if err != nil {
		return "", err
	}

	return media.Link, nil
}

func buildUpdateProfilePayload(request UpdateMe) (dbrepository.UpdateProfilePayload, error) {
	var payload dbrepository.UpdateProfilePayload

	if request.Name != nil {
		name, err := validation.Name(*request.Name)
		if err != nil {
			return payload, err
		}
		payload.Name = &name
	}

	if request.Description != nil {
		description, err := validation.Description(*request.Description)
		if err != nil {
			return payload, err
		}
		payload.Description = &description
	}

	if request.Website != nil {
		website, err := validation.Website(*request.Website)
		if err != nil {
			return payload, err
		}
		payload.Website = &website
	}

	return payload, nil
}

func UserHandler(ctx iris.Context, userRepo *dbrepository.UserRepository) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}

	sendUserProfile(ctx, userRepo, id, viewerID(ctx))
}

func FollowHandler(ctx iris.Context, userRepo *dbrepository.UserRepository) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}

	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	if err := userRepo.Follow(ctx, userID, id); err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	sendUserProfile(ctx, userRepo, id, userID)
}

func UnfollowHandler(ctx iris.Context, userRepo *dbrepository.UserRepository) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}

	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	if err := userRepo.Unfollow(ctx, userID, id); err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	sendUserProfile(ctx, userRepo, id, userID)
}

func FollowersHandler(ctx iris.Context, userRepo *dbrepository.UserRepository) {
	listRelationHandler(ctx, userRepo.ListFollowers)
}

func FollowingsHandler(ctx iris.Context, userRepo *dbrepository.UserRepository) {
	listRelationHandler(ctx, userRepo.ListFollowings)
}

type relationLister func(ctx context.Context, userID uint64, limit, offset int) ([]database.User, error)

func listRelationHandler(ctx iris.Context, list relationLister) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}

	page := pagination.FromContext(ctx)

	users, err := list(ctx, id, page.Limit, page.Offset)
	if err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	usersDto := make([]UserBrief, 0, len(users))
	for _, user := range users {
		usersDto = append(usersDto, userBriefDto(user))
	}

	response.SendOkResponse(ctx, usersDto)
}

func sendUserProfile(ctx iris.Context, userRepo *dbrepository.UserRepository, userID, viewer uint64) {
	user, err := userRepo.Get(ctx, dbrepository.GetUserPayload{UserID: userID})
	if err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	stats, err := userRepo.Stats(ctx, userID, viewer)
	if err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	response.SendOkResponse(ctx, userDto(*user, stats, viewer))
}

func userBriefDto(user database.User) UserBrief {
	return UserBrief{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  "@" + user.Username,
		Name:      user.Name,
		Pic:       user.Pic,
		PicCover:  user.PicCover,
		CreatedAt: user.CreatedAt,
	}
}

func userDto(user database.User, stats dbrepository.UserStats, viewer uint64) User {
	return User{
		UserBrief:       userBriefDto(user),
		Website:         user.Website,
		Description:     user.Description,
		FollowersCount:  stats.FollowersCount,
		FollowingsCount: stats.FollowingsCount,
		TweetsCount:     stats.TweetsCount,
		IsFollowing:     stats.IsFollowing,
		IsMe:            viewer != 0 && viewer == user.ID,
	}
}
