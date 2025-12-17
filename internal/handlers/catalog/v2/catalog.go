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

		filters := []catalog.FilterFunc{}

		if request.Query != "" {
			filters = append(filters, catalog.QueryFilter(request.Query))
		}
		if request.OnlyAvailable {
			filters = append(filters, catalog.OnlyAvailableFilter())
		}

		products, err := r.GetProductsNew(
			catalog.Page(request.Page), filters...,
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
			"products": products.Products,
			"pages":    products.Pages,
			"total":    products.Total,
			"page":     products.Page,
		})

	}
}

type Request struct {
	Page          int    `form:"page" binding:"gte=0"`
	Query         string `form:"query"`
	OnlyAvailable bool   `form:"oa"`
}
