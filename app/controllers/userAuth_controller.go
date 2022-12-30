package controllers

import (
	"github.com/bearts/go-fiber/app/services"
	"github.com/bearts/go-fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserAuth(app fiber.Router) {
	app.Post("/auth/otp", services.UserSendOtp)
	app.Post("/auth/otp/verify", services.UserVerifyOtp)
	app.Get("/profile", middleware.UserOnly, services.UserCurrent)
	app.Put("/profile", middleware.UserOnly, services.UserUpdate)
}
