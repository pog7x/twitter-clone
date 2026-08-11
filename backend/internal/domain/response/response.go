package response

import (
	"github.com/kataras/iris/v12"

	"twitter-clone/internal/services"
)

// InternalErrorMessage is what clients see for any 5xx. The real cause is only
// written to the server log, so database and driver details never leave the
// process.
const InternalErrorMessage = "internal server error"

type ErrorResponse struct {
	Success      bool   `json:"success"`
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`
}

type OkResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

// SendErrorResponse replies with a client-safe message. Never pass a raw
// error's text here for a 5xx: use SendInternalError instead.
func SendErrorResponse(ctx iris.Context, code int, ErrorMessage string) {
	ctx.StopWithStatus(code)
	_json(ctx, ErrorResponse{ErrorType: iris.StatusText(code), ErrorMessage: ErrorMessage})
}

// SendInternalError logs err and replies with a generic 500.
func SendInternalError(ctx iris.Context, err error) {
	if err != nil {
		ctx.Application().Logger().Errorf("%s %s: %v", ctx.Method(), ctx.Path(), err)
	}

	SendErrorResponse(ctx, iris.StatusInternalServerError, InternalErrorMessage)
}

func SendOkResponse(ctx iris.Context, result interface{}) {
	_json(ctx, OkResponse{Success: true, Result: result})
}

func _json(ctx iris.Context, v interface{}, opts ...iris.JSON) {
	err := ctx.JSON(v, opts...)
	if err != nil {
		if errHandler := ctx.Application().GetContextErrorHandler(); errHandler != nil {
			errHandler.HandleContextError(ctx, err)
		} else {
			ctx.Application().Logger().Error(err)
			ctx.StatusCode(iris.StatusInternalServerError)
		}
	}
}

func HandleServiceError(ctx iris.Context, err error) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case services.NotFoundServiceError:
		SendErrorResponse(ctx, iris.StatusNotFound, e.Entity+" not found")
	case services.InvalidCredentialsError:
		SendErrorResponse(ctx, iris.StatusUnauthorized, e.Error())
	case services.NotPermittedError:
		SendErrorResponse(ctx, iris.StatusForbidden, e.Error())
	default:
		SendInternalError(ctx, err)
	}
}
