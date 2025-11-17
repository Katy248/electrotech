package user

import (
	"electrotech/internal/repository/users"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateCompanyDataRequest struct {
	CompanyName       string `json:"company_name" binding:"required"`
	CompanyINN        string `json:"company_inn" binding:"required"`
	CompanyAddress    string `json:"company_address" binding:"required"`
	CompanyOKPO       string `json:"company_okpo" binding:"required"`
	PositionInCompany string `json:"position_in_company" binding:"required"`
}

func UpdateCompanyData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateCompanyDataRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		user.CompanyName = &req.CompanyName
		user.CompanyInn = &req.CompanyINN
		user.CompanyAddress = &req.CompanyAddress
		user.PositionInCompany = &req.PositionInCompany
		user.CompanyOkpo = &req.CompanyOKPO

		err = users.Update(user)

		if err != nil {
			log.Printf("Error updating company data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update company data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Company data updated successfully"})
	}
}

func GetCompanyData() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, user.CompanyData())
	}
}
