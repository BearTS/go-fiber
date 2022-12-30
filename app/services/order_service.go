package services

import (
	"github.com/bearts/go-fiber/app/dao"
	"github.com/bearts/go-fiber/app/models"
	"github.com/bearts/go-fiber/app/structs"
	"github.com/bearts/go-fiber/app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	price, err := dao.GetPriceFromTo("main_gate", body.Location)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	locationObj, err := dao.GetPlaceByCode(body.Location)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	// convert id to object id
	userid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	// create order object
	var order models.Order
	order.Delivery_app.NameOfApp = body.NameOfApp
	order.Delivery_app.NameOfRes = body.NameOfRestaurant
	order.Delivery_app.EstimatedTime = body.EstimatedTime
	order.Delivery_app.DeliveryPhone = body.DeliveryPhone
	order.Location = locationObj.Id
	if body.Otp > 0 {
		order.Delivery_app.Otp = body.Otp
	}
	order.Price = price
	order.Status = "pending"
	order.User = userid

	// create order to database
	Order, err := dao.CreateOrder(&order)

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

func GetAllOrdersByUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	id := user.Claims.(jwt.MapClaims)["id"].(string)
	status := c.Query("status")
	// convert id to object id
	userid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	orders, err := dao.GetAllOrdersOfUser(userid, status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"orders":  orders,
	})
}
