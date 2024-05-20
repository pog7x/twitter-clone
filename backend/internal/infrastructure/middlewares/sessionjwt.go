package middlewares

import (
	"errors"
	"time"

	"twitter-clone/internal/domain/response"

	"twitter-clone/internal/repository/dbrepository"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"golang.org/x/crypto/bcrypt"
)

func SessionJWTLoginMiddleware(
	userRepo *dbrepository.UserRepository,
	sessionRepo *dbrepository.SessionRepository,
	signer *jwt.Signer,
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

		// Create new UUID session_id for user
		sessionID, err := uuid.NewUUID()
		if err != nil {
			response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
			return
		}

		// Create JWT token for user session
		token, err := signer.Sign(sessionClaims{SessionID: sessionID.String()})
		if err != nil {
			response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
			return
		}

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

		response.SendOkResponse(ctx, string(token))
	}
}

func SessionJWTMiddleware(
	sessiondb *dbrepository.SessionRepository,
	verifier *jwt.Verifier,
) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		token := jwt.FromHeader(ctx)

		verifiedToken, err := verifier.VerifyToken([]byte(token))
		if err != nil {
			response.SendErrorResponse(ctx, iris.StatusUnauthorized, err.Error())
			return
		}

		var claims sessionClaims
		if err = verifiedToken.Claims(&claims); err != nil {
			response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		session, err := sessiondb.Get(ctx, dbrepository.GetSessionPayload{SessionID: claims.SessionID})
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
