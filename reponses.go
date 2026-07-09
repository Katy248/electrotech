package electrotech

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}

func Error(err error) ErrorResponse {
	return ErrorResponse{ErrorMessage: err.Error()}
}
func ErrorStr(err string) ErrorResponse {
	return ErrorResponse{ErrorMessage: err}
}
