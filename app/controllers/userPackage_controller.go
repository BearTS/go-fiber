package controllers

import (
	"github.com/bearts/go-fiber/app/services"
	"github.com/bearts/go-fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserPackage(app fiber.Router) {
	app.Post("/package", middleware.UserOnly, services.UserCreatePackage)
	app.Get("/package", middleware.UserOnly, services.UserGetAllPackage)
	app.Get("/package/:id", middleware.UserOnly, services.UserGetPackageById)
	app.Put("/package/:id", middleware.UserOnly, services.UserUpdatePackageStatus)
}
