package controllers

import (
	"github.com/bearts/go-fiber/app/services"
	"github.com/bearts/go-fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func RunnerOrder(app fiber.Router) {
	app.Put("/order/update/:id", middleware.RunnerOnly, services.RunnerAssignOrderById)
	app.Put("/order/update/:id/complete", middleware.RunnerOnly, services.RunnerDeliverOrder)
	app.Put("/order/update/:id/status", middleware.RunnerOnly, services.RunnerChangeOrderStatus)
	app.Get("/order/unassigned", middleware.RunnerOnly, services.RunnerGetAllUnassignedOrders)
	app.Get("/order/current", middleware.RunnerOnly, services.RunnerGetAllCurrentOrders)
	app.Get("/order/previous", middleware.RunnerOnly, services.RunnerGetAllPreviousOrders)
}
