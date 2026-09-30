package catalog

import (
	"electrotech"
	"electrotech/internal/models"
	"electrotech/internal/repository/catalog"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetProductRequest struct {
	ID string `binding:"required" uri:"id"`
}
type GetProductResponse struct {
	Code    int            `binding:"required" json:"code"`
	Product models.Product `binding:"required" json:"product"`
}

func (h *Handler) HandleGetProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request GetProductRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		product, err := h.repo.GetProduct(request.ID)
		if err != nil {
			if errors.Is(err, catalog.ErrNotFound) {
				ctx.AbortWithStatusJSON(http.StatusNotFound, electrotech.Error(err))

				return
			}

			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal server error",
				})

			return
		}

		ctx.JSON(http.StatusOK, GetProductResponse{
			Code:    http.StatusOK,
			Product: product,
		})
	}
}
