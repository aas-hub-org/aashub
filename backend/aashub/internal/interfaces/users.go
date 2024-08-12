package interfaces

type UserRepositoryInterface interface {
	RegisterUser(username string, email string, password string) error
	LoginUser(username string, password string) (string, error)
}
