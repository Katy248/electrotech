package common

import (
	"electrotech/pkg/problem"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrInternal = problem.New("err:electrotech/internal", http.StatusInternalServerError, "internal server error")

func SendError(ctx *gin.Context, err error) {
	if p, isProblem := errors.AsType[*problem.ProblemError](err); isProblem {
		sendProblem(ctx, p)
	} else {
		sendProblem(ctx, ErrInternal.WithDetails(err.Error()))
	}
}

func sendProblem(ctx *gin.Context, err *problem.ProblemError) {
	ctx.Header("content-type", "application/problem+json; charset=utf-8")
	ctx.Header("content-language", "en")
	ctx.JSON(err.Status, err)
}
