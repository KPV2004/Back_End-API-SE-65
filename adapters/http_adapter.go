package adapters

import (
	"encoding/hex"
	"fmt"
	"go-server/core"
	"go-server/middleware"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	"github.com/google/uuid"
)

type HttpUserHandler struct {
	service       core.UserService
	service_email core.EmailService
}

func NewHttpUserHandler(service core.UserService, service_email core.EmailService) *HttpUserHandler {
	return &HttpUserHandler{service: service, service_email: service_email}
}

// generateRandomFilename creates a random filename with the same extension
func generateRandomFilename(ext string) string {
	bytes := make([]byte, 8) // 8 bytes = 16 characters hex
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal("Failed to generate random filename")
	}
	return hex.EncodeToString(bytes) + ext
}

// @Summary Create a new user
// @Description Register a new user in the system
// @Tags User
// @Accept json
// @Produce json
// @Param user body core.User true "User Data"
// @Success 201 {object} core.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/user/register [post]
func (h *HttpUserHandler) RegisterUser(c *fiber.Ctx) error {
	var user core.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request"})
	}
	if err := h.service.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// @Summary Get user by email
// @Description Retrieve user details by email. This route is protected by Firebase authentication.
// @Tags User
// @Accept json
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {object} core.User
// @Failure 404 {object} map[string]string
// @Router /api/v1/user/getuser/{email} [get]
func (h *HttpUserHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Locals("userID") // Get userID from context
	fmt.Println("User ID from context:", userID)
	email := c.Params("email")
	fmt.Println(email)
	user, err := h.service.FindByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary Generate OTP
// @Description Generate a one-time password (OTP) for the user
// @Tags User
// @Accept json
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/user/genotp/{email} [get]
func (h *HttpUserHandler) GenOTP(c *fiber.Ctx) error {
	email := c.Params("email")
	fmt.Println(email)
	var verifly core.Verification

	//random 4 digit OTP
	rand.Seed(time.Now().UnixNano())
	otpInt := rand.Intn(9000) + 1000
	opt := strconv.Itoa(otpInt)

	verifly.Email = email
	verifly.Otp = opt

	if err := h.service.CreateVerifly(verifly); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}
	//send otp to user email

	err_mail := h.service_email.Message(opt, verifly.Email)
	if err_mail != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OTP sent to your email"})
}

// @Summary Verify OTP
// @Description Verify the one-time password (OTP) for the user
// @Tags User
// @Accept json
// @Produce json
// @Param verification body core.Verification true "Verification Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/user/verifyotp [post]
func (h *HttpUserHandler) VerifyOTP(c *fiber.Ctx) error {
	var verify core.Verification
	// email := c.Params("email")
	// otp := c.Params("otp")
	// fmt.Println(email)
	// fmt.Println(otp)
	if err := c.BodyParser(&verify); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request"})
	}

	if err := h.service.VerificationOTP(verify); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invalid OTP"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OTP verified"})
}

// @Summary Register a new admin
// @Description Register a new admin in the system
// @Tags Admin
// @Accept json
// @Produce json
// @Param admin body core.Admin true "Admin Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/admin/register [post]
func (h *HttpUserHandler) RegisterAdmin(c *fiber.Ctx) error {
	var admin core.Admin

	if err := c.BodyParser(&admin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request"})
	}
	if err := h.service.CreateAdmin(admin); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Register Sucessfully!"})
}

// @Summary Admin login
// @Description Login an admin into the system
// @Tags Admin
// @Accept json
// @Produce json
// @Param admin body core.Admin true "Admin Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/admin/login [post]
func (h *HttpUserHandler) LoginAdmin(c *fiber.Ctx) error {
	var admin core.Admin

	// อ่านข้อมูลจาก request
	if err := c.BodyParser(&admin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request"})
	}

	if err := h.service.LoginAdmin(admin); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"sub": admin.Username,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // หมดอายุภายใน 24 ชั่วโมง
		"iat": time.Now().Unix(),
		"jti": uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(middleware.SecretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": t})
}

func (h *HttpUserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}
	return c.JSON(users)
}

// @Summary Update User data
// @Description Update user data in the system
// @Tags User
// @Accept json
// @Produce json
// @Param email path string true "User Email"
// @Param user body core.User true "User Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/user/update/{email} [put]
func (h *HttpUserHandler) UserUpdate(c *fiber.Ctx) error {
	var user core.User
	email := c.Params("email")
	fmt.Println(email)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request"})
	}
	if err := h.service.UpdateUser(user, email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Update User Sucessfully!"})
}

func (h *HttpUserHandler) UploadImage(c *fiber.Ctx) error {
	imageDir := "./access/images"

	// Ensure the directory exists
	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create image directory",
		})
	}

	// Get file from request
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get image",
		})
	}

	// Generate a random filename
	randomName := generateRandomFilename(filepath.Ext(file.Filename))

	// Full path to save file
	savePath := filepath.Join(imageDir, randomName)
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save image",
		})
	}

	// Image URL
	imageURL := fmt.Sprintf("/access/images/%s", randomName)

	// Return success response
	return c.JSON(fiber.Map{
		"message": "Image uploaded successfully",
		"path":    imageURL,
	})
}

func (h *HttpUserHandler) GetAllPlans(c *fiber.Ctx) error {
	plans, err := h.service.GetAllPlans()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve plans"})
	}
	return c.JSON(plans)
}

// @Summary Update User Plan
// @Description Update User Plan by Email
// @Tags User
// @Accept json
// @Produce json
// @Param email path string true "User Email"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/user/updateuserplan/{email} [put]
func (h *HttpUserHandler) UserUpdatePlanByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	newPlanID := c.FormValue("userplan_id")
	if newPlanID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request"})
	}

	if err := h.service.UpdateUserPlanByEmail(email, newPlanID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Update UserPlanID Successfully!"})
}

// @Summary Create Plan
// @Description Create Plan
// @Tags Plan
// @Accept json
// @Produce json
// @Param user body core.Plan true "Plan Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/user/createplan [post]
func (h *HttpUserHandler) CreatePlanTrip(c *fiber.Ctx) error {
	var planData core.Plan
	if err := c.BodyParser(&planData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request"})
	}
	if err := h.service.CreatePlan(planData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Create Plan Successfully!"})
}

// @Summary AddTripLocation
// @Description AddTripLocation
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Plan Id"
// @Param plan body string true "Place Id"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/plan/addtriplocation/{id} [put]
func (h *HttpUserHandler) AddTripLocationHandler(c *fiber.Ctx) error {
	planID := c.Params("id")
	fmt.Print("This" + planID)
	var body core.TripLocation
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid Request",
			"details": err.Error(),
		})
	}

	// ส่ง struct ทั้งก้อนไปที่ AddTripLocation
	if err := h.service.AddTripLocation(planID, body, 0); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Server Error",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Trip location added successfully!"})
}

func (h *HttpUserHandler) GetTripLocationHandler(c *fiber.Ctx) error {
	planID := c.Params("id")
	locations, err := h.service.GetTripLocationByPlanID(planID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Plan not found", "details": err.Error()})
	}

	// Return the TripLocation slice as JSON
	return c.JSON(fiber.Map{"trip_location": locations})
}

func (h *HttpUserHandler) GetPlanByIDHandler(c *fiber.Ctx) error {
	planID := c.Params("id")
	if planID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "plan_id is required"})
	}

	plan, err := h.service.GetPlanByID(planID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Plan not found", "details": err.Error()})
	}

	return c.JSON(fiber.Map{"plan_data": plan})
}

func (h *HttpUserHandler) DeletePlanByIDHandler(c *fiber.Ctx) error {
	planID := c.Params("id")
	if planID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "plan_id is required"})
	}

	if err := h.service.DeletePlanByID(planID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete plan", "details": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Plan deleted successfully"})
}
func (h *HttpUserHandler) DeleteUserPlanByEmailHandler(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}

	type DeleteRequest struct {
		PlanID string `json:"plan_id"`
	}
	var req DeleteRequest
	// ใช้ BodyParser เพื่อแปลง JSON body เป็น struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.PlanID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "plan_id is required"})
	}

	if err := h.service.DeleteUserPlanByEmail(email, req.PlanID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user plan", "details": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Delete plan updated successfully"})
}
func (h *HttpUserHandler) GetVisiblePlansHandler(c *fiber.Ctx) error {
	plans, err := h.service.GetVisiblePlans()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve visible plans",
			"details": err.Error(),
		})
	}
	return c.JSON(plans)
}

func (h *HttpUserHandler) UpdatePlanByID(c *fiber.Ctx) error {
	planID := c.Params("id")
	var planData core.Plan

	if err := c.BodyParser(&planData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if error := h.service.UpdatePlan(planData, planID); error != nil {
		return error
	}
	return c.JSON(fiber.Map{"message": "Update Plan is Sucessfully"})

}

func (h *HttpUserHandler) DeleteTripLocationHandler(c *fiber.Ctx) error {
	planID := c.Params("id")
	if planID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "plan_id is required"})
	}

	// Parse JSON body to get the target place_id
	var body struct {
		PlaceID string `json:"place_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}
	if body.PlaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "place_id is required in the request body"})
	}

	// Call the service layer to delete the trip location
	if err := h.service.DeleteTripLocation(planID, body.PlaceID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete trip location",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Trip location deleted successfully"})
}

func (h *HttpUserHandler) UpdateAuthorImgHandler(c *fiber.Ctx) error {
	// Retrieve planID from URL parameter (e.g., /plan/:id/updateauthorimg)
	planID := c.Params("id")
	if planID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "plan_id is required"})
	}

	// Parse JSON body to get new author_img
	var body struct {
		AuthorImg string `json:"author_img"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if body.AuthorImg == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "author_img is required"})
	}

	// Call the service layer to update the author_img
	if err := h.service.UpdateAuthorImg(planID, body.AuthorImg); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update author_img",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "author_img updated successfully"})
}
