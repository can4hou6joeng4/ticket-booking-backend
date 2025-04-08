package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response 定义统一的响应结构
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 返回错误响应
func ErrorResponse(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(&Response{
		Status:  "fail",
		Message: err.Error(),
	})
}

// ErrorResponseWithData 返回带数据的错误响应
func ErrorResponseWithData(ctx *fiber.Ctx, status int, err error, data interface{}) error {
	return ctx.Status(status).JSON(&Response{
		Status:  "fail",
		Message: err.Error(),
		Data:    data,
	})
}

// SuccessResponse 返回成功响应
func SuccessResponse(ctx *fiber.Ctx, status int, message string, data interface{}) error {
	return ctx.Status(status).JSON(&Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// NoContentResponse 返回无内容响应
func NoContentResponse(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNoContent)
}
