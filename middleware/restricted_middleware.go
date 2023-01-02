package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func UserOnly(c *fiber.Ctx) error {
	_token := c.Request().Header.Peek("Authorization")
	if _token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Missing or malformed JWT",
		})
	}
	token := _token[7:]
	// Parse the token and store it in the "user" key of the Locals map
	user, err := jwt.Parse(string(token), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	c.Locals("user", user)
	// check if user.role is user
	if user.Claims.(jwt.MapClaims)["role"] != "user" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	}
	return c.Next()
}

func RunnerOnly(c *fiber.Ctx) error {
	_token := c.Request().Header.Peek("Authorization")
	if _token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Missing or malformed JWT",
		})
	}
	token := _token[7:]
	// Parse the token and store it in the "runner" key of the Locals map
	runner, err := jwt.Parse(string(token), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	c.Locals("runner", runner)
	// check if user.role is user
	if runner.Claims.(jwt.MapClaims)["role"] != "runner" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	}
	return c.Next()
}
