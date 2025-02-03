package core

type UserService interface {
	CreateUser(user User) error               // CreateUser creates a new user
	FindByEmail(email string) (User, error)   // FindByEmail returns a user by username
	CreateVerifly(verifly Verification) error // CreateVerifly creates a new verification
}

type userServiceImpl struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) CreateUser(user User) error {
	if err := s.repo.Save(user); err != nil {
		return err
	}
	return nil
}

func (s *userServiceImpl) CreateVerifly(verifly Verification) error {
	if err := s.repo.SaveVerifly(verifly); err != nil {
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
