package electrotech

import (
	"electrotech/pkg/problem"
	"net/http"
)

func Error(err error) *problem.ProblemError {
	return problem.NewFrom(err, "err:electrotech/internal", http.StatusInternalServerError, err.Error())
}
func ErrorStr(err string) *problem.ProblemError {
	return problem.New("err:electrotech/internal", http.StatusInternalServerError, err)
}
