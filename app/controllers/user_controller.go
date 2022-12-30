package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bearts/go-fiber/app/dao"
	"github.com/bearts/go-fiber/app/interfaces"
	"github.com/bearts/go-fiber/app/models"
	"github.com/bearts/go-fiber/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func VerifyOtp(c *fiber.Ctx) error {
	var body interfaces.Body_VerifyOtp
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	// validate body
	if err := services.Validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	// check if email already exists
	existingUser, err := dao.FindUserByEmail(body.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}
	// check about expiry time
	otpModel, err := dao.FindOtpByUser(*existingUser)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Otp not found",
		})
	}
	// check if otp is expired; expiresAt is a UnixMilli timestamp
	if otpModel.ExpiresAt < time.Now().UnixMilli() {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Otp is expired",
		})
	}
	// delete otp from db
	if _, err := dao.DeleteOtpByUser(*existingUser); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}
	// generate jwt token
	token, err := services.CreateJWTTokenUser(*existingUser)
	if err != nil {
		fmt.Println(err)
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

func SendOtp(c *fiber.Ctx) error {
	var body interfaces.Body_SendOtp
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if err := services.Validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	// check if email already exists
	existingUser, err := dao.FindUserByEmail(body.Email)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Existing user: ", existingUser)
	if existingUser != nil {
		data, err := genOtpAndSendOtp(c, body.Email, *existingUser)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.Status(200).JSON(data)
	}

	newUser := models.User{
		Id:    primitive.NewObjectID(),
		Email: body.Email,
	}
	if _, err := dao.CreateUser(newUser); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}
	data, err := genOtpAndSendOtp(c, body.Email, newUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}
	return c.Status(200).JSON(data)
}

func CurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	currentUser, err := dao.FindUserByEmail(email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"user":    currentUser,
	})
}

// genOtpAndSendOtp generates otp and sends it to the user
func genOtpAndSendOtp(c *fiber.Ctx, email string, user models.User) (interface{}, error) {
	var otp models.Otp
	otp_number := services.GenerateOtp()
	// find if otp already exists for the user then delete it
	if _, err := dao.DeleteOtpByUser(user); err != nil {
		return nil, err
	}
	otp = models.Otp{
		Id:        primitive.NewObjectID(),
		Otp:       otp_number,
		User:      user.Id,
		ExpiresAt: time.Now().Add(time.Minute * 5).UnixMilli(),
	}

	if _, err := dao.CreateOtp(otp); err != nil {
		return nil, err
	}
	subject := "OTP for login into your account"
	body := "Your OTP is " + strconv.Itoa(otp_number)
	go func() {
		if err := services.SendEmail(email, subject, body); err != nil {
			fmt.Println(err)
		}
	}()
	response := map[string]interface{}{
		"success": true,
		"message": "OTP sent to your email",
	}
	return response, nil
}
