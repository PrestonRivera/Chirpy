package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)


func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret"

	t.Run("valid token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("Error creating token: %v", err)
		}

		gotUID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("Error validating token: %v", err)
		}

		if gotUID != userID {
			t.Errorf("got user ID %v, want %v", gotUID, userID)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, -time.Hour)
		if err != nil {
			t.Fatalf("Error creating token: %v", err)
		}

		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Fatal("Expected error for expired token, got nil")
		}
	})

	t.Run("invalid secret", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("Error creating token: %v", err)
		}

		_, err = ValidateJWT(token, "invalid")
		if err == nil {
			t.Fatal("Expected error for invalid secret, got nil")
		}
	})
}


func TestGetBearerToken(t *testing.T) {
	const myToken string = "mytoken123"

	t.Run("valid Token", func(t *testing.T) {
		headers := http.Header{
			"Authorization": []string{"Bearer " + myToken},
		}

		token, err := GetBearerToken(headers)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token != myToken {
			t.Errorf("expected: %v, got: %v", myToken, token)
		}
	})

	t.Run("missing header", func(t *testing.T) {
		headers := http.Header{}
		
		_,  err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("malformed header", func(t *testing.T) {
		headers := http.Header{
			"Authorization": []string{"TheBearer " + myToken},
		}

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("malformed header 2", func(t *testing.T) {
		headers := http.Header{
			"Authorization": []string{" Bearer " + myToken},
		}

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}