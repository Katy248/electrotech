package filter

import (
	"electrotech/internal/repository/catalog"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFilters(repo *catalog.Repo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		parameters, err := repo.GetParameters()
		if err != nil {
			log.Printf("Failed get parameters: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":       http.StatusOK,
			"parameters": parameters,
		})
	}
}
