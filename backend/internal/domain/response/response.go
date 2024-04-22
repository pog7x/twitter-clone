package response

import (
	"github.com/kataras/iris/v12"
)

type ErrorResponse struct {
	Result       bool   `json:"result"`
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`
}

func SendErrorResponse(ctx iris.Context, code int, ErrorMessage string) {
	ctx.StopWithStatus(code)
	JSON(ctx, ErrorResponse{ErrorType: iris.StatusText(code), ErrorMessage: ErrorMessage})
}

func JSON(ctx iris.Context, v interface{}, opts ...iris.JSON) {
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
