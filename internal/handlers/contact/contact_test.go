package contact

import (
	"electrotech/internal/models"
	"os"
	"testing"
	"time"
)

const lorem = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

`

func TestBuildEmail(t *testing.T) {
	request := &models.Request{
		PersonName:   "John Doe",
		ID:           123,
		CreationDate: time.Now(),
		Message:      lorem,
	}

	mailAddr := "john.doe@example.com"
	phone := "8 800 555 35 35"
	request.Email = &mailAddr
	request.Phone = &phone

	email, err := buildEmail(request)
	if err != nil {
		t.Errorf("buildEmail failed: %v", err)
	}
	os.WriteFile("test.html", email, 0666)
}
