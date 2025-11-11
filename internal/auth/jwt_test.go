package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userId := uuid.New()
	secret := "supersecret"
	token, err := MakeJWT(userId, secret, time.Duration(time.Hour))

	if err != nil {
		t.Errorf("MakeJWT returned an error: %v", err)
	}

	if token == "" {
		t.Errorf("MakeJWT returned an empty token")
	}

	returnedUserId, err := ValidateJWT(token, secret)

	if err != nil {
		t.Errorf("ValidateJWT returned an error: %v", err)
	}

	if returnedUserId != userId {
		t.Errorf("ValidateJWT returned wrong userId: got %v, want %v", returnedUserId, userId)
	}
}
