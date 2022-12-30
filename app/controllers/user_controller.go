package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
	"github.com/bearts/go-fiber/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")
var otpCollection *mongo.Collection = database.GetCollection(database.DB, "otp")

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type verifyOtp_Body struct {
	Email string `json:"email" validate:"required,email"`
	Otp   int    `json:"otp" validate:"required,min=4,number"`
}

type sendOtp_Body struct {
	Email string `json:"email" validate:"required,email"`
}

var validate = validator.New()

func VerifyOtp(c *fiber.Ctx) error {
	var body verifyOtp_Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	// validate body
	if err := validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	// check if email already exists
	var existingUser models.User
	if err := userCollection.FindOne(c.Context(), bson.M{"email": body.Email}).Decode(&existingUser); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	var otpModel models.Otp
	// check about expiry time
	if err := otpCollection.FindOne(c.Context(), bson.M{"user": existingUser.Id, "otp": body.Otp}).Decode(&otpModel); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Otp is invalid",
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
	if _, err := otpCollection.DeleteOne(c.Context(), bson.M{"user": existingUser.Id, "otp": body.Otp}); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}
	// generate jwt token
	token, err := services.CreateJWTTokenUser(existingUser)
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
	var body sendOtp_Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if err := validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	// check if email already exists
	var existingUser models.User
	userCollection.FindOne(c.Context(), bson.M{"email": body.Email}).Decode(&existingUser)
	if existingUser.Email != "" {
		data, err := genOtpAndSendOtp(c, body.Email, existingUser.Id)
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
	if _, err := userCollection.InsertOne(c.Context(), newUser); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}
	data, err := genOtpAndSendOtp(c, body.Email, newUser.Id)
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
	var currentUser models.User
	if err := userCollection.FindOne(c.Context(), models.User{Email: email}).Decode(&currentUser); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"user":    currentUser,
	})
}

// genOtpAndSendOtp generates otp and sends it to the user
func genOtpAndSendOtp(c *fiber.Ctx, email string, id primitive.ObjectID) (interface{}, error) {
	var otp models.Otp
	otp_number := services.GenerateOtp()
	// find if otp already exists for the user then delete it
	otpCollection.FindOneAndDelete(c.Context(), bson.M{"user": id})
	otp = models.Otp{
		Id:        primitive.NewObjectID(),
		Otp:       otp_number,
		User:      id,
		ExpiresAt: time.Now().Add(time.Minute * 5).UnixMilli(),
	}

	if _, err := otpCollection.InsertOne(c.Context(), otp); err != nil {
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
