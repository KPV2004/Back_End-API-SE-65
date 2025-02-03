package core

type UserService interface {
	CreateUser(user User) error                   // CreateUser creates a new user
	FindByUsername(username string) (User, error) // FindByUsername returns a user by username
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

func (s *userServiceImpl) FindByUsername(username string) (User, error) {
	user, err := s.repo.GetUserData(username)
	if err != nil {
		return user, err
	}
	return user, nil
}
