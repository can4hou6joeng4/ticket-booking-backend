package main

import (
	"fmt"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/config"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/db"
	_ "github.com/can4hou6joeng4/ticket-booking-project-v1/docs" // swagger docs
	"github.com/can4hou6joeng4/ticket-booking-project-v1/handlers"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/middlewares"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/repositories"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title           Ticket Booking API
// @version         1.0
// @description     A ticket booking system API server.

// @contact.name   API Support
// @contact.url    https://github.com/can4hou6joeng4
// @contact.email  can4hou6joeng4@163.com

// @host      localhost:8081
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// 以下是自定义的 Swagger 配置，可以通过修改这些注释来更新文档信息
// 注意：修改这些注释后需要重新运行 swag init 命令来生成新的文档

// @info.title           Ticket Booking API
// @info.version         1.0
// @info.description     A ticket booking system API server.
// @info.contact.url     https://github.com/can4hou6joeng4
// @info.contact.email   can4hou6joeng4@163.com

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "TickBooking",
		ServerHeader: "Fiber",
	})

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Config
	envConfig := config.NewEnvConfig()
	redis := db.InitRedis(envConfig)
	db := db.InitDatabase(envConfig, db.DBMigrator)

	// Repository
	eventRepository := repositories.NewEventRepository(db)
	ticketRepository := repositories.NewTicketRepository(db)
	authRepository := repositories.NewAuthRepository(db)
	statisticsRepository := repositories.NewStatisticsRepository(db)
	// Service
	authService := services.NewAuthService(authRepository, redis)
	// Routing
	server := app.Group("/api")
	handlers.NewAuthHandler(server.Group("/auth"), authService)

	privateRoutes := server.Use(middlewares.AuthProtected(db, redis))
	handlers.NewAuthProtectedHandler(privateRoutes.Group("/auth"), authService)

	handlers.NewEventHandler(privateRoutes.Group("/event"), eventRepository)
	handlers.NewTicketHandler(privateRoutes.Group("/ticket"), ticketRepository, eventRepository, envConfig, redis)
	handlers.NewStatisticsHandler(privateRoutes.Group("/statistics"), statisticsRepository)

	app.Listen(fmt.Sprintf(":%s", envConfig.ServerPort))
}
