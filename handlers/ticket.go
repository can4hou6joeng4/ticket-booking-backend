package handlers

import (
	"encoding/json"
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
		fmt.Sprintf("qrCode:ticketId:%d,ownerId:%d", ticket.ID, userId),
		getQRLevel(h.config.QRConfig.QRLevel),
		h.config.QRConfig.QRSize,
	)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	// 获取活动信息并设置过期时间
	expiration := time.Until(event.EndDate)
	if expiration < 0 {
		expiration = 0
	}

	// 异步缓存票据信息和QRCode
	go func() {
		// 创建新的上下文用于异步操作
		asyncCtx, cancel := utils.CreateTimeoutContext(60 * time.Second)
		defer cancel()

		// 1. 缓存票据信息（不包含QRCode）
		// 包含Event信息
		eventData := map[string]interface{}{
			"ID":                    event.ID,
			"Name":                  event.Name,
			"Location":              event.Location,
			"Date":                  event.Date,
			"EndDate":               event.EndDate,
			"TotalTicketsPurchased": event.TotalTicketsPurchased,
			"TotalTicketsEntered":   event.TotalTicketsEntered,
		}

		ticketData := map[string]interface{}{
			"ID":        ticket.ID,
			"UserID":    ticket.UserID,
			"EventID":   ticket.EventID,
			"Entered":   ticket.Entered,
			"CreatedAt": ticket.CreatedAt,
			"UpdatedAt": ticket.UpdatedAt,
			"Event":     eventData, // 添加Event信息
		}

		// 序列化票据数据
		ticketJSON, err := json.Marshal(ticketData)
		if err == nil {
			// 票据信息缓存键
			ticketCacheKey := fmt.Sprintf("ticket:info:%d:user:%d", ticket.ID, userId)
			// 缓存票据信息，与活动结束时间相同的过期时间
			h.redis.Set(asyncCtx, ticketCacheKey, ticketJSON, expiration)
		}

		// 2. 缓存QRCode
		qrCodeKey := fmt.Sprintf("qrCode:ticketId:%d,ownerId:%d", ticket.ID, userId)
		h.redis.Set(asyncCtx, qrCodeKey, QRcode, expiration)

		// 3. 删除用户票据列表缓存，确保下次获取列表时能获取最新数据
		userTicketsKey := fmt.Sprintf("tickets:user:%d", userId)
		h.redis.Del(asyncCtx, userTicketsKey)
	}()

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

	var ticket *models.Ticket
	var QRcode []byte
	var err error

	// 尝试从Redis获取票据信息
	ticketInfoKey := fmt.Sprintf("ticket:info:%d:user:%d", ticketId, userId)
	cachedTicketData, err := h.redis.Get(context, ticketInfoKey).Bytes()
	if err == nil {
		// 票据信息缓存命中
		var ticketData map[string]interface{}
		if err := json.Unmarshal(cachedTicketData, &ticketData); err == nil {
			// 手动构建票据对象
			ticket = &models.Ticket{
				ID:      uint(ticketData["ID"].(float64)),
				UserID:  uint(ticketData["UserID"].(float64)),
				EventID: uint(ticketData["EventID"].(float64)),
				Entered: ticketData["Entered"].(bool),
			}

			// 处理Event信息
			if eventData, ok := ticketData["Event"].(map[string]interface{}); ok {
				ticket.Event = models.Event{
					ID:       uint(eventData["ID"].(float64)),
					Name:     eventData["Name"].(string),
					Location: eventData["Location"].(string),
				}

				// 处理票券统计信息
				if tp, ok := eventData["TotalTicketsPurchased"].(float64); ok {
					ticket.Event.TotalTicketsPurchased = int64(tp)
				}

				if te, ok := eventData["TotalTicketsEntered"].(float64); ok {
					ticket.Event.TotalTicketsEntered = int64(te)
				}

				// 处理Event的日期字段
				if dateStr, ok := eventData["Date"].(string); ok {
					date, err := time.Parse(time.RFC3339, dateStr)
					if err == nil {
						ticket.Event.Date = date
					}
				}

				if endDateStr, ok := eventData["EndDate"].(string); ok {
					endDate, err := time.Parse(time.RFC3339, endDateStr)
					if err == nil {
						ticket.Event.EndDate = endDate
					}
				}
			}

			// 处理时间字段
			if createdAtStr, ok := ticketData["CreatedAt"].(string); ok {
				createdAt, err := time.Parse(time.RFC3339, createdAtStr)
				if err == nil {
					ticket.CreatedAt = createdAt
				}
			}

			if updatedAtStr, ok := ticketData["UpdatedAt"].(string); ok {
				updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
				if err == nil {
					ticket.UpdatedAt = updatedAt
				}
			}
		}
	}

	// 如果票据信息缓存未命中，从数据库获取
	if ticket == nil {
		ticket, err = h.ticketRepository.GetOne(context, userId, uint(ticketId))
		if err != nil {
			return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
		}

		// 异步缓存票据信息
		go func() {
			// 创建新的上下文用于异步操作
			ctx, cancel := utils.CreateTimeoutContext(60 * time.Second)
			defer cancel()

			// 包含Event信息
			eventData := map[string]interface{}{
				"ID":                    ticket.Event.ID,
				"Name":                  ticket.Event.Name,
				"Location":              ticket.Event.Location,
				"Date":                  ticket.Event.Date,
				"EndDate":               ticket.Event.EndDate,
				"TotalTicketsPurchased": ticket.Event.TotalTicketsPurchased,
				"TotalTicketsEntered":   ticket.Event.TotalTicketsEntered,
			}

			ticketData := map[string]interface{}{
				"ID":        ticket.ID,
				"UserID":    ticket.UserID,
				"EventID":   ticket.EventID,
				"Entered":   ticket.Entered,
				"CreatedAt": ticket.CreatedAt,
				"UpdatedAt": ticket.UpdatedAt,
				"Event":     eventData, // 添加Event信息
			}

			ticketJSON, err := json.Marshal(ticketData)
			if err == nil {
				h.redis.Set(ctx, ticketInfoKey, ticketJSON, time.Hour)
			}
		}()
	}

	// 从Redis获取二维码
	qrCodeKey := fmt.Sprintf("qrCode:ticketId:%d,ownerId:%d", ticketId, userId)
	QRcode, err = h.redis.Get(context, qrCodeKey).Bytes()
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	// 若QRCode为空，则表示二维码已过期
	if len(QRcode) == 0 {
		return utils.ErrorResponseWithData(ctx, fiber.StatusBadRequest, fmt.Errorf("QR code expired"), map[string]interface{}{
			"message": "活动已过期",
		})
	}

	// 构建响应数据
	responseData := map[string]interface{}{
		"ticket": ticket,
		"qrcode": QRcode,
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "", responseData)
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

	// 尝试从Redis获取缓存
	cacheKey := fmt.Sprintf("tickets:user:%d", userId)
	cachedData, err := h.redis.Get(context, cacheKey).Bytes()
	if err == nil {
		// 缓存命中
		var tickets []*models.Ticket
		if err := json.Unmarshal(cachedData, &tickets); err == nil {
			return utils.SuccessResponse(ctx, fiber.StatusOK, "", tickets)
		}
	}

	// 从数据库获取
	tickets, err := h.ticketRepository.GetMany(context, userId)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	// 异步缓存结果
	go func() {
		ctx, cancel := utils.CreateTimeoutContext(60 * time.Second)
		defer cancel()

		// 序列化票据数据
		ticketsJSON, err := json.Marshal(tickets)
		if err != nil {
			return
		}

		// 设置一个合理的过期时间（例如1小时）
		h.redis.Set(ctx, cacheKey, ticketsJSON, time.Hour)
	}()

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

	// 异步更新或删除相关缓存
	go func() {
		ctx, cancel := utils.CreateTimeoutContext(60 * time.Second)
		defer cancel()

		// 票据相关缓存键
		ticketInfoKey := fmt.Sprintf("ticket:info:%d:user:%d", validateBody.TicketId, validateBody.OwnerId)
		userTicketsKey := fmt.Sprintf("tickets:user:%d", validateBody.OwnerId)

		// 事件相关缓存键
		eventCacheKey := fmt.Sprintf("event:%d", ticket.EventID)

		// 删除缓存以确保下次获取最新数据
		h.redis.Del(ctx, ticketInfoKey, userTicketsKey, eventCacheKey)
	}()

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
