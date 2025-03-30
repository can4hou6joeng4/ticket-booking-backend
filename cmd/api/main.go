package main

import (
	"fmt"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/config"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/db"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/handlers"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/middlewares"
	repositorys "github.com/can4hou6joeng4/ticket-booking-project-v1/repositories"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "TickBooking",
		ServerHeader: "Fiber",
	})
	// Config
	envConfig := config.NewEnvConfig()
	db := db.Init(envConfig, db.DBMigrator)
	// Repository
	eventRepository := repositorys.NewEventRepository(db)
	ticketRepository := repositorys.NewTicketRepository(db)
	authRepository := repositorys.NewAuthRepository(db)
	// Service
	authService := services.NewAuthService(authRepository)
	// Routing
	server := app.Group("/api")
	handlers.NewAuthHandler(server.Group("/auth"), authService)

	privateRoutes := server.Use(middlewares.AuthProtected(db))

	handlers.NewEventHandler(privateRoutes.Group("/event"), eventRepository)
	handlers.NewTicketHandler(privateRoutes.Group("/ticket"), ticketRepository)

	app.Listen(fmt.Sprintf(":" + envConfig.ServerPort))
}
