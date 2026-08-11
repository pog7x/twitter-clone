package auth

import (
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"

	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/services/authservice"
)

func SessionJWTMiddleware(authService *authservice.AuthService, verifier *jwt.Verifier) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		token := jwt.FromHeader(ctx)

		verifiedToken, err := verifier.VerifyToken([]byte(token))
		if err != nil {
			response.SendErrorResponse(ctx, iris.StatusUnauthorized, "invalid or expired token")
			return
		}

		var claims sessionClaims
		if err = verifiedToken.Claims(&claims); err != nil {
			response.SendErrorResponse(ctx, iris.StatusUnauthorized, "invalid token claims")
			return
		}

		session, err := authService.GetSession(ctx, claims.SessionID)
		if err != nil {
			response.HandleServiceError(ctx, err)
			return
		}

		ctx.Values().Set(UserIDKey, session.UserID)
		ctx.Values().Set(SessionIDKey, session.SessionID)

		ctx.Next()
	}
}

func SessionSecureCookieMiddleware(authService *authservice.AuthService, sc *securecookie.SecureCookie) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		sessionID := ctx.GetCookie(cookieNameForSessionID, iris.CookieEncoding(sc))

		session, err := authService.GetSession(ctx, sessionID)
		if err != nil {
			response.HandleServiceError(ctx, err)
			return
		}

		ctx.Values().Set(UserIDKey, session.UserID)
		ctx.Values().Set(SessionIDKey, session.SessionID)

		ctx.Next()
	}
}
