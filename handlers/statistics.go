package handlers

import (
	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/utils"
	"github.com/gofiber/fiber/v2"
)

type StatisticsHandler struct {
	repository models.StatisticsRepository
}

// @Summary      Get dashboard statistics
// @Description  Retrieve statistics for the dashboard
// @Tags         statistics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/statistics/dashboard [get]
func (h *StatisticsHandler) GetDashboardStatistics(ctx *fiber.Ctx) error {
	context, cancel := utils.CreateTimeoutContext(0)
	defer cancel()
	count, err := h.repository.GetCount(context)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "", count)
}

func NewStatisticsHandler(router fiber.Router, repository models.StatisticsRepository) {
	handler := &StatisticsHandler{
		repository: repository,
	}

	router.Get("/dashboard", handler.GetDashboardStatistics)
}
