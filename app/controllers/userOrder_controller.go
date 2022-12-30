package controllers

import (
	"github.com/bearts/go-fiber/app/services"
	"github.com/bearts/go-fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserOrder(app fiber.Router) {
	app.Post("/order", middleware.UserOnly, services.UserCreateOrder)
	app.Get("/order", middleware.UserOnly, services.UserGetAllOrdersByUser)
	app.Get("/order/:id", middleware.UserOnly, services.UserGetOrderById)
}
