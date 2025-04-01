package handlers

import (
	"context"
	"time"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"github.com/gofiber/fiber/v2"
)

type StatisticsHandler struct {
	repository models.StatisticsRepository
}

func (h *StatisticsHandler) GetDashboardStatistics(ctx *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := h.repository.GetCount(context)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get statistics",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   "success",
		"messages": "",
		"data":     count,
	})
}
func NewStatisticsHandler(router fiber.Router, repository models.StatisticsRepository) {
	handler := &StatisticsHandler{
		repository: repository,
	}

	router.Get("/dashboard", handler.GetDashboardStatistics)
}
