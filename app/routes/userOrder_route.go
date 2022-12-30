package routes

import (
	"github.com/bearts/go-fiber/app/controllers"
	"github.com/bearts/go-fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserOrder(app *fiber.App) {
	app.Post("/order/create", middleware.UserOnly, controllers.CreateOrder)
}
