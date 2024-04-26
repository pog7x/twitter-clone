package dbrepository

import (
	"context"
	"errors"
	"time"
	"twitter-clone/internal/infrastructure/database"

	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db     *database.Database
	logger *logrus.Logger
}

func NewSessionDBRepository(i *do.Injector) (*SessionRepository, error) {
	database := do.MustInvoke[*database.Database](i)
	logger := do.MustInvoke[*logrus.Logger](i)

	return &SessionRepository{
		db:     database,
		logger: logger,
	}, nil
}

type CreateSessionPayload struct {
	SessionID string
	UserID    uint64
	ExpiredAt time.Time
}

func (r *SessionRepository) Create(ctx context.Context, payload CreateSessionPayload) (*database.Session, error) {
	session := database.Session{
		SessionID: payload.SessionID,
		UserID:    payload.UserID,
		ExpiredAt: payload.ExpiredAt,
	}

	result := r.db.WithContext(ctx).Create(&session)
	if err := result.Error; err != nil {
		return nil, ErrInternal
	}

	return &session, nil
}

type GetSessionPayload struct {
	SessionID string
}

func (r *SessionRepository) Get(ctx context.Context, payload GetSessionPayload) (*database.Session, error) {
	var session database.Session

	result := r.db.WithContext(ctx).
		Joins("User").
		Where("session_id = ? AND expired_at > ?", payload.SessionID, time.Now()).
		First(&session)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	return &session, nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	result := r.db.WithContext(ctx).Delete(&database.Session{SessionID: sessionID})
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}

	return nil
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uint64) error {
	result := r.db.WithContext(ctx).Where(&database.Session{UserID: userID}).Delete(&database.Session{})
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}

	return nil
}
