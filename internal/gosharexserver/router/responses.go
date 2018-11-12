package router

import (
	"github.com/kataras/iris"
	"github.com/labstack/gommon/log"
	"net/http"
)

func Unauthorized(ctx iris.Context) {
	HttpError(ctx, http.StatusUnauthorized, "unauthorized request")
}

func InternalServerError(ctx iris.Context) {
	HttpError(ctx, http.StatusUnauthorized, "internal server error")
}

func HttpError(ctx iris.Context, statusCode int, errMessage string) {
	ctx.StatusCode(statusCode)
	if _, err := ctx.JSON(iris.Map{"error_message": errMessage}); err != nil {
		log.Printf("an error occurred while writing the error message to %s: %e", ctx.RemoteAddr(), err)
	}
	ctx.StopExecution()
}
