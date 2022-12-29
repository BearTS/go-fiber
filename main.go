package main

import (
	"log"

	"github.com/bearts/go-fiber/app/routes"
	"github.com/bearts/go-fiber/pkg/configs"
	"github.com/bearts/go-fiber/pkg/middleware"
	"github.com/bearts/go-fiber/pkg/utils"
	"github.com/bearts/go-fiber/platform/database"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Define Fiber Config
	config := configs.FiberConfig()
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	// Define new app with Fiber config
	app := fiber.New(config)
	// use middleware
	middleware.FiberMiddleware(app)

	database.PostgreSQLConnection()

	// add routes
	routes.UserRoute(app)

	utils.StartServer(app)
}
