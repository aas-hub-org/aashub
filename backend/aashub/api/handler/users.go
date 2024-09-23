package api

import (
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	repositories "github.com/aas-hub-org/aashub/internal/database/repositories"
	interfaces "github.com/aas-hub-org/aashub/internal/interfaces"
)

type APIUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserHandler struct {
	Repo interfaces.UserRepositoryInterface
}

type VerificationHandler struct {
	VerificationRepository interfaces.VerificationRepositoryInterface
}

// RegisterUser registers a new user in the system.
// @Summary Register a new user
// @Description Registers a new user with the provided username, email, and password.
// @Tags users
// @Accept json
// @Produce json
// @Param user body APIUser true "User to register"
// @Success 201 {string} string "Successfully registered the user"
// @Failure 400 {string} string "Invalid request parameters"
// @Failure 500 {string} string "Internal server error"
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user APIUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if any of the required fields are empty
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Missing required field(s)", http.StatusBadRequest)
		return
	}

	err = h.Repo.RegisterUser(user.Username, user.Email, user.Password)
	if err != nil {
		switch err {
		case repositories.ErrUserRepoEmailUsernameExists:
			http.Error(w, "Email or username already exists", http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// VerifyUser godoc
// @Summary Verify user
// @Description Verifies a user using base64 URL encoded email and verification code.
// @Tags verification
// @Accept  json
// @Produce  json
// @Param   email   query    string     true  "Base64 URL Encoded Email"
// @Param   code    query    string     true  "Base64 URL Encoded Verification Code"
// @Success 200  {string}  string  "User verified successfully"
// @Failure 400  {string}  string  "Invalid email or code"
// @Failure 500  {string}  string  "Verification failed"
// @Router /verify [get]
func (h *VerificationHandler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters
	email_byte, mail_decode_err := b64.RawURLEncoding.DecodeString(r.URL.Query().Get("email"))
	code_byte, code_decode_err := b64.RawURLEncoding.DecodeString(r.URL.Query().Get("code"))

	if mail_decode_err != nil || code_decode_err != nil {
		http.Error(w, "Invalid email or code", http.StatusBadRequest)
		return
	}

	email := string(email_byte)
	code := string(code_byte)

	error_type, err := h.VerificationRepository.Verify(email, code)
	if err != nil {
		if error_type == "system" {
			http.Error(w, "Verification failed", http.StatusInternalServerError)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Write success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User verified successfully"))
}

// LoginUser logs in a user, sets a cookie with a JWT token, and returns the token in the response
// @Summary User login and set cookie
// @Description Logs in a user by identifier (username or email) and password, sets a cookie with a JWT token if successful, and returns the JWT token in the response.
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param identifier formData string true "Username or Email"
// @Param password formData string true "Password"
// @Success 204 "Successfully logged in"
// @Failure 400 {object} map[string]string "Missing required field(s) or bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/login [post]
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	// Check if any of the required fields are empty
	if identifier == "" || password == "" {
		http.Error(w, "Missing required field(s)", http.StatusBadRequest)
		return
	}

	token, err := h.Repo.LoginUser(identifier, password)
	if err != nil {
		switch err {
		case repositories.ErrUserRepoNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
		case repositories.ErrUserRepoNotVerified:
			http.Error(w, "User not verified", http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// Create a cookie
	expiration := time.Now().Add(24 * time.Hour) // Set expiration to 24 hours from now
	cookie := http.Cookie{
		Name:     "token",    // Name of the cookie
		Value:    token,      // Token value
		Expires:  expiration, // Expiration time
		HttpOnly: true,       // Make the cookie HTTP-only (not accessible via JavaScript)
		Path:     "/",        // Cookie path
		// Secure:   true,     // Uncomment this if you are serving your site over HTTPS
	}

	// Set the cookie in the response header
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusNoContent)
}
