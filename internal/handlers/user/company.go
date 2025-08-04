package user

import (
	"database/sql"
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

func UpdateCompanyData(repo *users.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateCompanyDataRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repo.GetByEmail(c.Request.Context(), c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		err = repo.UpdateCompanyData(c.Request.Context(), users.UpdateCompanyDataParams{
			Email:             c.GetString("email"),
			CompanyName:       sql.NullString{String: req.CompanyName, Valid: true},
			CompanyInn:        sql.NullString{String: req.CompanyINN, Valid: true},
			CompanyAddress:    sql.NullString{String: req.CompanyAddress, Valid: true},
			CompanyOkpo:       sql.NullString{String: req.CompanyOKPO, Valid: true},
			PositionInCompany: sql.NullString{String: req.PositionInCompany, Valid: true},
		})
		if err != nil {
			log.Printf("Error updating company data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update company data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Company data updated successfully"})
	}
}

func GetCompanyData(repo *users.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := repo.GetByEmail(c.Request.Context(), c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, newCompanyData(user))
	}
}

type CompanyData struct {
	CompanyName       string `json:"companyName"`
	CompanyINN        string `json:"companyINN"`
	CompanyAddress    string `json:"companyAddress"`
	CompanyOKPO       string `json:"companyOKPO"`
	PositionInCompany string `json:"positionInCompany"`

	AllRequiredFields bool `json:"allRequiredFields"`
}

func newCompanyData(u users.User) *CompanyData {
	data := &CompanyData{
		CompanyName:       u.CompanyName.String,
		CompanyINN:        u.CompanyInn.String,
		CompanyAddress:    u.CompanyAddress.String,
		PositionInCompany: u.PositionInCompany.String,
	}

	data.AllRequiredFields = CheckUserHasCompanyData(u)

	return data
}

func CheckUserHasCompanyData(user users.User) bool {
	return user.CompanyName.Valid &&
		user.CompanyAddress.Valid &&
		user.PositionInCompany.Valid &&
		user.CompanyInn.Valid &&
		user.CompanyOkpo.Valid
}
