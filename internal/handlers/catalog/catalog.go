package catalog

import (
	"electrotech/internal/models"
	"electrotech/internal/repository/catalog"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProducts(r *catalog.Repo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
			log.Printf("Error getting products: %v", err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"products": products,
		})
	}
}
func GetProductsFilter(r *catalog.Repo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request Request
		if err := ctx.ShouldBindJSON(&request); err != nil {
			log.Printf("Error binding request: %v", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}

		filters := []catalog.FilterFunc{}

		for _, filter := range request.Filters {
			switch filter.Type {
			case models.ParameterTypeNumber:
				filters = append(filters, catalog.RangeFilter(filter.Parameter, filter.Min, filter.Max))
			case models.ParameterTypeList:
				filters = append(filters, catalog.ListFilter(filter.Parameter, filter.Values))
			}
		}

		products, err := r.GetProducts(
			catalog.Page(request.Page),
			filters...,
		)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
			})
			log.Printf("Error getting products: %v", err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"products": products,
		})
	}
}

type Request struct {
	Page    int             `json:"page" binding:"gte=0"`
	Filters []RequestFilter `json:"filters" binding:""`
}

type RequestFilter struct {
	Parameter string               `json:"parameter" binding:"required"`
	Type      models.ParameterType `json:"type" binding:"required"`
	Min       float64              `json:"min" binding:"gte=0"`
	Max       float64              `json:"max" binding:"gte=0"`
	Values    []string             `json:"values" binding:"required"`
}

type Response struct {
	Code     int              `json:"code"`
	PageSize int              `json:"pageSize"`
	Page     int              `json:"page"`
	Products []models.Product `json:"products"`
}
