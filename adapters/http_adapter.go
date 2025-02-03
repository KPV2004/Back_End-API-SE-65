package adapters

import (
	"fmt"
	"go-server/core"

	"github.com/gofiber/fiber/v2"
)

type HttpUserHandler struct {
	service core.UserService
}

func NewHttpUserHandler(service core.UserService) *HttpUserHandler {
	return &HttpUserHandler{service: service}
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
	email := c.Params("email")
	fmt.Println(email)
	user, err := h.service.FindByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}
