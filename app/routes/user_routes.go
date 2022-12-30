package routes

import (
	"github.com/bearts/go-fiber/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	user := app.Group("/v1/user")
	controllers.UserPanel(user)
	controllers.UserOrder(user)
}
