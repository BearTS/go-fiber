package routes

import (
	"github.com/bearts/go-fiber/app/controllers"
	"github.com/bearts/go-fiber/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	app.Post("/user", controllers.SignUp)
	app.Post("/user/login", controllers.Login)
	app.Get("/user", middleware.UserOnly, controllers.GetUsers)
	app.Get("/user/current", middleware.UserOnly, controllers.CurrentUser)
}
