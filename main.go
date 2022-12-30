package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/routes"
	"github.com/bearts/go-fiber/middleware"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	// Define Fiber Config
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	config := fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
		AppName:     "Tez Backend",
	}
	// Define new app with Fiber config
	app := fiber.New(config)
	// use middleware
	middleware.FiberMiddleware(app)

	database.MongoConnectDB()

	// add routes
	routes.UserRoutes(app)
	routes.RunnerRoutes(app)

	if err := app.Listen(":" + os.Getenv("SERVER_PORT")); err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

}
