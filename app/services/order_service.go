package services

import (
	"github.com/bearts/go-fiber/app/dao"
	"github.com/bearts/go-fiber/app/structs"
	"github.com/bearts/go-fiber/app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func CreateOrder(c *fiber.Ctx) error {
	var body structs.BodyCreateOrder
	user := c.Locals("user").(*jwt.Token)
	id := user.Claims.(jwt.MapClaims)["id"].(string)
	claims := user.Claims.(jwt.MapClaims)
	// parse body
	PhoneAvailable := claims["PhoneAvailable"].(bool)
	if !PhoneAvailable {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Phone number is required",
		})
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	// validate body
	if err := utils.Validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}

	Order, err := dao.CreateOrder(id, body.NameOfApp, body.NameOfRestaurant, body.EstimatedTime, body.DeliveryPhone, body.Location, body.Otp)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Create order error",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Order created",
		"order":   Order,
	})
}
