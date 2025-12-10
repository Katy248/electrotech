package v2

import (
	"electrotech/internal/repository/catalog"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func GetProducts(r *catalog.Repo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request Request
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}

		products, pages, err := r.GetProductsNew(
			catalog.Page(request.Page),
		)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
			})
			log.Error("Error getting products", "error", err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"products": products,
			"pages":    pages,
			"page":     request.Page,
		})
	}
}

type Request struct {
	Page int `form:"page" binding:"gte=0"`
}
