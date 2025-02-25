package adapters

import (
	"fmt"
	"go-server/core"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HttpUserHandler struct {
	service       core.UserService
	service_email core.EmailService
}

func NewHttpUserHandler(service core.UserService, service_email core.EmailService) *HttpUserHandler {
	return &HttpUserHandler{service: service, service_email: service_email}
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
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/admin/login [post]
func (h *HttpUserHandler) LoginAdmin(c *fiber.Ctx) error {
	var admin core.Admin

	if err := c.BodyParser(&admin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request"})
	}
	if err := h.service.LoginAdmin(admin); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Login Sucessfully!"})
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
