//go:build unit
// +build unit

package unit_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	b64 "encoding/base64"

	api "github.com/aas-hub-org/aashub/api/handler"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (m *MockRepository) RegisterUser(username, email, password string) error {
	return nil
}

func (repo *MockRepository) LoginUser(username string, password string) (string, error) {
	return "", nil
}

var (
	// VerifyFunc is a package-level variable that can be overridden in tests.
	VerifyFunc func(email, code string) (string, error)
)

func (m *MockRepository) CreateVerification(email string) (string, error) {
	return "", nil
}

func (m *MockRepository) Verify(email, code string) (string, error) {
	if VerifyFunc != nil {
		return VerifyFunc(email, code)
	}
	// Default behavior if VerifyFunc is not set
	return "", nil
}

func (m *MockRepository) IsVerified(email string) (bool, error) {
	return true, nil
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := &MockRepository{}
	handler := api.UserHandler{Repo: mockRepo}

	user := api.APIUser{
		Username: "testUser",
		Email:    "test@example.com",
		Password: "password123",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/users/register", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/register", handler.RegisterUser)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status code 201")
}

func TestRegisterUser_InvalidRequest(t *testing.T) {
	mockRepo := &MockRepository{}
	handler := api.UserHandler{Repo: mockRepo}

	user := api.APIUser{
		Username: "testUser",
		Email:    "test@example.com",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/users/register", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/register", handler.RegisterUser)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code 400")
}

func TestVerifyUser_Success(t *testing.T) {
	// Set up the mock behavior
	originalVerifyFunc := VerifyFunc
	VerifyFunc = func(email, code string) (string, error) {
		return "", nil // Simulate successful verification
	}
	defer func() { VerifyFunc = originalVerifyFunc }() // Reset VerifyFunc after the test

	mockRepo := &MockRepository{}
	handler := api.VerificationHandler{VerificationRepository: mockRepo}

	// Simulate the request
	email := "test@example.com"
	code := "verificationCode"
	emailEncoded := url.QueryEscape(b64.RawURLEncoding.EncodeToString([]byte(email)))
	codeEncoded := url.QueryEscape(b64.RawURLEncoding.EncodeToString([]byte(code)))

	req, err := http.NewRequest("GET", "/verify?email="+emailEncoded+"&code="+codeEncoded, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/verify", handler.VerifyUser)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "User verified successfully"
	if rr.Body.String() != expected {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestVerifyUser_Failure(t *testing.T) {
	// Set up the mock behavior for failure
	originalVerifyFunc := VerifyFunc
	VerifyFunc = func(email, code string) (string, error) {
		return "user", errors.New("Invalid email or code") // Simulate verification failure due to invalid input
	}
	defer func() { VerifyFunc = originalVerifyFunc }() // Reset VerifyFunc after the test

	mockRepo := &MockRepository{}
	handler := api.VerificationHandler{VerificationRepository: mockRepo}

	// Simulate the request with invalid email and code
	email := "invalidEmail"
	code := "invalidCode"
	emailEncoded := url.QueryEscape(b64.RawURLEncoding.EncodeToString([]byte(email)))
	codeEncoded := url.QueryEscape(b64.RawURLEncoding.EncodeToString([]byte(code)))

	req, err := http.NewRequest("GET", "/verify?email="+emailEncoded+"&code="+codeEncoded, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/verify", handler.VerifyUser)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect for failure.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check the response body is what we expect for failure.
	actual := strings.TrimSpace(rr.Body.String())
	expected := strings.TrimSpace("Invalid email or code")

	if actual != expected {
		t.Fatalf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}
