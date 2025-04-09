package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

type EventHandler struct {
	repository models.EventRepository
	redis      *redis.Client
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

	// 尝试从Redis获取缓存
	keys, err := h.redis.Keys(context, "event:*").Result()
	if err != nil && err != redis.Nil {
		log.Error("从Redis获取事件键失败:", err)
	} else if len(keys) > 0 {
		cachedEvents := make([]*models.Event, 0, len(keys))
		for _, key := range keys {
			event, err := h.getEventFromCache(context, extractEventId(key))
			if err != nil {
				// 只有在非缓存未命中的情况下记录错误
				if err != redis.Nil {
					log.Error(fmt.Sprintf("获取缓存事件失败 %s: %v", key, err))
				}
				continue
			}
			cachedEvents = append(cachedEvents, event)
		}
		if len(cachedEvents) > 0 {
			return utils.SuccessResponse(ctx, fiber.StatusOK, "", cachedEvents)
		}
	}

	// 从数据库获取并缓存
	events, err := h.repository.GetMany(context)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	// 异步缓存事件
	go func() {
		ctx, cancel := utils.CreateTimeoutContext(60 * time.Second)
		defer cancel()
		for _, event := range events {
			if err := h.cacheEvent(ctx, event); err != nil {
				log.Error(err)
			}
			// 增加间隔，避免Redis过载
			time.Sleep(100 * time.Millisecond)
		}
	}()

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

	// 尝试从缓存获取
	event, err := h.getEventFromCache(context, eventId)
	if err == nil {
		return utils.SuccessResponse(ctx, fiber.StatusOK, "", event)
	}

	// 从数据库获取
	event, err = h.repository.GetOne(context, eventId)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	// 异步缓存
	go func() {
		ctx, cancel := utils.CreateTimeoutContext(60 * time.Second)
		defer cancel()
		if err := h.cacheEvent(ctx, event); err != nil {
			log.Error(err)
		}
	}()

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

	// 异步缓存事件
	go func() {
		ctx, cancel := utils.CreateTimeoutContext(60 * time.Second)
		defer cancel()
		if err := h.cacheEvent(ctx, event); err != nil {
			log.Error(err)
		}
	}()

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

	// 如果有endDate字段，修改为数据库列名end_date
	if endDate, ok := updateData["endDate"]; ok {
		updateData["end_date"] = endDate
		delete(updateData, "endDate")
	}

	event, err := h.repository.UpdateOne(context, eventId, updateData)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	// 异步缓存更新后的事件
	go func() {
		ctx, cancel := utils.CreateTimeoutContext(60 * time.Second)
		defer cancel()
		if err := h.cacheEvent(ctx, event); err != nil {
			log.Error(err)
		}
	}()

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

	// 删除事件的缓存
	go func() {
		ctx, cancel := utils.CreateTimeoutContext(60 * time.Second)
		defer cancel()
		key := fmt.Sprintf("event:%d", eventId)
		if err := h.redis.Del(ctx, key).Err(); err != nil && err != redis.Nil {
			log.Error(fmt.Sprintf("删除事件缓存失败 ID=%d: %v", eventId, err))
		}
	}()

	return utils.NoContentResponse(ctx)
}

// cacheEvent 将单个事件缓存到Redis
func (h *EventHandler) cacheEvent(ctx context.Context, event *models.Event) error {
	// 创建简化版的事件对象，只包含必要字段以减小数据大小
	cacheEvent := map[string]interface{}{
		"ID":                    event.ID,
		"Name":                  event.Name,
		"Location":              event.Location,
		"Date":                  event.Date,
		"EndDate":               event.EndDate,
		"TotalTicketsPurchased": event.TotalTicketsPurchased,
		"TotalTicketsEntered":   event.TotalTicketsEntered,
	}

	// 序列化简化版事件对象
	eventJSON, err := json.Marshal(cacheEvent)
	if err != nil {
		return fmt.Errorf("序列化事件数据失败 ID=%d: %v", event.ID, err)
	}

	// 设置过期时间
	expiration := time.Until(event.EndDate)
	if expiration < 0 {
		expiration = time.Hour
	}

	// 生成缓存的key
	key := fmt.Sprintf("event:%d", event.ID)

	// 存储到Redis
	if err := h.redis.Set(ctx, key, eventJSON, expiration).Err(); err != nil {
		return fmt.Errorf("缓存事件到Redis失败 ID=%d: %v", event.ID, err)
	}

	return nil
}

// getEventFromCache 从Redis缓存中获取事件
func (h *EventHandler) getEventFromCache(ctx context.Context, eventId int) (*models.Event, error) {
	key := fmt.Sprintf("event:%d", eventId)

	eventData, err := h.redis.Get(ctx, key).Result()
	if err != nil {
		// 仅在调试级别记录缓存未命中
		if err != redis.Nil {
			log.Error(fmt.Sprintf("从Redis获取事件数据失败 ID=%d: %v", eventId, err))
		}
		return nil, err
	}

	var cacheEvent map[string]interface{}
	if err := json.Unmarshal([]byte(eventData), &cacheEvent); err != nil {
		log.Error(fmt.Sprintf("解析事件数据失败 ID=%d: %v", eventId, err))
		return nil, err
	}

	// 创建一个完整的事件对象
	event := &models.Event{
		ID:       uint(cacheEvent["ID"].(float64)), // JSON将数字解析为float64
		Name:     cacheEvent["Name"].(string),
		Location: cacheEvent["Location"].(string),
	}

	// 处理票券统计信息
	if tp, ok := cacheEvent["TotalTicketsPurchased"].(float64); ok {
		event.TotalTicketsPurchased = int64(tp)
	}

	if te, ok := cacheEvent["TotalTicketsEntered"].(float64); ok {
		event.TotalTicketsEntered = int64(te)
	}

	// 处理日期字段
	if dateStr, ok := cacheEvent["Date"].(string); ok {
		date, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			log.Error(fmt.Sprintf("解析Date字段失败 ID=%d: %v", eventId, err))
			return nil, err
		}
		event.Date = date
	}

	if endDateStr, ok := cacheEvent["EndDate"].(string); ok {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			log.Error(fmt.Sprintf("解析EndDate字段失败 ID=%d: %v", eventId, err))
			return nil, err
		}
		event.EndDate = endDate
	}

	return event, nil
}

// extractEventId 从Redis key中提取事件ID
func extractEventId(key string) int {
	idStr := strings.TrimPrefix(key, "event:")
	id, _ := strconv.Atoi(idStr)
	return id
}

func NewEventHandler(router fiber.Router, repository models.EventRepository, redis *redis.Client) {
	handler := &EventHandler{
		repository: repository,
		redis:      redis,
	}
	router.Get("/", handler.GetMany)
	router.Post("/", handler.CreateOne)
	router.Get("/:eventId", handler.GetOne)
	router.Put("/:eventId", handler.UpdateOne)
	router.Delete("/:eventId", handler.DeleteOne)
}
