package models

import "time"

type Request struct {
	ID           int       `json:"id"`
	CreationDate time.Time `json:"creationDate"`
	PersonName   string    `json:"personName"`
	Email        *string   `json:"email"`
	Phone        *string   `json:"phone"`
	Message      string    `json:"message"`
}
