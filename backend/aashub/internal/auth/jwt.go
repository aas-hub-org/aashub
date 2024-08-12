package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Define a struct to hold your payload. You can add more fields as needed.
type CustomClaims struct {
	Payload string `json:"payload"`
	jwt.StandardClaims
}

// Function to generate a JWT token with a string payload
func GenerateJWT(payload string, secretKey string) (string, error) {
	// Set expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the claims with the payload and standard claims
	claims := CustomClaims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
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

func IsTokenValid(token string, secretkey string) (bool, error) {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretkey), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
