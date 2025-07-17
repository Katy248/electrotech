package user

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ChangeEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ChangePhoneNumberRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type UpdateUserDataRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	Surname   string `json:"surname" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type UpdateCompanyDataRequest struct {
	CompanyName       string `json:"company_name" binding:"required"`
	CompanyINN        string `json:"company_inn" binding:"required"`
	CompanyAddress    string `json:"company_address" binding:"required"`
	PositionInCompany string `json:"position_in_company" binding:"required"`
}
