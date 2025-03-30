package handlers

import (
	"context"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

type TicketHandler struct {
	repository models.TicketRepository
}

func (h *TicketHandler) CreateOne(ctx *fiber.Ctx) error {
	ticket := &models.Ticket{}
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	// ctx.BodyParser(ticket) 是 Fiber 框架提供的一个方法
	// 它的作用是将 HTTP 请求的 请求体（Body） 自动解析并绑定到 Go 结构体（event 变量）上。
	if err := ctx.BodyParser(ticket); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	ticket, err := h.repository.CreateOne(context, ticket)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Ticket created successfully",
		"data":    ticket,
	})
}
func (h *TicketHandler) GetOne(ctx *fiber.Ctx) error {
	ticketId, _ := strconv.Atoi(ctx.Params("ticketId"))
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ticket, err := h.repository.GetOne(context, uint(ticketId))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    ticket,
	})
}
func (h *TicketHandler) GetMany(ctx *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tickets, err := h.repository.GetMany(context)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    tickets,
	})
}
func (h *TicketHandler) ValidateOne(ctx *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	validateBody := &models.ValidateTicket{}
	if err := ctx.BodyParser(validateBody); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	validateData := make(map[string]interface{})
	validateData["entered"] = true
	ticket, err := h.repository.UpdateOne(context, validateBody.TicketId, validateData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Welcome to the show",
		"data":    ticket,
	})
}
func NewTicketHandler(router fiber.Router, repository models.TicketRepository) {
	handler := &TicketHandler{
		repository: repository,
	}
	router.Post("/", handler.CreateOne)
	router.Get("/:ticketId", handler.GetOne)
	router.Get("/", handler.GetMany)
	router.Post("/validate", handler.ValidateOne)
}
