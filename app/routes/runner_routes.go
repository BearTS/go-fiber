package routes

import (
	"github.com/bearts/go-fiber/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func RunnerRoutes(app *fiber.App) {
	runner := app.Group("/v1/runner")
	controllers.RunnerAuth(runner)
	controllers.RunnerOrder(runner)
}
