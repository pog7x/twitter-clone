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
	// Search user by username from payload. A missing user and a wrong password
	// produce the same error so the endpoint cannot be used to probe for
	// existing usernames.
	user, err := s.userRepo.Get(ctx, dbrepository.GetUserPayload{Username: login.Username})
	if err != nil {
		if errors.Is(err, dbrepository.ErrNotFound) {
			return nil, services.InvalidCredentialsError{}
		}
		return nil, services.InternalServiceError{Err: err}
	}

	// Validate password from request body
	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password)); err != nil {
		return nil, services.InvalidCredentialsError{}
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

// AuthLogout invalidates the session behind the current token, so a logged out
// token cannot be reused until it expires on its own.
func (s *AuthService) AuthLogout(ctx context.Context, sessionID string) error {
	if err := s.sessionRepo.Delete(ctx, sessionID); err != nil {
		if errors.Is(err, dbrepository.ErrNotFound) {
			return nil
		}
		return services.InternalServiceError{Err: err}
	}

	return nil
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
