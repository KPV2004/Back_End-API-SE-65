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

func (r *GormUserRepository) SaveUser(user core.User) error {
	if result := r.db.Create(&user); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormUserRepository) SaveAdmin(admin core.Admin) error {
	if result := r.db.Create(&admin); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormUserRepository) LoginAdmin(username string, password string) error {
	var admin core.Admin
	if result := r.db.Where("username = ? AND password = ? ", username, password).First(&admin); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormUserRepository) GetUserData(email string) (core.User, error) {
	var user core.User
	if result := r.db.Where("email = ?", email).First(&user); result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (r *GormUserRepository) SaveVerify(verifly core.Verification) error {
	if result := r.db.Create(&verifly); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormUserRepository) VerificationOTP(email string, OTP string) error {
	var verifly core.Verification
	if result := r.db.Where("email = ? AND otp = ?", email, OTP).First(&verifly); result.Error != nil {
		return result.Error
	}
	return nil
}
