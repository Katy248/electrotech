package catalog

import (
	repo "electrotech/internal/repository/catalog"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(r *repo.Repo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := r.GetProducts(0)

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
