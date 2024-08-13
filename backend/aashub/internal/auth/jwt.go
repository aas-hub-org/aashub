package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Define a struct to hold your payload. You can add more fields as needed.
type CustomClaims struct {
	Payload string `json:"payload"`
	jwt.RegisteredClaims
}

// Function to generate a JWT token with a string payload
func GenerateJWT(payload string, secretKey string) (string, error) {
	// Set expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the claims with the payload and registered claims
	claims := CustomClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func IsTokenValid(tokenString string, secretKey string) (bool, error) {
	_, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
