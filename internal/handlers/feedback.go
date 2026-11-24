package handlers

import (
	"context"
	"electrotech/internal/gen/handler_gen"
	"electrotech/internal/models"
	"net/http"
)

type FeedbackService interface {
	NewUserQuestion(ctx context.Context, q *models.CreateUserQuestion) error
}

// PostFeedback implements [handler_gen.Handler].
//
//nolint:nilerr
func (h *Handler) PostFeedback(
	ctx context.Context,
	req *handler_gen.PostFeedbackReq,
) (handler_gen.PostFeedbackRes, error) {
	err := h.feedback.NewUserQuestion(
		ctx,
		&models.CreateUserQuestion{
			PersonName: req.Name,
			Email:      req.Email.Or(""),
			Phone:      req.Phone.Or(""),
			Message:    req.Message,
		})
	if err != nil {
		return &handler_gen.BadResponseStatusCode{StatusCode: http.StatusInternalServerError, Response: handler_gen.Problem{
			Type:   "new-feedback-message",
			Status: http.StatusInternalServerError,
			Title:  "Create new feedback message",
			Detail: handler_gen.OptString{Value: "", Set: false},
		}}, nil
	}

	panic("unimplemented")
}
