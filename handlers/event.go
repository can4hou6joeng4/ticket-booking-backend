package handlers

import (
	"strconv"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/utils"
	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	repository models.EventRepository
}

// @Summary      Get all events
// @Description  Retrieve all events from the system
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/event [get]
func (h *EventHandler) GetMany(ctx *fiber.Ctx) error {
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	events, err := h.repository.GetMany(context)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}
	return utils.SuccessResponse(ctx, fiber.StatusOK, "", events)
}

// @Summary      Get event by ID
// @Description  Retrieve a specific event by its ID
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        eventId path int true "Event ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /api/event/{eventId} [get]
func (h *EventHandler) GetOne(ctx *fiber.Ctx) error {
	eventId, _ := strconv.Atoi(ctx.Params("eventId"))
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	event, err := h.repository.GetOne(context, eventId)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	return utils.SuccessResponse(ctx, fiber.StatusOK, "", event)
}

// @Summary      Create new event
// @Description  Create a new event in the system
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        event body models.Event true "Event object"
// @Success      201  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      422  {object}  utils.Response
// @Router       /api/event [post]
func (h *EventHandler) CreateOne(ctx *fiber.Ctx) error {
	event := &models.Event{}
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	// ctx.BodyParser(event) 是 Fiber 框架提供的一个方法
	// 它的作用是将 HTTP 请求的 请求体（Body） 自动解析并绑定到 Go 结构体（event 变量）上。
	if err := ctx.BodyParser(event); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusUnprocessableEntity, err)
	}
	event, err := h.repository.CreateOne(context, event)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	return utils.SuccessResponse(ctx, fiber.StatusCreated, "Event created successfully", event)
}

// @Summary      Update event
// @Description  Update an existing event by its ID
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        eventId path int true "Event ID"
// @Param        event body map[string]interface{} true "Event update data"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      422  {object}  utils.Response
// @Router       /api/event/{eventId} [put]
func (h *EventHandler) UpdateOne(ctx *fiber.Ctx) error {
	eventId, _ := strconv.Atoi(ctx.Params("eventId"))
	updateData := make(map[string]interface{})
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	if err := ctx.BodyParser(&updateData); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusUnprocessableEntity, err)
	}
	event, err := h.repository.UpdateOne(context, eventId, updateData)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	return utils.SuccessResponse(ctx, fiber.StatusOK, "Event updated successfully", event)
}

// @Summary      Delete event
// @Description  Delete an event by its ID
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        eventId path int true "Event ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /api/event/{eventId} [delete]
func (h *EventHandler) DeleteOne(ctx *fiber.Ctx) error {
	eventId, _ := strconv.Atoi(ctx.Params("eventId"))
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	if err := h.repository.DeleteOne(context, eventId); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	return utils.NoContentResponse(ctx)
}

func NewEventHandler(router fiber.Router, repository models.EventRepository) {
	handler := &EventHandler{
		repository: repository,
	}
	router.Get("/", handler.GetMany)
	router.Post("/", handler.CreateOne)
	router.Get("/:eventId", handler.GetOne)
	router.Put("/:eventId", handler.UpdateOne)
	router.Delete("/:eventId", handler.DeleteOne)
}
