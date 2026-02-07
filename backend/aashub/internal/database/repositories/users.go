package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"

	auth "github.com/aas-hub-org/aashub/internal/auth"
	interfaces "github.com/aas-hub-org/aashub/internal/interfaces"
	utils "github.com/aas-hub-org/aashub/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserRepoNotFound = errors.New("identifier or password wrong")
var ErrUserRepoNotVerified = errors.New("user not verified")
var ErrUserRepoEmailUsernameExists = errors.New("email or username already exists")

type UserRepository struct {
	DB                     *sql.DB
	VerificationRepository interfaces.VerificationRepositoryInterface
}

type User struct {
	ID       string
	Username string
	Email    string
	Password string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (repo *UserRepository) RegisterUser(username string, email string, password string) error {
	userid := uuid.New().String()
	hashedpassword, err := HashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	_, err = repo.DB.Exec("INSERT INTO Users (id, username, email, password_hash) VALUES (?, ?, ?, ?)", userid, username, email, hashedpassword)

	if err != nil {
		// Check if error includes "Duplicate"
		if strings.Contains(err.Error(), "Duplicate") {
			return ErrUserRepoEmailUsernameExists
		}
		return err
	}
	if IsVerificationEnabled() {
		_, err = repo.VerificationRepository.CreateVerification(email)
	}
	if err != nil {
		log.Printf("Error inserting verification: %v", err)
		return err
	}

	return nil
}

func (repo *UserRepository) LoginUser(identifier string, password string) (string, error) {
	// Changed the error message to 'identifier' to generalize username/email
	var user User

	// Adjust the SQL query to check both the username and email fields
	err := repo.DB.QueryRow("SELECT * FROM Users WHERE username = ? OR email = ?", identifier, identifier).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return "", ErrUserRepoNotFound
	}
	if IsVerificationEnabled() {
		isVerified, err := repo.VerificationRepository.IsVerified(user.Email)
		if err != nil {
			log.Printf("Error checking verification: %v", err)
			return "", err
		}
		if !isVerified {
			return "", ErrUserRepoNotVerified
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrUserRepoNotFound
	}

	secret, fileReadError := utils.ReadFile("/workspace/backend/aashub/privatekey.txt")

	if fileReadError != nil {
		log.Printf("Error reading file: %v", fileReadError)
		return "", fileReadError
	}

	jwt, err := auth.GenerateJWT(user.ID, secret)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return "", err
	}

	return jwt, nil
}

func IsVerificationEnabled() bool {
	enabled, found := os.LookupEnv("VERIFICATION_ENABLED")
	if !found {
		return false
	}
	return enabled == "true"
}
