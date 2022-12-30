package routes

import (
	"github.com/bearts/go-fiber/app/controllers"
	"github.com/bearts/go-fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserPanel(app *fiber.App) {
	app.Post("/user/otp", controllers.SendOtp)
	app.Post("/user/otp/verify", controllers.VerifyOtp)
	app.Get("/user/current", middleware.UserOnly, controllers.CurrentUser)
}
