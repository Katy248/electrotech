package user

import (
	"electrotech/internal/repository/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(repo *users.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}
