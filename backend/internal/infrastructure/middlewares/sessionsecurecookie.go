package middlewares

import (
	"errors"
	"time"

	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
)

func SessionSecureCookieMiddleware(
	sessiondb *dbrepository.SessionRepository,
	sc *securecookie.SecureCookie,
) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		sessionID := ctx.GetCookie(cookieNameForSessionID, iris.CookieEncoding(sc))

		session, err := sessiondb.Get(ctx, dbrepository.GetSessionPayload{SessionID: sessionID})
		if err != nil {
			if errors.Is(err, dbrepository.ErrNotFound) {
				response.SendErrorResponse(ctx, iris.StatusUnauthorized, "Please, start your session.")
				return
			}
			response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
			return
		}

		ctx.Values().Set(UserIDKey, session.UserID)

		ctx.Next()
	}
}

func SessionSecureCookieLoginMiddleware(
	userRepo *dbrepository.UserRepository,
	sessionRepo *dbrepository.SessionRepository,
	sc *securecookie.SecureCookie,
) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		// Parse request body
		var login loginRequest
		err := ctx.ReadJSON(&login)
		if err != nil {
			response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		// Search user by username from request body
		user, err := userRepo.Get(ctx, dbrepository.GetUserPayload{Username: login.Username})
		if err != nil {
			if errors.Is(err, dbrepository.ErrNotFound) {
				response.SendErrorResponse(ctx, iris.StatusNotFound, err.Error())
				return
			}
			response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
			return
		}

		// Validate password from request body
		if err = bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password)); err != nil {
			response.SendErrorResponse(ctx, iris.StatusUnauthorized, err.Error())
			return
		}

		// Delete all active user's sessions
		err = sessionRepo.DeleteByUserID(ctx, user.ID)
		if err != nil {
			response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
			return
		}

		// Create new UUID session for user
		sessionID, err := uuid.NewUUID()
		if err != nil {
			response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
			return
		}

		// Set response secure cookie header
		ctx.SetCookie(
			&iris.Cookie{
				Name:  cookieNameForSessionID,
				Value: sessionID.String(),
			},
			iris.CookieEncoding(sc),
		)

		//  Save session_id value
		_, err = sessionRepo.Create(
			ctx,
			dbrepository.CreateSessionPayload{
				SessionID: sessionID.String(),
				UserID:    user.ID,
				ExpiredAt: time.Now().Add(time.Hour * 24), // TODO mb into env config (service logic)
			},
		)
		if err != nil {
			response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
			return
		}

		response.SendOkResponse(ctx, nil)
	}
}
