package catalog

import (
	"electrotech/internal/models"
	"electrotech/internal/repository/catalog"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func GetProducts(r *catalog.Repo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.Contains(ctx.Request.URL.String(), "filter") {
			log.Warn("Deprecated url, should be removed", "url", ctx.Request.URL.String())
		}
		pageParam, _ := ctx.Params.Get("page")

		page, err := strconv.Atoi(pageParam)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}

		products, err := r.GetProducts(
			catalog.Page(page),
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
		})
	}
}

type Request struct {
	Page int `json:"page" binding:"gte=0"`
}

type Response struct {
	Code     int              `json:"code"`
	PageSize int              `json:"pageSize"`
	Page     int              `json:"page"`
	Products []models.Product `json:"products"`
}
