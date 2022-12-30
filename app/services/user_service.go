package services

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/bearts/go-fiber/app/dao"
	"github.com/bearts/go-fiber/app/models"
	"github.com/bearts/go-fiber/app/structs"
	"github.com/bearts/go-fiber/app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func VerifyOtp(c *fiber.Ctx) error {
	var body structs.BodyVerifyOtp
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
	existingUser, err := dao.GetUserByEmail(body.Email)
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
	token, err := utils.CreateJWTTokenUser(*existingUser)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}
	PhoneAvailable := false
	if *existingUser.Phone != "" {
		PhoneAvailable = true
	}
	return c.Status(200).JSON(fiber.Map{
		"success":        true,
		"token":          token,
		"PhoneAvailable": PhoneAvailable,
	})
}

func SendOtp(c *fiber.Ctx) error {
	var body structs.BodySendOtp
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
	existingUser, err := dao.GetUserByEmail(body.Email)
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
	currentUser, err := dao.GetUserByEmail(email)
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

func UpdateUser(c *fiber.Ctx) error {
	var body structs.BodyUpdateUser
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
	re := regexp.MustCompile(`^[0-9]{2}[A-Za-z]{3}[0-9]{4}$`)
	if !re.MatchString(body.RegistrationNumber) {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid registration number",
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
	user, err := dao.GetUserById(userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}
	if body.RegistrationNumber != "" {
		user.RegistrationNumber = &body.RegistrationNumber
	}
	if body.Phone != "" {
		user.Phone = &body.Phone
	}
	if body.Name != "" {
		user.Name = body.Name
	}
	if body.DefaultAddress != "" {
		location, err := dao.GetPlaceByCode(body.DefaultAddress)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid location",
			})
		}
		user.DefaultAddress = &location.Id
	}
	if _, err = dao.UpdateUser(*user); err != nil {
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
