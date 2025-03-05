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

func (r *GormUserRepository) AddTripLocation(planID string, newPlaceID string, index int) error {
	var plan core.Plan
	// ดึงข้อมูลแผนจากฐานข้อมูลตาม planID
	if err := r.db.First(&plan, "plan_id = ?", planID).Error; err != nil {
		return err
	}

	// ตรวจสอบว่า TripLocation เป็น nil หรือไม่ ถ้าเป็น nil ให้ตั้งค่าเป็น array ว่าง ๆ
	if plan.TripLocation == nil {
		plan.TripLocation = [][]string{}
	}

	// ตรวจสอบว่า index อยู่ในขอบเขตของ TripLocation หรือไม่
	if index >= len(plan.TripLocation) {
		// ถ้าไม่อยู่ในขอบเขต ให้เพิ่ม index ใหม่ใน TripLocation
		plan.TripLocation = append(plan.TripLocation, []string{})
	}

	// เพิ่ม newPlaceID ลงในตำแหน่งที่ต้องการ (ตาม index)
	plan.TripLocation[index] = append(plan.TripLocation[index], newPlaceID)

	// ใช้ Save แทน Updates() เพื่อให้ GORM อัปเดตข้อมูลใน TripLocation ได้ถูกต้อง
	if err := r.db.Save(&plan).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) GetTripLocationByPlanID(planID string) ([][]string, error) {
	var plan core.Plan
	// ดึงข้อมูลแผนจากฐานข้อมูลตาม planID
	if err := r.db.First(&plan, "plan_id = ?", planID).Error; err != nil {
		return nil, err
	}
	// คืนค่าข้อมูล TripLocation ที่เป็น 2D array
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
