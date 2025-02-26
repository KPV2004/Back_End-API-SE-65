package core

type UserRepository interface {
	SaveUser(user User) error
	SaveAdmin(admin Admin) error
	LoginAdmin(username string, password string) error
	SaveVerify(verify Verification) error
	GetUserData(email string) (User, error)
	VerificationOTP(email string, OTP string) error
	UpdateUser(user User, email string) error
	UpdateUserPlanByEmail(email string, newPlanID string) error
	CreatePlan(userPlan Plan) error
	AddTripLocation(planID string, newPlaceID string) error
	GetTripLocationByPlanID(planID string) ([]string, error)
	GetPlanByID(planID string) (Plan, error)
}
