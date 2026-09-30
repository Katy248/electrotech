package user

import (
	"electrotech"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateCompanyDataRequest struct {
	CompanyName       string `binding:"required" json:"company_name"`
	CompanyINN        string `binding:"required" json:"company_inn"`
	CompanyAddress    string `binding:"required" json:"company_address"`
	CompanyOKPO       string `binding:"required" json:"company_okpo"`
	PositionInCompany string `binding:"required" json:"position_in_company"`
}

func (h *Handler) HandleUpdateCompanyData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateCompanyDataRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		user, err := h.usersRepo.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			h.logger.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("invalid credentials"))

			return
		}

		user.CompanyName = &req.CompanyName
		user.CompanyInn = &req.CompanyINN
		user.CompanyAddress = &req.CompanyAddress
		user.PositionInCompany = &req.PositionInCompany
		user.CompanyOkpo = &req.CompanyOKPO

		err = h.usersRepo.Update(user)
		if err != nil {
			h.logger.Printf("Error updating company data: %v", err)
			c.JSON(http.StatusInternalServerError, electrotech.ErrorStr("failed to update company data"))

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Company data updated successfully"})
	}
}

func (h *Handler) HandleGetCompanyData() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := h.usersRepo.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			h.logger.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("invalid credentials"))

			return
		}

		c.JSON(http.StatusOK, user.CompanyData())
	}
}
