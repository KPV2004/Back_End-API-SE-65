package adapters

import (
	"errors"
	"fmt"
	"go-server/core"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
	hashedPassword, err := hashPassword(admin.Password)
	if err != nil {
		return err
	}
	admin.Password = hashedPassword
	if result := r.db.Create(&admin); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormUserRepository) LoginAdmin(username string, password string) error {
	var admin core.Admin

	// Query database for admin by username
	if result := r.db.Select("password").Where("username = ?", username).First(&admin); result.Error != nil {
		return errors.New("login fail") // Return "login fail" if user is not found
	}
	fmt.Println(admin.Password)
	// Check if the password matches the stored hash
	if !CheckPasswordHash(password, admin.Password) {
		return errors.New("login fail") // Return "login fail" if password is incorrect
	}

	return nil // Login successful
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
