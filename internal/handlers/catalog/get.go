package catalog

import (
	"electrotech/internal/models"
	"electrotech/internal/repository/catalog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetProductRequest struct {
	ID string `uri:"id" binding:"required"`
}
type GetProductResponse struct {
	Code    int            `json:"code" binding:"required"`
	Product models.Product `json:"product" binding:"required"`
}

func GetProduct(repo *catalog.Repo) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var request GetProductRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid request body"})
			return
		}

		product, err := repo.GetProduct(request.ID)

		if err != nil {
			if err == catalog.ErrNotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "Product not found"})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Internal server error"})
			return
		}

		ctx.JSON(http.StatusOK, GetProductResponse{
			Code:    http.StatusOK,
			Product: product,
		})
	}
}
