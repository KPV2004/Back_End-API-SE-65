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
func (r *GormUserRepository) UpdateUser(user core.User, email string) error {
	if result := r.db.Model(&user).Where("email = ?", email).Updates(&user); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormUserRepository) UpdatePlane(plan core.Plan, planID string) error {
	if result := r.db.Model(&plan).Where("plan_id = ?", planID).Updates(&plan); result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *GormUserRepository) UpdateUserPlanByEmail(email string, newPlanID string) error {
	var user core.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	for _, plan := range user.UserPlanID {
		if plan == newPlanID {
			// หากมีอยู่แล้ว ไม่ต้องอัปเดต
			return nil
		}
	}
	// เพิ่ม newPlanID เข้าไปใน slice ของ UserPlanID
	user.UserPlanID = append(user.UserPlanID, newPlanID)

	// อัปเดต field user_plan_id โดยระบุเงื่อนไขด้วย email แทนการใช้ Save()
	if err := r.db.Model(&core.User{}).
		Where("email = ?", email).
		Update("user_plan_id", user.UserPlanID).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) CreatePlan(userPlan core.Plan) error {
	if result := r.db.Create(&userPlan); result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *GormUserRepository) AddTripLocation(planID string, newLocation core.TripLocation, index int) error {
	var plan core.Plan

	// ดึงข้อมูลแผนจากฐานข้อมูลตาม planID
	if err := r.db.First(&plan, "plan_id = ?", planID).Error; err != nil {
		return err
	}

	// ตรวจสอบว่า TripLocation เป็น nil หรือไม่ ถ้าเป็น nil ให้ตั้งค่าเป็น array ว่าง ๆ
	if plan.TripLocation == nil {
		plan.TripLocation = []core.TripLocation{}
	}

	// ถ้า index ไม่อยู่ในขอบเขต (น้อยกว่า 0 หรือมากกว่าความยาวของ slice) ให้ append
	if index < 0 || index > len(plan.TripLocation) {
		plan.TripLocation = append(plan.TripLocation, newLocation)
	} else {
		// Insert newLocation ที่ตำแหน่ง index โดยเลื่อน element ที่มีอยู่ไปด้านหลัง
		plan.TripLocation = append(
			plan.TripLocation[:index],
			append([]core.TripLocation{newLocation}, plan.TripLocation[index:]...)...,
		)
	}

	// Save the updated plan
	if err := r.db.Save(&plan).Error; err != nil {
		return err
	}

	return nil
}

// Updated repository method using the new model type
func (r *GormUserRepository) GetTripLocationByPlanID(planID string) ([]core.TripLocation, error) {
	var plan core.Plan
	// Retrieve the plan from the database by planID
	if err := r.db.First(&plan, "plan_id = ?", planID).Error; err != nil {
		return nil, err
	}
	// Return the TripLocation slice
	return plan.TripLocation, nil
}

func (r *GormUserRepository) GetPlanByID(planID string) (core.Plan, error) {
	var plan core.Plan
	if err := r.db.First(&plan, "plan_id = ?", planID).Error; err != nil {
		return plan, err
	}
	return plan, nil
}
func (r *GormUserRepository) DeletePlanByID(planID string) error {
	// Delete the plan record from the "plans" table
	if err := r.db.Delete(&core.Plan{}, "plan_id = ?", planID).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormUserRepository) DeleteUserPlanByEmail(email, planID string) error {
	err := r.db.Exec(
		"UPDATE users SET user_plan_id = array_remove(user_plan_id, ?) WHERE email = ?",
		planID, email,
	).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *GormUserRepository) GetVisiblePlans() ([]core.Plan, error) {
	var plans []core.Plan
	if err := r.db.Where("visibility = ?", true).Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}
func (r *GormUserRepository) DeleteTripLocation(planID, targetPlaceID string) error {
	var plan core.Plan
	if err := r.db.First(&plan, "plan_id = ?", planID).Error; err != nil {
		return err
	}

	var updatedLocations []core.TripLocation
	found := false
	for _, loc := range plan.TripLocation {
		if loc.PlaceID == targetPlaceID {
			found = true
			continue
		}
		updatedLocations = append(updatedLocations, loc)
	}

	if !found {
		return fmt.Errorf("TripLocation with PlaceID %s not found", targetPlaceID)
	}

	// อัปเดต TripLocation ด้วย slice ที่ถูกกรองแล้ว
	plan.TripLocation = updatedLocations

	if err := r.db.Save(&plan).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) UpdateAuthorImg(planID, newImg string) error {
	var plan core.Plan
	// Retrieve the plan by planID
	if err := r.db.First(&plan, "plan_id = ?", planID).Error; err != nil {
		return err
	}

	// Update the AuthorImg field
	plan.AuthorImg = newImg

	// Save the changes
	if err := r.db.Save(&plan).Error; err != nil {
		return err
	}
	return nil
}
