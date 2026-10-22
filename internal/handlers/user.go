package handlers

import (
	"context"
	"electrotech/internal/gen/handler_gen"
	"electrotech/internal/models"
	"fmt"
)

// GetUser implements [handler_gen.Handler].
//
//nolint:ireturn
func (h *Handler) GetUser(ctx context.Context) (handler_gen.GetUserRes, error) {
	userID, _ := GetUserID(ctx)

	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return userToResponse(user), nil
}

// UpdateUser implements [handler_gen.Handler].
//
//nolint:ireturn
func (h *Handler) UpdateUser(ctx context.Context, req *handler_gen.UpdateUserReq) (handler_gen.UpdateUserRes, error) {
	update := models.UserUpdate{
		FirstName: optStringToNilString(req.FirstName),
		LastName:  optStringToNilString(req.LastName),
		Surname:   optStringToNilString(req.Surname),
	}

	if req.Company.Set {
		update.CompanyAddress = optStringToNilString(req.Company.Value.Address)
		update.CompanyInn = optStringToNilString(req.Company.Value.Inn)
		update.CompanyName = optStringToNilString(req.Company.Value.Name)
		update.CompanyOkpo = optStringToNilString(req.Company.Value.Okpo)
		update.PositionInCompany = optStringToNilString(req.Company.Value.Position)
	}

	userID, _ := GetUserID(ctx)

	err := h.userService.Update(ctx, userID, &update)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get updated user: %w", err)
	}

	return userToResponse(user), nil
}

// UpdateUserEmail implements [handler_gen.Handler].
//
//nolint:ireturn
func (h *Handler) UpdateUserEmail(
	ctx context.Context,
	req *handler_gen.UpdateUserEmailReq,
) (handler_gen.UpdateUserEmailRes, error) {
	userId, _ := GetUserID(ctx)

	err := h.userService.UpdateEmail(ctx, userId, req.Email)
	if err != nil {
		return nil, fmt.Errorf("update email: %w", err)
	}

	return &handler_gen.UpdateUserEmailOK{}, nil
}

// UpdateUserPassword implements [handler_gen.Handler].
//
//nolint:ireturn
func (h *Handler) UpdateUserPassword(
	ctx context.Context,
	req *handler_gen.UpdateUserPasswordReq,
) (handler_gen.UpdateUserPasswordRes, error) {
	userId, _ := GetUserID(ctx)

	err := h.userService.UpdatePassword(ctx, userId, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("update password: %w", err)
	}

	return &handler_gen.UpdateUserPasswordOK{}, nil
}

// UpdateUserPhone implements [handler_gen.Handler].
//
//nolint:ireturn
func (h *Handler) UpdateUserPhone(
	ctx context.Context,
	req *handler_gen.UpdateUserPhoneReq,
) (handler_gen.UpdateUserPhoneRes, error) {
	userId, _ := GetUserID(ctx)

	err := h.userService.UpdatePhone(ctx, userId, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("update phone: %w", err)
	}

	return &handler_gen.UpdateUserPhoneOK{}, nil
}

func optStringToNilString(opt handler_gen.OptString) *string {
	if opt.Set {
		return &opt.Value
	}

	return nil
}

func userToResponse(user *models.User) *handler_gen.User {
	return &handler_gen.User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		Surname: handler_gen.OptString{
			Value: user.Surname,
			Set:   user.Surname != "",
		},
		LastName: handler_gen.OptString{
			Value: user.LastName,
			Set:   user.LastName != "",
		},
		Phone: user.PhoneNumber,
		Company: handler_gen.OptCompany{
			Set: true,
			Value: handler_gen.Company{
				Name: handler_gen.OptString{
					Value: *user.CompanyName,
					Set:   user.CompanyName != nil,
				},
				Address: handler_gen.OptString{
					Value: *user.CompanyAddress,
					Set:   user.CompanyAddress != nil,
				},
				Inn: handler_gen.OptString{
					Value: *user.CompanyInn,
					Set:   user.CompanyInn != nil,
				},
				Okpo: handler_gen.OptString{
					Value: *user.CompanyOkpo,
					Set:   user.CompanyOkpo != nil,
				},
				Position: handler_gen.OptString{
					Value: *user.PositionInCompany,
					Set:   user.PositionInCompany != nil,
				},
			},
		},
	}
}
