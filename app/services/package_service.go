package services

import (
	"github.com/bearts/go-fiber/app/dao"
	"github.com/bearts/go-fiber/app/structs"
	"github.com/bearts/go-fiber/app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func UserCreatePackage(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	// id := user.Claims.(jwt.MapClaims)["id"].(string)
	claims := user.Claims.(jwt.MapClaims)
	PhoneAvailable := claims["PhoneAvailable"].(bool)
	if !PhoneAvailable {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Phone number is required",
		})
	}
	var body structs.UserCreatePackage
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := utils.Validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	delivery, err := dao.GetPlaceByCode(body.DeliveryLocation)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	pickup, err := dao.GetPlaceByCode(body.Package.Location)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	price, err := dao.GetPriceFromToById(pickup.Id, delivery.Id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    price,
	})
	// ! TO DO
}
