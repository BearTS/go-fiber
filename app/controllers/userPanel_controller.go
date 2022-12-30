package controllers

import (
	"github.com/bearts/go-fiber/app/services"
	"github.com/bearts/go-fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserPanel(app fiber.Router) {
	app.Post("/auth/otp", services.SendOtp)
	app.Post("/auth/otp/verify", services.VerifyOtp)
	app.Get("/profile", middleware.UserOnly, services.CurrentUser)
	app.Put("/profile", middleware.UserOnly, services.UpdateUser)
}
