package user

import (
	"context"
	"electrotech"
	"electrotech/internal/config"
	"electrotech/internal/models"
	"electrotech/pkg/problem"
	"fmt"
	"net/http"
	"strings"

	"charm.land/log/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	ByID(ctx context.Context, id int64) (*models.User, error)
	ByEmail(ctx context.Context, email string) (*models.User, error)
	InsertNew(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	UpdatePartial(ctx context.Context, userID int64, updates map[string]any) error
}

type Service struct {
	logger *log.Logger
	repo   UserRepository
	config *config.AuthConfig
}

func New(
	logger *log.Logger,
	repo UserRepository,
	config *config.AuthConfig,
) (*Service, error) {
	service := &Service{
		logger: logger,
		repo:   repo,
		config: config,
	}
	if _, err := service.getSecretKey(); err != nil {
		return nil, fmt.Errorf("bad secret key specified: %w", err)
	}

	return service, nil
}

var (
	errAuth = problem.New(
		"err:electrotech/auth", http.StatusUnauthorized, "auth error")
	errInternal = problem.New(
		"err:electrotech/internal",
		http.StatusInternalServerError,
		"internal server error",
	)
	errInvalidCredentials = errAuth.Extend(
		"/invalid-credentials",
		http.StatusUnauthorized,
		"invalid credentials",
	)
	errEmailTaken = errAuth.Extend(
		"/email-taken",
		http.StatusConflict,
		"email already taken",
	)
	errCreateAuthResponse = errInternal.Extend(
		"/",
		http.StatusInternalServerError,
		"internal server error",
	).WithDetails("creating auth response")
	errHashPassword = errInternal.Extend(
		"/",
		http.StatusInternalServerError,
		"internal server error",
	).WithDetails("hash user password")
	errFormatPhone = errAuth.Extend(
		"/phone",
		http.StatusBadRequest,
		"phone number format invalid",
	)
	errInsertNewUser = errInternal.Extend("/", http.StatusInternalServerError, "internal server error").
				WithDetails("insert user")
	errInvalidToken = errAuth.Extend("/invalid-token", http.StatusUnauthorized, "invalid refresh token")
)

func (s *Service) Login(ctx context.Context, email, password string) (*models.AuthResponse, error) {
	user, err := s.repo.ByEmail(ctx, email)
	if err != nil || user.Email == "" {
		s.logger.Error("Failed get user", "user", email, "error", err)

		return nil, errInvalidCredentials.WithError(err)
	}

	if !user.CheckPassword(password) {
		s.logger.Error("Password doesn't match", "user", email)

		return nil, errInvalidCredentials
	}

	r, err := s.newAuthResponse(user)
	if err != nil {
		return nil, errCreateAuthResponse.WithError(err)
	}

	return r, nil
}

func (s *Service) Register(ctx context.Context, req *models.RegisterRequest) error {
	phone, err := electrotech.FormatPhoneNumber(req.PhoneNumber)
	if err != nil {
		s.logger.Error("Error formatting phone number (probably number invalid)", "phone", req.PhoneNumber, "error", err)

		return errFormatPhone.WithError(err).WithDetails(err.Error())
	}

	req.Email = strings.ToLower(req.Email)

	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.repo.ByEmail(ctx, req.Email)
	if err == nil && existingUser.Email != "" {
		s.logger.Error("Attempt to create user with email already taken", "email", req.Email)

		return errEmailTaken.WithError(err)
	}

	user := &models.User{ //nolint:exhaustruct_v5
		Email:       req.Email,
		FirstName:   req.FirstName,
		Surname:     req.Surname,
		LastName:    req.LastName,
		PhoneNumber: phone,
	}
	if err := user.SetPassword(req.Password); err != nil {
		s.logger.Error("Failed set (hash) user password")

		return errHashPassword.WithError(err)
	}

	// Создаем нового пользователя
	err = s.repo.InsertNew(ctx, user)
	if err != nil {
		s.logger.Errorf("Error creating user: %v", err)

		return errInsertNewUser.WithError(err)
	}

	return nil
}

func (s *Service) Refresh(ctx context.Context, token string) (*models.AuthResponse, error) {
	claimsUser, err := s.ValidateToken(token)
	if err != nil {
		s.logger.Error("Invalid refresh token", "error", err)

		return nil, errInvalidToken.WithError(err)
	}

	user, err := s.repo.ByID(ctx, claimsUser.Id)
	if err != nil {
		s.logger.Error("Error getting user by id", "userID", claimsUser.Id, "error", err)

		return nil, errInvalidToken.WithDetails("no such user").WithError(err)
	}

	response, err := s.newAuthResponse(user)
	if err != nil {
		s.logger.Error("Error generating auth response", "error", err)

		return nil, errCreateAuthResponse
	}

	return response, nil
}

func (s *Service) UpdateEmail(ctx context.Context, id int64, email string) error {
	user, err := s.repo.ByID(ctx, id)
	if err != nil {
		return fmt.Errorf("select user by id: %w", err)
	}

	user.Email = email

	err = s.repo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("save changes: %w", err)
	}

	return nil
}

func (s *Service) UpdatePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	user, err := s.repo.ByID(ctx, id)
	if err != nil {
		return fmt.Errorf("select user by id: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return fmt.Errorf("passwords don't match: %w", err)
	}

	err = user.SetPassword(newPassword)
	if err != nil {
		return fmt.Errorf("set password: %w", err)
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("save changes: %w", err)
	}

	return nil
}

func (s *Service) UpdatePhone(ctx context.Context, id int64, phone string) error {
	user, err := s.repo.ByID(ctx, id)
	if err != nil || user.Email == "" {
		return fmt.Errorf("select user by id: %w", err)
	}

	phone, err = electrotech.FormatPhoneNumber(phone)
	if err != nil {
		return fmt.Errorf("format phone: %w", err)
	}

	user.PhoneNumber = phone

	err = s.repo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("save changes: %w", err)
	}

	return nil
}

func (s *Service) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.repo.ByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("select user by id: %w", err)
	}

	return user, nil
}

// Update updates user fields.
//
// TODO: Check empty string setting.
func (s *Service) Update(ctx context.Context, id int64, update *models.UserUpdate) error {
	if !update.HasChanges() {
		return nil // nothing to do — no DB call
	}

	err := s.repo.UpdatePartial(ctx, id, update.Values())
	if err != nil {
		return fmt.Errorf("save changes: %w", err)
	}

	return nil
}
