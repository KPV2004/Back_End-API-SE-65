package core

type UserRepository interface {
	Save(user User) error
	GetUserData(email string) (User, error)
}
