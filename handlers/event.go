package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
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
// @Success      200  {object}  Response
// @Failure      500  {object}  Response
// @Router       /api/event [get]
func (h *EventHandler) GetMany(ctx *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	events, err := h.repository.GetMany(context)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":   "fail",
			"messages": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   "success",
		"messages": "",
		"data":     events,
	})
}

// @Summary      Get event by ID
// @Description  Retrieve a specific event by its ID
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        eventId path int true "Event ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Router       /api/event/{eventId} [get]
func (h *EventHandler) GetOne(ctx *fiber.Ctx) error {
	eventId, _ := strconv.Atoi(ctx.Params("eventId"))
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	event, err := h.repository.GetOne(context, eventId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":   "fail",
			"messages": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   "success",
		"messages": "",
		"data":     event,
	})
}

// @Summary      Create new event
// @Description  Create a new event in the system
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        event body models.Event true "Event object"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Failure      422  {object}  Response
// @Router       /api/event [post]
func (h *EventHandler) CreateOne(ctx *fiber.Ctx) error {
	event := &models.Event{}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// ctx.BodyParser(event) 是 Fiber 框架提供的一个方法
	// 它的作用是将 HTTP 请求的 请求体（Body） 自动解析并绑定到 Go 结构体（event 变量）上。
	if err := ctx.BodyParser(event); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":   "fail",
			"messages": err.Error(),
		})
	}
	event, err := h.repository.CreateOne(context, event)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":   "fail",
			"messages": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":   "success",
		"messages": "Event created successfully",
		"data":     event,
	})
}

// @Summary      Update event
// @Description  Update an existing event by its ID
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        eventId path int true "Event ID"
// @Param        event body map[string]interface{} true "Event update data"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      422  {object}  Response
// @Router       /api/event/{eventId} [put]
func (h *EventHandler) UpdateOne(ctx *fiber.Ctx) error {
	eventId, _ := strconv.Atoi(ctx.Params("eventId"))
	updateData := make(map[string]interface{})
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ctx.BodyParser(&updateData); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":   "fail",
			"messages": err.Error(),
			"data":     nil,
		})
	}
	event, err := h.repository.UpdateOne(context, eventId, updateData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":   "fail",
			"messages": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   "success",
		"messages": "Event updated successfully",
		"data":     event,
	})

}

// @Summary      Delete event
// @Description  Delete an event by its ID
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        eventId path int true "Event ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Router       /api/event/{eventId} [delete]
func (h *EventHandler) DeleteOne(ctx *fiber.Ctx) error {
	eventId, _ := strconv.Atoi(ctx.Params("eventId"))
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.repository.DeleteOne(context, eventId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":   "fail",
			"messages": err.Error(),
		})
	}
	return ctx.SendStatus(fiber.StatusNoContent)
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
