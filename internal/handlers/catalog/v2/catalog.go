package v2

import (
	"electrotech"
	"electrotech/internal/models"
	"electrotech/internal/repository/catalog"
	"net/http"

	"charm.land/log/v2"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *catalog.Repo
}

func NewHandler(repo *catalog.Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) HandleGetProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request Request
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		filters := []catalog.FilterFunc{}

		if request.Query != "" {
			filters = append(filters, catalog.QueryFilter(request.Query))
		}

		if request.OnlyAvailable {
			filters = append(filters, catalog.OnlyAvailableFilter())
		}

		products, err := h.repo.GetProductsNew(
			catalog.Page(request.Page), filters...,
		)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
			})
			log.Error("Error getting products", "error", err)

			return
		}

		ctx.JSON(http.StatusOK, Response{
			Code:     http.StatusOK,
			Products: products.Products,
			Pages:    products.Pages,
			Page:     products.Page,
			Total:    products.Total,
		})
	}
}

type Request struct {
	Page          int    `binding:"gte=0" form:"page"`
	Query         string `form:"query"`
	OnlyAvailable bool   `form:"oa"`
}

type Response struct {
	Code     int              `json:"code"`
	Products []models.Product `json:"products"`
	Pages    int              `json:"pages"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
}
