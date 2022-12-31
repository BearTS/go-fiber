package middleware

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func FiberMiddleware(a *fiber.App) {
	file, err := os.OpenFile("./combined.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	a.Use(
		// Add CORS to each route.
		cors.New(),
		recover.New(),
		// Add simple logger.
		logger.New(logger.Config{
			Format:     "${pid} ${locals:requestid} ${status} - ${method} ${path} - ${latency} - ${ip} - ${ua} - ${error}\n",
			TimeFormat: "02-Jan-2006",
			TimeZone:   "Asia/Calcutta",
			Output:     file,
		}),
	)
	JWTProtected()
}
