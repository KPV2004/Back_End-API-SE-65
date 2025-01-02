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
