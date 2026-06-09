package mail

import (
	"crypto/tls"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func SendEmail(to, subject, body string) error {
	log.Printf("Sending email to %s", to)
	// Load .env
	env_err := godotenv.Load("/workspace/backend/aashub/.env")
	if env_err != nil {
		log.Printf("Error loading .env file")
	}

	// Get from .env
	from := os.Getenv("MAIL_ADDRESS")
	pass := os.Getenv("MAIL_PASSWORD")

	// SMTP server configuration.
	smtpHost := os.Getenv("MAIL_SMTP")
	smtpPort := os.Getenv("SMTP_PORT")

	// Message.
	message := []byte("From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		body)

	// TLS configuration
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// Connect to the SMTP Server
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		log.Printf("Error connecting to SMTP server: %v", err)
		return err
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Printf("Error creating SMTP client: %v", err)
		return err
	}

	// Authentication
	auth := smtp.PlainAuth("", from, pass, smtpHost)
	if err = client.Auth(auth); err != nil {
		log.Printf("Error authenticating: %v", err)
		return err
	}

	// To && From
	if err = client.Mail(from); err != nil {
		log.Printf("Error setting sender: %v", err)
		return err
	}
	if err = client.Rcpt(to); err != nil {
		log.Printf("Error setting recipient: %v", err)
		return err
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Printf("Error getting SMTP data writer: %v", err)
		return err
	}

	_, err = w.Write(message)
	if err != nil {
		log.Printf("Error writing message: %v", err)
		return err
	}

	err = w.Close()
	if err != nil {
		log.Printf("Error closing SMTP data writer: %v", err)
		return err
	}

	client.Quit()

	log.Printf("Email Sent Successfully!")
	return nil
}
