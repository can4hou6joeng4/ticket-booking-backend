package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/config"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/skip2/go-qrcode"
)

type TicketHandler struct {
	ticketRepository models.TicketRepository
	eventRepository  models.EventRepository
	config           *config.EnvConfig
	redis            *redis.Client
}

// @Summary      Create new ticket
// @Description  Create a new ticket for an event
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ticket body models.Ticket true "Ticket object"
// @Success      201  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      422  {object}  utils.Response
// @Router       /api/ticket [post]
func (h *TicketHandler) CreateOne(ctx *fiber.Ctx) error {
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	ticket := &models.Ticket{}
	userId := ctx.Locals("userId").(uint)
	if err := ctx.BodyParser(ticket); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusUnprocessableEntity, err)
	}
	eventId := int(ticket.EventID)
	// 验证活动是否已经结束
	event, err := h.eventRepository.GetOne(context, eventId)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	if event.EndDate.Before(time.Now()) {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, fmt.Errorf("活动已结束"))
	}
	ticket, err = h.ticketRepository.CreateOne(context, userId, ticket)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	// 生成二维码
	var QRcode []byte
	QRcode, err = qrcode.Encode(
		fmt.Sprintf("ticketId:%d,ownerId:%d", ticket.ID, userId),
		getQRLevel(h.config.QRConfig.QRLevel),
		h.config.QRConfig.QRSize,
	)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	ticket.QRCode = QRcode

	// 获取活动信息并设置过期时间
	expiration := time.Until(event.EndDate)
	if expiration < 0 {
		expiration = 0
	}

	// 存储到Redis
	key := fmt.Sprintf("ticket:%d,ownerId:%d", ticket.ID, userId)
	if err := h.redis.Set(context, key, ticket.QRCode, expiration).Err(); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, fmt.Errorf("failed to store QR code in Redis"))
	}

	return utils.SuccessResponse(ctx, fiber.StatusCreated, "Ticket created successfully", ticket)
}

// @Summary      Get ticket by ID
// @Description  Retrieve a specific ticket by its ID with QR code
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ticketId path int true "Ticket ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /api/ticket/{ticketId} [get]
func (h *TicketHandler) GetOne(ctx *fiber.Ctx) error {
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	ticketId, _ := strconv.Atoi(ctx.Params("ticketId"))
	userId := ctx.Locals("userId").(uint)
	ticket, err := h.ticketRepository.GetOne(context, userId, uint(ticketId))
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	// 从Redis中获取二维码
	QRcode, err := h.redis.Get(context, fmt.Sprintf("ticket:%d,ownerId:%d", ticketId, userId)).Bytes()
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	// 若QRCode为空，则表示二维码已过期
	if len(QRcode) == 0 {
		return utils.ErrorResponseWithData(ctx, fiber.StatusBadRequest, fmt.Errorf("QR code expired"), map[string]interface{}{
			"message": "活动已过期",
		})
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "", map[string]interface{}{
		"ticket": ticket,
		"qrcode": QRcode,
	})
}

func getQRLevel(level string) qrcode.RecoveryLevel {
	switch level {
	case "Low":
		return qrcode.Low
	case "Medium":
		return qrcode.Medium
	case "High":
		return qrcode.High
	case "Highest":
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

// @Summary      Get all tickets
// @Description  Retrieve all tickets for the authenticated user
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /api/ticket [get]
func (h *TicketHandler) GetMany(ctx *fiber.Ctx) error {
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	userId := ctx.Locals("userId").(uint)
	tickets, err := h.ticketRepository.GetMany(context, userId)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	return utils.SuccessResponse(ctx, fiber.StatusOK, "", tickets)
}

// @Summary      Validate ticket
// @Description  Validate a ticket by its ID
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ticketId path int true "Ticket ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /api/ticket/{ticketId}/validate [post]
func (h *TicketHandler) ValidateOne(ctx *fiber.Ctx) error {
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	validateBody := &models.ValidateTicket{}
	if err := ctx.BodyParser(validateBody); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusUnprocessableEntity, err)
	}
	validateData := make(map[string]interface{})
	validateData["entered"] = true
	ticket, err := h.ticketRepository.UpdateOne(context, validateBody.OwnerId, validateBody.TicketId, validateData)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	return utils.SuccessResponse(ctx, fiber.StatusOK, "Welcome to the show", ticket)
}

func NewTicketHandler(router fiber.Router, ticketRepository models.TicketRepository, eventRepository models.EventRepository, config *config.EnvConfig, redis *redis.Client) {
	handler := &TicketHandler{
		ticketRepository: ticketRepository,
		eventRepository:  eventRepository,
		config:           config,
		redis:            redis,
	}
	router.Post("/", handler.CreateOne)
	router.Get("/:ticketId", handler.GetOne)
	router.Get("/", handler.GetMany)
	router.Post("/validate", handler.ValidateOne)
}
