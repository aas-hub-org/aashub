package interfaces

type VerificationRepositoryInterface interface {
	CreateVerification(email string) (string, error)
	Verify(email string, verificationCode string) (string, error)
	IsVerified(email string) (bool, error)
}
