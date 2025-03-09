package core

type UserService interface {
	CreateUser(user User) error               // CreateUser creates a new user
	CreateAdmin(admin Admin) error            // CreateAdmin creates a new admin
	FindByEmail(email string) (User, error)   // FindByEmail returns a user by username
	CreateVerifly(verifly Verification) error // CreateVerifly creates a new verification
	VerificationOTP(verifly Verification) error
	LoginAdmin(admin Admin) error
	UpdateUser(user User, email string) error
	UpdateUserPlanByEmail(email string, newPlanID string) error
	CreatePlan(planData Plan) error // CreatePlan creates a new plan
	AddTripLocation(planID string, newLocation TripLocation, index int) error
  UpdatePlan(planData Plan, planID string) error
	GetTripLocationByPlanID(planID string) ([]TripLocation, error)
	GetPlanByID(planID string) (Plan, error)
	DeletePlanByID(planID string) error
	DeleteUserPlanByEmail(email, planID string) error
	GetVisiblePlans() ([]Plan, error)
}

type userServiceImpl struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) CreateUser(user User) error {
	if err := s.repo.SaveUser(user); err != nil {
		return err
	}
	return nil
}

func (s *userServiceImpl) CreateAdmin(admin Admin) error {
	if err := s.repo.SaveAdmin(admin); err != nil {
		return err
	}
	return nil
}

func (s *userServiceImpl) CreatePlan(planData Plan) error {
	if err := s.repo.CreatePlan(planData); err != nil {
		return err
	}
	return nil
}

func (s *userServiceImpl) CreateVerifly(verifly Verification) error {
	if err := s.repo.SaveVerify(verifly); err != nil {
		return err
	}
	return nil
}

func (s *userServiceImpl) FindByEmail(email string) (User, error) {

	user, err := s.repo.GetUserData(email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *userServiceImpl) VerificationOTP(verifly Verification) error {
	if err := s.repo.VerificationOTP(verifly.Email, verifly.Otp); err != nil {
		return err
	}
	return nil
}

func (s *userServiceImpl) LoginAdmin(admin Admin) error {
	if err := s.repo.LoginAdmin(admin.Username, admin.Password); err != nil {
		return err
	}
	return nil
}
func (s *userServiceImpl) UpdateUser(user User, email string) error {
	if err := s.repo.UpdateUser(user, email); err != nil {
		return err
	}
	return nil
}
func (s *userServiceImpl) UpdatePlan(planData Plan, planID string) error{
  if err:= s.repo.UpdatePlane(planData, planID); err != nil {
    return err
  }
  return nil
}
func (s *userServiceImpl) UpdateUserPlanByEmail(email string, newPlanID string) error {
	if err := s.repo.UpdateUserPlanByEmail(email, newPlanID); err != nil {
		return err
	}
	return nil
}
func (s *userServiceImpl) AddTripLocation(planID string, newLocation TripLocation, index int) error {
	if err := s.repo.AddTripLocation(planID, newLocation, index); err != nil {
		return err
	}
	return nil
}
func (s *userServiceImpl) GetTripLocationByPlanID(planID string) ([]TripLocation, error) {
	locations, err := s.repo.GetTripLocationByPlanID(planID)
	if err != nil {
		return nil, err
	}
	return locations, nil
}
func (s *userServiceImpl) GetPlanByID(planID string) (Plan, error) {
	plan, err := s.repo.GetPlanByID(planID)
	if err != nil {
		return plan, err
	}
	return plan, nil
}
func (s *userServiceImpl) DeletePlanByID(planID string) error {
	if err := s.repo.DeletePlanByID(planID); err != nil {
		return err
	}
	return nil
}
func (s *userServiceImpl) DeleteUserPlanByEmail(email, planID string) error {
	if err := s.repo.DeleteUserPlanByEmail(email, planID); err != nil {
		return err
	}
	return nil
}
func (s *userServiceImpl) GetVisiblePlans() ([]Plan, error) {
	plans, err := s.repo.GetVisiblePlans()
	if err != nil {
		return nil, err
	}
	return plans, nil
}
