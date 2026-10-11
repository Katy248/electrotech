package main

import (
	"electrotech/pkg/problem"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func main() {
	errInvalidData := errors.New("invalid data")
	p1 := problem.New("err:problem/1", 400, "something bad happened")

	fmt.Println(p1)

	p2 := p1.WithDetails("no rows").WithDetails("db err").WithDetails("occurs today")
	fmt.Println(p2)
	fmt.Println("is p1:", errors.Is(p2, p1))

	p3 := problem.NewFrom(errInvalidData, "err:problem/3", 400, "invalid data error")
	fmt.Println(p3)
	fmt.Println("is errInvalidData:", errors.Is(p3, errInvalidData))

	json.NewEncoder(os.Stdout).Encode(p3)

	baseP := problem.New("err:electrotech", http.StatusBadRequest, "base error")

	authP := baseP.Extend("/auth", http.StatusUnauthorized, "auth error")
	invalidTokenP := authP.Extend("/invalid-token", http.StatusUnauthorized, "invalid auth token provided")

	fmt.Println(invalidTokenP)
	json.NewEncoder(os.Stdout).Encode(invalidTokenP.WithDetails("token is empty"))
	fmt.Println("is auth error:", errors.Is(invalidTokenP.WithDetails("token is empty"), authP))
}
