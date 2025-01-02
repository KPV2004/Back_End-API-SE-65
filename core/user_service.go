package core

type UserService interface {
	CreateUser(user User) error
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
