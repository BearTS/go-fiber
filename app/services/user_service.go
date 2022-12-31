package services

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/bearts/go-fiber/app/dbFunctions"
	"github.com/bearts/go-fiber/app/models"
	"github.com/bearts/go-fiber/app/structs"
	"github.com/bearts/go-fiber/app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UserVerifyOtp(c *fiber.Ctx) error {
	var body structs.UserVerifyOtp
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
	// check if email already exists
	existingUser, err := dbFunctions.GetUserByEmail(body.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}
	// check about expiry time
	otpModel, err := dbFunctions.FindOtpByUser(*existingUser)
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
	if _, err := dbFunctions.DeleteOtpByUser(*existingUser); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}
	PhoneAvailable := false
	if existingUser.Phone != "" {
		PhoneAvailable = true
	}
	// generate jwt token
	token, err := utils.CreateJWTTokenUser(*existingUser, PhoneAvailable)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success":        true,
		"token":          token,
		"PhoneAvailable": PhoneAvailable,
	})
}

func UserSendOtp(c *fiber.Ctx) error {
	var body structs.UserSendOtp
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if err := utils.Validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	// check if email already exists
	existingUser, err := dbFunctions.GetUserByEmail(body.Email)
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
		Phone: "",
	}
	if _, err := dbFunctions.CreateUser(newUser); err != nil {
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

func UserCurrent(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	currentUser, err := dbFunctions.GetUserByEmail(email)
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

func UserUpdate(c *fiber.Ctx) error {
	var body structs.UserUpdateUser
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if err := utils.Validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	currentuser := c.Locals("user").(*jwt.Token)
	claims := currentuser.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user id",
		})
	}
	user, err := dbFunctions.GetUserById(userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}
	if body.RegistrationNumber != "" {
		re := regexp.MustCompile(`^[0-9]{2}[A-Za-z]{3}[0-9]{4}$`)
		if !re.MatchString(body.RegistrationNumber) {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid registration number",
			})
		}
		user.RegistrationNumber = &body.RegistrationNumber
	}
	if body.Phone != "" {
		user.Phone = body.Phone
	}
	if body.Name != "" {
		user.Name = body.Name
	}
	if body.DefaultAddress != "" {
		location, err := dbFunctions.GetPlaceByCode(body.DefaultAddress)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid location",
			})
		}
		user.DefaultAddress = &location.Id
	}
	if _, err = dbFunctions.UpdateUser(*user); err != nil {
		fmt.Println("Error: ", err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

// genOtpAndSendOtp generates otp and sends it to the user
func genOtpAndSendOtp(c *fiber.Ctx, email string, user models.User) (interface{}, error) {
	var otp models.Otp
	otp_number := utils.GenerateOtp()
	// find if otp already exists for the user then delete it
	if _, err := dbFunctions.DeleteOtpByUser(user); err != nil {
		return nil, err
	}
	otp = models.Otp{
		Id:        primitive.NewObjectID(),
		Otp:       otp_number,
		User:      user.Id,
		ExpiresAt: time.Now().Add(time.Minute * 5).UnixMilli(),
	}

	if _, err := dbFunctions.CreateOtp(otp); err != nil {
		return nil, err
	}
	subject := "OTP for login into your account"
	body := "Your OTP is " + strconv.Itoa(otp_number)
	go func() {
		if err := utils.SendEmail(email, subject, body); err != nil {
			fmt.Println(err)
		}
	}()
	response := map[string]interface{}{
		"success": true,
		"message": "OTP sent to your email",
	}
	return response, nil
}
