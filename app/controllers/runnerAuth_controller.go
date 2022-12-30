package controllers

import (
	"github.com/bearts/go-fiber/app/services"
	"github.com/gofiber/fiber/v2"
)

func RunnerAuth(app fiber.Router) {
	app.Post("/auth/signin", services.RunnerSignIn)
	app.Post("/auth/signup", services.RunnerSignUp)
}
