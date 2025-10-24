//go:build unit
// +build unit

package unit_test

import (
	"strings"
	"testing"

	"github.com/aas-hub-org/aashub/internal/auth"
)

func TestGenerateJWTAndValidate(t *testing.T) {
	// Define a payload and a secret key for testing
	payload := "testPayload"
	secretKey := "testSecretKey"

	// Generate a JWT token
	tokenString, err := auth.GenerateJWT(payload, secretKey)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Check if the token is valid
	isValid, err := auth.IsTokenValid(tokenString, secretKey)
	if err != nil {
		t.Fatalf("Error validating token: %v", err)
	}

	if !isValid {
		t.Fatalf("The token was expected to be valid")
	}
}

func TestGenerateJWTAndValidateWithManipulatedPayload(t *testing.T) {
	expectedPayload := "testPayload"
	manipulatedPayload := "manipulated"
	secretKey := "testSecret"

	// Generate a JWT token
	tokenString, err := auth.GenerateJWT(expectedPayload, secretKey)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Generate a second JWT token
	secondTokenString, err := auth.GenerateJWT(manipulatedPayload, secretKey)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Check if the token is valid
	isValid, err := auth.IsTokenValid(tokenString, secretKey)
	if err != nil {
		t.Fatalf("Error validating token: %v", err)
	}

	if !isValid {
		t.Fatalf("The token was expected to be valid")
	}

	// Split token at .
	tokenParts := strings.Split(tokenString, ".")
	secondTokenParts := strings.Split(secondTokenString, ".")

	newToken := tokenParts[0] + "." + secondTokenParts[1] + "." + tokenParts[2]

	// Check if the token is valid
	_, err = auth.IsTokenValid(newToken, secretKey)
	if err == nil {
		t.Fatalf("Expected an error when validating token")
	}
}
