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

func UserCreateOrder(c *fiber.Ctx) error {
	var body structs.UserCreateOrder
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
	RunnerOtp := utils.GenerateOtp()
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
	order.RunnerOtp = RunnerOtp
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

func UserGetAllOrdersByUser(c *fiber.Ctx) error {
	status := c.Query("status")
	user := c.Locals("user").(*jwt.Token)
	id := user.Claims.(jwt.MapClaims)["id"].(string)
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

func UserGetOrderById(c *fiber.Ctx) error {
	id := c.Params("id")
	// add validator to check id
	if err := utils.Validate.Var(id, "required,len=24"); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	user := c.Locals("user").(*jwt.Token)
	_id := user.Claims.(jwt.MapClaims)["id"].(string)
	userid, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	// convert id to object id
	orderid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	order, err := dao.GetOrderById(orderid, userid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	if order.Runner != nil {
		runner, err := dao.GetRunnerById(*order.Runner)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   "Internal server error",
				"message": err.Error(),
			})
		}
		runner.Password = ""
		runner.Email = ""
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"order":   order,
			"runner":  runner,
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"order":   order,
	})

}

// Runner Panel
func RunnerAssignOrderById(c *fiber.Ctx) error {
	id := c.Params("id")
	// add validator to check id
	if err := utils.Validate.Var(id, "required,len=24"); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	runner := c.Locals("runner").(*jwt.Token)
	_id := runner.Claims.(jwt.MapClaims)["id"].(string)
	runnerid, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	// convert id to object id
	orderid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	order, err := dao.AssignOrderById(orderid, runnerid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"order":   order,
	})
}

func RunnerDeliverOrder(c *fiber.Ctx) error {
	var body structs.RunnerDeliverOrder
	id := c.Params("id")
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	// validate body
	if err := utils.Validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	// add validator to check id
	if err := utils.Validate.Var(id, "required,len=24"); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	orderid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	runner := c.Locals("runner").(*jwt.Token)
	_id := runner.Claims.(jwt.MapClaims)["id"].(string)
	runnerid, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	orderObj, err := dao.GetOrderById(orderid, primitive.NilObjectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	if orderObj.RunnerOtp != body.Otp {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid OTP",
		})
	}
	order, err := dao.UpdateOrderStatus(orderid, runnerid, "delivered")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"order":   order,
	})
}

func RunnerChangeOrderStatus(c *fiber.Ctx) error {
	orderid := c.Params("id")
	if err := utils.Validate.Var(orderid, "required,len=24"); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	var body structs.RunnerChangeStatus
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	// validate body
	if err := utils.Validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	runner := c.Locals("runner").(*jwt.Token)
	_id := runner.Claims.(jwt.MapClaims)["id"].(string)
	runnerid, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	// convert id to object id
	id, err := primitive.ObjectIDFromHex(orderid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	order, err := dao.UpdateOrderStatus(id, runnerid, body.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"order":   order,
	})
}

func RunnerGetAllUnassignedOrders(c *fiber.Ctx) error {
	orders, err := dao.GetAllUnassignedOrders()
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

func RunnerGetAllPreviousOrders(c *fiber.Ctx) error {
	runner := c.Locals("runner").(*jwt.Token)
	_id := runner.Claims.(jwt.MapClaims)["id"].(string)
	runnerid, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	orders, err := dao.GetAllPreviousOrders(runnerid)
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

func RunnerGetAllCurrentOrders(c *fiber.Ctx) error {
	runner := c.Locals("runner").(*jwt.Token)
	_id := runner.Claims.(jwt.MapClaims)["id"].(string)
	runnerid, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	orders, err := dao.GetAllCurrentOrders(runnerid)
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
