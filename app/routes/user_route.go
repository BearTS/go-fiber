package routes

import (
	"github.com/bearts/go-fiber/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	app.Post("/user", controllers.SignUp)
	app.Get("/user", controllers.GetUsers)
	app.Post("/user/login", controllers.Login)
}
