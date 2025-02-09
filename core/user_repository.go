package core

type UserRepository interface {
	SaveUser(user User) error
	SaveAdmin(admin Admin) error
	LoginAdmin(username string, password string) error
	SaveVerify(verify Verification) error
	GetUserData(email string) (User, error)
	VerificationOTP(email string, OTP string) error
}
