package adapters

import (
	"go-server/core"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) core.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Save(user core.User) error {
	if result := r.db.Create(&user); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormUserRepository) FindByUsername(username string) (core.User, error) {
	var user core.User
	if result := r.db.Where("username = ?", username).First(&user); result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
