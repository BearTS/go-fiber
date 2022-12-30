package services

import (
	"github.com/bearts/go-fiber/app/dao"
	"github.com/bearts/go-fiber/app/models"
	"github.com/bearts/go-fiber/app/structs"
	"github.com/bearts/go-fiber/app/utils"
	"github.com/gofiber/fiber/v2"
)

func RunnerSignUp(c *fiber.Ctx) error {
	var body structs.RunnerSignUp
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
	// check if email already exists
	_, err := dao.GetRunnerByEmail(body.Email)
	if err == nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Email already exists",
		})
	}
	// hash password
	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}
	// create runner
	runner := models.Runner{
		Name:     body.Name,
		Email:    body.Email,
		Password: hashedPassword,
		Phone:    body.Phone,
	}
	// save runner
	runnerobj, err := dao.CreateRunner(runner)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}
	// generate token
	token, err := utils.CreateJWTTokenRunner(runner)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"token":   token,
		"runner":  runnerobj,
	})
}

func RunnerSignIn(c *fiber.Ctx) error {
	var body structs.RunnerSignIn
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
	runner, err := dao.GetRunnerByEmail(body.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Runner not found",
		})
	}
	if !utils.CheckPasswordHash(body.Password, runner.Password) {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid password",
		})
	}
	// generate token
	token, err := utils.CreateJWTTokenRunner(*runner)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"token":   token,
	})
}
