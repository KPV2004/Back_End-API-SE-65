package core

type UserService interface {
	CreateUser(user User) error               // CreateUser creates a new user
	CreateAdmin(admin Admin) error            // CreateAdmin creates a new admin
	FindByEmail(email string) (User, error)   // FindByEmail returns a user by username
	CreateVerifly(verifly Verification) error // CreateVerifly creates a new verification
	CreatePlan(planData Plan) error           // CreatePlan creates a new plan
	VerificationOTP(verifly Verification) error
	LoginAdmin(admin Admin) error
	UpdateUser(user User, email string) error
	UpdateUserPlanByEmail(email string, newPlanID string) error
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
func (s *userServiceImpl) UpdateUserPlanByEmail(email string, newPlanID string) error {
	if err := s.repo.UpdateUserPlanByEmail(email, newPlanID); err != nil {
		return err
	}
	return nil
}
