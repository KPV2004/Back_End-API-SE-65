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

	//random 4 digit OTP
	rand.Seed(time.Now().UnixNano())
	otpInt := rand.Intn(9000) + 1000
	opt := strconv.Itoa(otpInt)

	//set otp to user
	user.Otp = opt
	//send otp to user email
	err_mail := h.service_email.Message(opt, user.Email)
	if err_mail != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *HttpUserHandler) GetUser(c *fiber.Ctx) error {
	email := c.Params("email")
	fmt.Println(email)
	user, err := h.service.FindByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}
