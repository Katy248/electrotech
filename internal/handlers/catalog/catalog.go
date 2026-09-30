package catalog

import (
	"electrotech"
	"electrotech/internal/models"
	"electrotech/internal/repository/catalog"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"charm.land/log/v2"
	"github.com/gin-gonic/gin"
)

func GetProducts(r *catalog.Repo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.Contains(ctx.Request.URL.String(), "filter") {
			log.Warn("Deprecated url, should be removed", "url", ctx.Request.URL.String())
		}

		log.Warn("Deprecated url, will be removed sooner, move to new version",
			"url", ctx.Request.URL.String(),
			"newVersion", "/api/v2/products")
		pageParam, _ := ctx.Params.Get("page")

		page, err := strconv.Atoi(pageParam)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, electrotech.Error(fmt.Errorf("failed parse page paraneter: %w", err)))

			return
		}

		products, err := r.GetProducts(
			catalog.Page(page),
		)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, electrotech.Error(err))
			log.Error("Error getting products", "error", err)

			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"products": products,
		})
	}
}

type Request struct {
	Page int `binding:"gte=0" json:"page"`
}

type Response struct {
	Code     int              `json:"code"`
	PageSize int              `json:"pageSize"`
	Page     int              `json:"page"`
	Products []models.Product `json:"products"`
}
