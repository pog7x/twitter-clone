package authservice

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"twitter-clone/config"
	"twitter-clone/internal/infrastructure/database"
	"twitter-clone/internal/repository/dbrepository"
	"twitter-clone/internal/services"
)

type AuthService struct {
	logger      *logrus.Logger
	sessionRepo *dbrepository.SessionRepository
	userRepo    *dbrepository.UserRepository

	sessionExpiredAt time.Duration
}

func NewAuthService(i *do.Injector) (*AuthService, error) {
	cfg := do.MustInvoke[*config.Config](i)

	return &AuthService{
		logger:      do.MustInvoke[*logrus.Logger](i),
		sessionRepo: do.MustInvoke[*dbrepository.SessionRepository](i),
		userRepo:    do.MustInvoke[*dbrepository.UserRepository](i),

		sessionExpiredAt: cfg.SessionExpiredAt,
	}, nil
}

type LoginPayload struct {
	Username, Password string
}

func (s *AuthService) AuthLogin(ctx context.Context, login LoginPayload) (*database.Session, error) {
	// Search user by username from request body
	user, err := s.userRepo.Get(ctx, dbrepository.GetUserPayload{Username: login.Username})
	if err != nil {
		if errors.Is(err, dbrepository.ErrNotFound) {
			return nil, services.NotFoundServiceError{Err: err, Entity: "user"}
		}
		return nil, services.InternalServiceError{Err: err}
	}

	// Validate password from request body
	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password)); err != nil {
		return nil, services.NotPermittedError{}
	}

	// Delete all active user's sessions
	err = s.sessionRepo.DeleteByUserID(ctx, user.ID)
	if err != nil {
		return nil, services.InternalServiceError{Err: err}
	}

	// Create new UUID session for user
	sessionID, err := uuid.NewUUID()
	if err != nil {
		return nil, services.InternalServiceError{Err: err}
	}

	//  Save session_id value
	session, err := s.sessionRepo.Create(
		ctx,
		dbrepository.CreateSessionPayload{
			SessionID: sessionID.String(),
			UserID:    user.ID,
			ExpiredAt: time.Now().Add(s.sessionExpiredAt),
		},
	)
	if err != nil {
		return nil, services.InternalServiceError{Err: err}
	}

	return session, nil
}

func (s *AuthService) GetSession(ctx context.Context, sessionID string) (*database.Session, error) {
	session, err := s.sessionRepo.Get(ctx, dbrepository.GetSessionPayload{SessionID: sessionID})
	if err != nil {
		if errors.Is(err, dbrepository.ErrNotFound) {
			return nil, services.NotPermittedError{}
		}
		return nil, services.InternalServiceError{Err: err}
	}

	return session, nil
}
