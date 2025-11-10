package ginger

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler func(*gin.Context) (any, error)

type Engine struct {
	isProduction bool
}

func Default() *Engine {
	isProd := gin.Mode() == gin.ReleaseMode
	return &Engine{
		isProduction: isProd,
	}
}

func (e *Engine) AsHandler(fn func(*gin.Context) (any, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := fn(ctx)

		if err != nil {
			e.handleErr(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, response)
	}
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message,omitzero"`
}

func (err *ErrorResponse) hideMessage() {
	err.Message = ""
}

func (e *Engine) NewErrorResponse(err Error) *ErrorResponse {
	resp := &ErrorResponse{
		StatusCode: err.statusCode,
		Message:    err.message,
	}
	if err.statusCode >= http.StatusInternalServerError && e.isProduction {
		resp.hideMessage()
	}

	return resp
}
func (e *Engine) NewUnknownErrorResponse(err error) *ErrorResponse {
	resp := &ErrorResponse{
		StatusCode: http.StatusInternalServerError,
	}
	if e.isProduction {
		resp.Message = "unknown internal server error occurred, see server logs for details"
	} else {

		resp.Message = fmt.Sprintf("unknown internal server error occurred: %s", err)
	}

	return resp
}

func (e *Engine) handleErr(ctx *gin.Context, err error) {
	if gingerErr, ok := err.(Error); ok {
		ctx.AbortWithStatusJSON(gingerErr.statusCode, e.NewErrorResponse(gingerErr))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusInternalServerError, e.NewUnknownErrorResponse(err))
}

type Error struct {
	statusCode int
	message    string
}

func (err Error) Error() string {
	return fmt.Sprintf("status code %d: %s", err.statusCode, err.message)
}

func BadRequest(msg string) Error {
	return Error{
		statusCode: http.StatusBadRequest,
		message:    msg,
	}
}

func Unauthorized(msg string) Error {
	return Error{
		statusCode: http.StatusUnauthorized,
		message:    msg,
	}
}
