package models

import (
	"electrotech"
	"time"

	"charm.land/log/v2"
)

type UserQuestion struct {
	ID           int       `json:"id"`
	CreationDate time.Time `json:"creationDate"`
	PersonName   string    `json:"personName"`
	Email        *string   `json:"email"`
	Phone        *string   `json:"phone"`
	Message      string    `json:"message"`
	ClientIP     string    `json:"clientIp"`
}

func NewUserQuestion(name, email, phone, message, ip string) *UserQuestion {
	if phone != "" {
		formatted, err := electrotech.FormatPhoneNumber(phone)
		if err != nil {
			log.Warn("Failed format phone number")
		} else {
			phone = formatted
		}
	}

	return &UserQuestion{ //nolint:exhaustruct_v5
		CreationDate: time.Now(),

		PersonName: name,
		Email:      strToPointer(email),
		Phone:      strToPointer(phone),
		Message:    message,
		ClientIP:   ip,
	}
}

func strToPointer(str string) *string {
	if str == "" {
		return nil
	}

	return &str
}
