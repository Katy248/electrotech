package main

import (
	"context"
	"database/sql"
	"electrotech/internal/handlers/auth"
	"electrotech/internal/repository/users"
	"electrotech/storage"
	"flag"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

const DefaultPassword = "P@$$w0rd"

var (
	password = flag.String("password", DefaultPassword, "password")
	userRepo *users.Queries
	ctx      = context.Background()
)

func init() {
	godotenv.Load()
	flag.Parse()
}

func main() {

	db, err := storage.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	userRepo = users.New(db)
	if err := exec(); err != nil {
		log.Fatal("Failed to execute", "error", err)
	}
}

func exec() error {
	passwordHash, err := auth.GeneratePasswordHash(*password)
	if err != nil {
		return err
	}

	userRepo.InsertNew(ctx, users.InsertNewParams{
		FirstName:    "Снейк",
		LastName:     "БигБоссович",
		Surname:      "Веном",
		Email:        "test1@mail.ru",
		PasswordHash: passwordHash,
		PhoneNumber:  "+79999999999",
	})
	fmt.Printf("Created new user:\ntest1@mail.ru\nPassword: %s\n", *password)
	userRepo.InsertNew(ctx, users.InsertNewParams{
		FirstName:    "Снейк",
		LastName:     "БигБоссович",
		Surname:      "Веном",
		Email:        "test2@mail.ru",
		PasswordHash: passwordHash,
		PhoneNumber:  "+79999999998",
	})
	userRepo.UpdateCompanyData(ctx, users.UpdateCompanyDataParams{
		CompanyName:       sql.NullString{String: "Test Company", Valid: true},
		Email:             "test2@mail.ru",
		CompanyAddress:    sql.NullString{String: "Test Address", Valid: true},
		CompanyInn:        sql.NullString{String: "Test INN", Valid: true},
		CompanyOkpo:       sql.NullString{String: "Test OKPO", Valid: true},
		PositionInCompany: sql.NullString{String: "Test Position", Valid: true},
	})
	fmt.Printf("Created new user with company data:\ntest2@mail.ru\nPassword: %s\n", *password)
	return nil

}
