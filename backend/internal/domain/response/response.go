package response

import "github.com/kataras/iris/v12"

type ErrorResponse struct {
	Success      bool   `json:"success"`
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`
}

type OkResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

func SendErrorResponse(ctx iris.Context, code int, ErrorMessage string) {
	ctx.StopWithStatus(code)
	JSON(ctx, ErrorResponse{ErrorType: iris.StatusText(code), ErrorMessage: ErrorMessage})
}

func SendOkResponse(ctx iris.Context, result interface{}) {
	JSON(ctx, OkResponse{Success: true, Result: result})
}

func JSON(ctx iris.Context, v interface{}, opts ...iris.JSON) {
	// ctx.ContentType(irisctx.ContentJSONHeaderValue)
	// err := irisctx.WriteJSON(ctx, v, &irisctx.DefaultJSONOptions)
	// err := json.NewEncoder(ctx.ResponseWriter()).Encode(v)
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
