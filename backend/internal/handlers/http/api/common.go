package api

import (
	"errors"

	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/infrastructure/auth"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/kataras/iris/v12"
)

// authenticatedUserID returns the user id put in the context by the auth
// middleware. It answers 401 and reports false when it is missing, so a handler
// can never fall through to work with a zero user id.
func authenticatedUserID(ctx iris.Context) (uint64, bool) {
	id, err := ctx.Values().GetUint64(auth.UserIDKey)
	if err != nil || id == 0 {
		response.SendErrorResponse(ctx, iris.StatusUnauthorized, "unauthenticated")
		return 0, false
	}

	return id, true
}

// viewerID returns the authenticated user id, or 0 when it is unavailable.
// Unlike authenticatedUserID it does not write a response: it is meant for
// read-only endpoints that merely personalize their output.
func viewerID(ctx iris.Context) uint64 {
	id, err := ctx.Values().GetUint64(auth.UserIDKey)
	if err != nil {
		return 0
	}

	return id
}

// pathID reads the {id} route parameter, answering 400 when it is malformed.
func pathID(ctx iris.Context) (uint64, bool) {
	id, err := ctx.Params().GetUint64("id")
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, "invalid id")
		return 0, false
	}

	return id, true
}

// readJSON decodes the request body, answering 400 when it is malformed. The
// decoder error itself is not echoed back to the client.
func readJSON(ctx iris.Context, dst any) bool {
	if err := ctx.ReadJSON(dst); err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, "invalid request body")
		return false
	}

	return true
}

// sendRepositoryError maps repository errors onto status codes. Anything that
// is not a known domain error is reported as a generic 500 and logged.
func sendRepositoryError(ctx iris.Context, err error) {
	switch {
	case errors.Is(err, dbrepository.ErrNotFound):
		response.SendErrorResponse(ctx, iris.StatusNotFound, "not found")
	case errors.Is(err, dbrepository.ErrSelfFollow):
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
	case errors.Is(err, dbrepository.ErrInvalidMedia):
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
	default:
		response.SendInternalError(ctx, err)
	}
}
