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
