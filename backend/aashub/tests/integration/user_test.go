//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	api "github.com/aas-hub-org/aashub/api/handler"
	db "github.com/aas-hub-org/aashub/internal/database"
	repositories "github.com/aas-hub-org/aashub/internal/database/repositories"
)

type testCase struct {
	name           string
	user           api.APIUser
	expectedStatus int
	expectedBody   string
}

func teardown(database *sql.DB) {
	// SQL statement to delete the test user, using ? as the placeholder
	query := "DELETE FROM Users WHERE username = ?"

	// Execute the query for the test username
	if _, err := database.Exec(query, "testuser"); err != nil {
		log.Fatalf("Failed to clean up test user: %v", err)
	}
}

func TestRegisterUser(t *testing.T) {
	// Initialize the database connection
	database, err := db.NewDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Ensure teardown is called no matter what happens in the test
	defer teardown(database)

	// Instantiate the repository
	verifyRepo := &repositories.VerificationRepository{DB: database}
	userRepo := &repositories.UserRepository{DB: database, VerificationRepository: verifyRepo}

	// Instantiate the handler struct with the repository
	userHandler := &api.UserHandler{Repo: userRepo}
	verifyHandler := &api.VerificationHandler{VerificationRepository: verifyRepo}

	// Create a new ServeMux.
	mux := http.NewServeMux()

	// Register handlers for different endpoints.
	mux.HandleFunc("/register", userHandler.RegisterUser)
	mux.HandleFunc("/verify", verifyHandler.VerifyUser)

	// Create a new HTTP test server
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Define test cases
	tests := []struct {
		name           string
		user           api.APIUser
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful Registration",
			user:           api.APIUser{Username: "testuser", Email: "test@example.com", Password: "password123"},
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
		},
		{
			name:           "Missing Fields",
			user:           api.APIUser{Username: "", Email: "test@example.com", Password: "password123"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing required field(s)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal the user object to JSON
			body, err := json.Marshal(tc.user)
			if err != nil {
				t.Fatalf("Failed to marshal user: %v", err)
			}

			// Create a new POST request
			req, err := http.NewRequest("POST", ts.URL+"/register", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Perform the request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to perform request: %v", err)
			}
			defer resp.Body.Close()

			// Read the response body
			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			// Assert the status code
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			// Assert the response body if expected
			if tc.expectedBody != "" && !strings.Contains(string(responseBody), tc.expectedBody) {
				t.Errorf("Expected response body to contain %q, got %q", tc.expectedBody, string(responseBody))
			}

			// Call the verify user function
			verifyUser(t, tc, ts, database)
		})
	}
}

func verifyUser(t *testing.T, tc testCase, ts *httptest.Server, database *sql.DB) {
	// Assuming successful registration, proceed to verification
	if tc.expectedStatus == http.StatusCreated {
		// get the verification code from the database
		var verificationCode string
		query := "SELECT verification_code FROM Verifications WHERE email = ?"
		if err := database.QueryRow(query, tc.user.Email).Scan(&verificationCode); err != nil {
			t.Fatalf("Failed to get verification code: %v", err)
		}

		// Base64 URL encode the email and verification code
		emailEncoded := base64.RawURLEncoding.EncodeToString([]byte(tc.user.Email))
		codeEncoded := base64.RawURLEncoding.EncodeToString([]byte(verificationCode))

		// Construct the verification URL with query parameters
		verifyURL := fmt.Sprintf("%s/verify?email=%s&code=%s", ts.URL, url.QueryEscape(emailEncoded), url.QueryEscape(codeEncoded))

		log.Printf("Verification URL: %s", verifyURL)

		// Create a new GET request to the verification endpoint
		verifyReq, err := http.NewRequest("GET", verifyURL, nil)
		if err != nil {
			t.Fatalf("Failed to create verification request: %v", err)
		}

		// Perform the verification request
		verifyResp, err := http.DefaultClient.Do(verifyReq)
		if err != nil {
			t.Fatalf("Failed to perform verification request: %v", err)
		}
		defer verifyResp.Body.Close()

		// Read the verification response body
		verifyResponseBody, err := io.ReadAll(verifyResp.Body)
		if err != nil {
			t.Fatalf("Failed to read verification response body: %v", err)
		}

		log.Printf("Verification response: %s", verifyResponseBody)

		// Assert the verification status code and response body
		if verifyResp.StatusCode != http.StatusOK {
			t.Errorf("Expected verification status code %d, got %d", http.StatusOK, verifyResp.StatusCode)
		}
		if !strings.Contains(string(verifyResponseBody), "User verified successfully") {
			t.Errorf("Expected verification response body to contain %q, got %q", "User verified successfully", string(verifyResponseBody))
		}
	}
}

func TestLoginUser(t *testing.T) {
	// Initialize the database connection
	database, err := db.NewDB()
	if err != nil {
		t.Fatalf("Could not connect to the database: %v", err)
	}

	// Instantiate the repository
	verifyRepo := &repositories.VerificationRepository{DB: database}
	userRepo := &repositories.UserRepository{DB: database, VerificationRepository: verifyRepo}

	// Instantiate the handler struct with the repository
	userHandler := &api.UserHandler{Repo: userRepo}

	// Create a buffer to write our multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the form fields
	_ = writer.WriteField("identifier", "test")
	_ = writer.WriteField("password", "test")
	// Close the writer to finalize the multipart body
	writer.Close()

	// Create a request to pass to our handler
	req, err := http.NewRequest("POST", "/login", body)
	if err != nil {
		t.Fatal(err)
	}
	// Set the content type to multipart/form-data with the boundary parameter
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// We create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandler.LoginUser)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}
