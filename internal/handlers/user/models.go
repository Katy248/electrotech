package user

type ChangePasswordRequest struct {
	OldPassword string `binding:"required"       json:"old_password"`
	NewPassword string `binding:"required,min=8" json:"new_password"`
}

type ChangeEmailRequest struct {
	Email string `binding:"required,email" json:"email"`
}

type ChangePhoneNumberRequest struct {
	PhoneNumber string `binding:"required" json:"phone_number"`
}

type UpdateUserDataRequest struct {
	FirstName string `binding:"required" json:"first_name"`
	Surname   string `binding:"required" json:"surname"`
	LastName  string `binding:"required" json:"last_name"`
}
