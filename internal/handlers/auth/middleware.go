package auth

import (
	"charm.land/log/v2"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(errNoHeader.Status, errNoHeader)

			return
		}

		claims, err := h.service.ValidateToken(tokenString)
		if err != nil {
			log.Error("Failed to validate token", "error", err)
			c.AbortWithStatusJSON(errInvalidToken.Status, errInvalidToken)

			return
		}

		c.Set("email", claims.Email)
		c.Set("user_id", claims.Id)
		c.Next()
	}
}
