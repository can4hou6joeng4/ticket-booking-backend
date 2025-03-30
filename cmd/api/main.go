package main

import (
	"fmt"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/config"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/db"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/handlers"
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
	// TODO add privateRoutes
	// privateRoutes := server.Use(middlewares.AuthProtected(db))
	// Handlers
	handlers.NewEventHandler(server.Group("/event"), eventRepository)
	handlers.NewTicketHandler(server.Group("/ticket"), ticketRepository)
	handlers.NewAuthHandler(server.Group("/auth"), authService)

	app.Listen(fmt.Sprintf(":" + envConfig.ServerPort))
}
