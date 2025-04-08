package handlers

import (
	"fmt"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service models.AuthService
}

var validate = validator.New()

// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body models.AuthCredentials true "Login credentials"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	creds := &models.AuthCredentials{}
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	if err := ctx.BodyParser(&creds); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	if err := validate.Struct(creds); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	token, user, err := h.service.Login(context, creds)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	return utils.SuccessResponse(ctx, fiber.StatusOK, "Successfully logged in", map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// @Summary      Register new user
// @Description  Create a new user account and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body models.AuthCredentials true "Registration credentials"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /api/auth/register [post]
func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	creds := &models.AuthCredentials{}
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	if err := ctx.BodyParser(&creds); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	if err := validate.Struct(creds); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, fmt.Errorf("please provide a valid email and password"))
	}
	token, user, err := h.service.Register(context, creds)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	return utils.SuccessResponse(ctx, fiber.StatusOK, "Successfully registered", map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// @Summary      Logout user
// @Description  Logout user and invalidate session
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /api/auth/logout [post]
func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()

	userId := ctx.Locals("userId").(uint)
	if err := h.service.Logout(context, userId); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "Successfully logged out", nil)
}

func NewAuthHandler(router fiber.Router, service models.AuthService) {
	handler := &AuthHandler{
		service: service,
	}
	// 公开路由
	router.Post("/login", handler.Login)
	router.Post("/register", handler.Register)
}

func NewAuthProtectedHandler(router fiber.Router, service models.AuthService) {
	handler := &AuthHandler{
		service: service,
	}
	// 需要认证的路由
	router.Post("/logout", handler.Logout)
}
