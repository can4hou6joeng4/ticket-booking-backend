package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service models.AuthService
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var validate = validator.New()

// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body models.AuthCredentials true "Login credentials"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	creds := &models.AuthCredentials{}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ctx.BodyParser(&creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}
	if err := validate.Struct(creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}
	token, user, err := h.service.Login(context, creds)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully logged in",
		"data": &fiber.Map{
			"token": token,
			"user":  user,
		},
	})
}

// @Summary      Register new user
// @Description  Create a new user account and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body models.AuthCredentials true "Registration credentials"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Router       /api/auth/register [post]
func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	creds := &models.AuthCredentials{}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ctx.BodyParser(&creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}
	if err := validate.Struct(creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Errorf("please provide a valid email and password").Error(),
		})
	}
	token, user, err := h.service.Register(context, creds)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully logged in",
		"data": &fiber.Map{
			"token": token,
			"user":  user,
		},
	})
}

func NewAuthHandler(router fiber.Router, service models.AuthService) {
	handler := &AuthHandler{
		service: service,
	}
	router.Post("/login", handler.Login)
	router.Post("/register", handler.Register)
}
