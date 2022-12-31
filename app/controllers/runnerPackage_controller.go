package controllers

import (
	"github.com/bearts/go-fiber/app/services"
	"github.com/bearts/go-fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func RunnerPackage(app fiber.Router) {
	app.Get("/package/unassigned", middleware.RunnerOnly, services.RunnerGetAllUnAssignedPackage)
	app.Get("/package/previous", middleware.RunnerOnly, services.RunnerGetAllPreviousPackage)
	app.Get("/package/id/:id", middleware.RunnerOnly, services.RunnerGetPackageById)

	app.Put("/package/assign/:id", middleware.RunnerOnly, services.RunnerAssignPackage)
	app.Put("/package/update/:id", middleware.RunnerOnly, services.RunnerUpdatePackage)
	app.Put("/package/complete/:id", middleware.RunnerOnly, services.RunnerDeliverPackage)
}
