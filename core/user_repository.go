package core

type UserRepository interface {
	Save(user User) error
	GetUserData(username string) (User, error)
}
