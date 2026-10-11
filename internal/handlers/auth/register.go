package auth

import (
	"electrotech/internal/handlers/common"
	"electrotech/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			common.SendError(c, errBadRequest.WithDetails(err.Error()))

			return
		}

		err := h.service.Register(c.Request.Context(), &req)
		if err != nil {
			common.SendError(c, err)

			return
		}

		c.Status(http.StatusCreated)
	}
}
