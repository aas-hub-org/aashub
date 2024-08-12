package database

import (
	b64 "encoding/base64"
	"os"

	mail "github.com/aas-hub-org/aashub/internal/mail"
)

type EmailVerificationRepository struct {
	VerificationRepository *VerificationRepository
}

func (e *EmailVerificationRepository) CreateVerification(email string) (string, error) {
	verificationCode, err := e.VerificationRepository.CreateVerification(email)
	if err != nil {
		return "", err
	}

	var server = os.Getenv("SERVER_ADDRESS")

	encodedMail := b64.RawURLEncoding.EncodeToString([]byte(email))
	encodedCode := b64.RawURLEncoding.EncodeToString([]byte(verificationCode))

	link := server + "/verify?email=" + encodedMail + "&code=" + encodedCode
	println(link)
	sendMailError := mail.SendEmail(email, "Verification Code", "<a href='"+link+"'>Click here to verify your email</a>")
	if sendMailError != nil {
		return "", sendMailError
	}

	return verificationCode, nil
}

func (e *EmailVerificationRepository) Verify(email string, verificationCode string) (string, error) {
	if errtype, err := e.VerificationRepository.Verify(email, verificationCode); err != nil {
		return errtype, err
	}
	return "", nil
}
