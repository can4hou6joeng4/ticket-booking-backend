package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

type TicketHandler struct {
	repository models.TicketRepository
}

// @Summary      Create new ticket
// @Description  Create a new ticket for an event
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ticket body models.Ticket true "Ticket object"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      422  {object}  Response
// @Router       /api/ticket [post]
func (h *TicketHandler) CreateOne(ctx *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ticket := &models.Ticket{}
	userId := ctx.Locals("userId").(uint)
	// ctx.BodyParser(ticket) 是 Fiber 框架提供的一个方法
	// 它的作用是将 HTTP 请求的 请求体（Body） 自动解析并绑定到 Go 结构体（event 变量）上。
	if err := ctx.BodyParser(ticket); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	ticket, err := h.repository.CreateOne(context, userId, ticket)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status":  "success",
		"message": "Ticket created successfully",
		"data":    ticket,
	})
}

// @Summary      Get ticket by ID
// @Description  Retrieve a specific ticket by its ID with QR code
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ticketId path int true "Ticket ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Router       /api/ticket/{ticketId} [get]
func (h *TicketHandler) GetOne(ctx *fiber.Ctx) error {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ticketId, _ := strconv.Atoi(ctx.Params("ticketId"))
	userId := ctx.Locals("userId").(uint)
	ticket, err := h.repository.GetOne(context, userId, uint(ticketId))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	var QRcode []byte
	QRcode, err = qrcode.Encode(
		fmt.Sprintf("ticketId:%d,ownerId:%d", ticketId, userId),
		qrcode.Medium,
		256,
	)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "",
		"data": &fiber.Map{
			"ticket": ticket,
			"qrcode": QRcode,
		},
	})
}

// @Summary      Get all tickets
// @Description  Retrieve all tickets for the authenticated user
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Router       /api/ticket [get]
func (h *TicketHandler) GetMany(ctx *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	userId := ctx.Locals("userId").(uint)
	tickets, err := h.repository.GetMany(context, userId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "",
		"data":    tickets,
	})
}

// @Summary      Validate ticket
// @Description  Validate a ticket by its ID
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ticketId path int true "Ticket ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Router       /api/ticket/{ticketId}/validate [post]
func (h *TicketHandler) ValidateOne(ctx *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	validateBody := &models.ValidateTicket{}
	if err := ctx.BodyParser(validateBody); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	validateData := make(map[string]interface{})
	validateData["entered"] = true
	ticket, err := h.repository.UpdateOne(context, validateBody.OwnerId, validateBody.TicketId, validateData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
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
