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
	UpdatePlane(plan Plan, planID string) error
	CreatePlan(userPlan Plan) error
	AddTripLocation(planID string, newLocation TripLocation, index int) error
	GetTripLocationByPlanID(planID string) ([]TripLocation, error)
	GetPlanByID(planID string) (Plan, error)
	DeletePlanByID(planID string) error
	DeleteUserPlanByEmail(email, planID string) error
	GetVisiblePlans() ([]Plan, error)
	DeleteTripLocation(planID, targetPlaceID string) error
	UpdateAuthorImg(planID, newImg string) error
}
