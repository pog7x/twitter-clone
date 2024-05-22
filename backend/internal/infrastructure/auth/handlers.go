package auth

import (
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"

	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/services/authservice"
)

func SessionJWTLoginHandler(authService *authservice.AuthService, signer *jwt.Signer) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		// Parse request body
		var login loginRequest
		err := ctx.ReadJSON(&login)
		if err != nil {
			response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		session, err := authService.AuthLogin(ctx, authservice.LoginPayload{
			Username: login.Username,
			Password: login.Password,
		})
		if err != nil {
			response.HandleServiceError(ctx, err)
			return
		}

		// Create JWT token for user session
		token, err := signer.Sign(sessionClaims{SessionID: session.SessionID})
		if err != nil {
			response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
			return
		}

		response.SendOkResponse(ctx, string(token))
	}
}

func SessionSecureCookieLoginHandler(authService *authservice.AuthService, sc *securecookie.SecureCookie) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		// Parse request body
		var login loginRequest
		err := ctx.ReadJSON(&login)
		if err != nil {
			response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
			return
		}

		session, err := authService.AuthLogin(ctx, authservice.LoginPayload{
			Username: login.Username,
			Password: login.Password,
		})
		if err != nil {
			response.HandleServiceError(ctx, err)
			return
		}

		// Set response secure cookie header
		ctx.SetCookie(
			&iris.Cookie{
				Name:  cookieNameForSessionID,
				Value: session.SessionID,
			},
			iris.CookieEncoding(sc),
		)

		response.SendOkResponse(ctx, nil)
	}
}
