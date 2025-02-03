package core

type UserRepository interface {
	Save(user User) error
	SaveVerifly(verifly Verification) error
	GetUserData(email string) (User, error)
	VerificationOTP(email string, OTP string) error
}
