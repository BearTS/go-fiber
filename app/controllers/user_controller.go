package controllers

import (
	"fmt"
	"strconv"

	"github.com/bearts/go-fiber/app/models"
	"github.com/bearts/go-fiber/app/services"
	"github.com/bearts/go-fiber/platform/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")
var otpCollection *mongo.Collection = database.GetCollection(database.DB, "otp")
var validate = validator.New()

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type verifyBody struct {
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

func VerifyOtp(c *fiber.Ctx) error {
	var body verifyBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	email := body.Email
	otp := body.Otp
	if email == "" {
		return c.Status(400).JSON("Email is required")
	}
	if len(strconv.Itoa(otp)) != 4 {
		return c.Status(400).JSON("Otp is required")
	}
	// check if email already exists
	var existingUser models.User
	err := userCollection.FindOne(c.Context(), bson.M{"email": email}).Decode(&existingUser)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}
	var otpModel models.Otp
	// find in db where user.id and otp match
	err = otpCollection.FindOne(c.Context(), bson.M{"user": existingUser.Id, "otp": otp}).Decode(&otpModel)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Otp is invalid",
		})
	}
	// delete otp from db
	_, err = otpCollection.DeleteOne(c.Context(), bson.M{"user": existingUser.Id, "otp": otp})
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	// generate jwt token
	token, err := services.CreateJWTTokenUser(existingUser)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"token":   token,
	})
}

func SendOtp(c *fiber.Ctx) error {
	var user verifyBody
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	email := user.Email
	if email == "" {
		return c.Status(400).JSON("Email is required")
	}

	// check if email already exists
	var existingUser models.User
	userCollection.FindOne(c.Context(), bson.M{"email": email}).Decode(&existingUser)
	if existingUser.Email != "" {
		data, err := genOtpAndSendOtp(c, email, existingUser.Id)
		if err != nil {
			return c.Status(400).JSON(err.Error())
		}
		return c.Status(200).JSON(data)
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(400).JSON(validationErr.Error())
	}
	newUser := models.User{
		Id:    primitive.NewObjectID(),
		Email: user.Email,
	}
	_, err := userCollection.InsertOne(c.Context(), newUser)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	data, err := genOtpAndSendOtp(c, email, newUser.Id)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}

func CurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	var currentUser models.User
	err := userCollection.FindOne(c.Context(), models.User{Email: email}).Decode(&currentUser)
	if err != nil {
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
	subject := "OTP for login into your account"
	body := "Your OTP is " + strconv.Itoa(otp_number)
	// find if otp already exists for the user then delete it
	otpCollection.FindOneAndDelete(c.Context(), bson.M{"user": id})
	otp = models.Otp{
		Id:   primitive.NewObjectID(),
		Otp:  otp_number,
		User: id,
	}
	_, err := otpCollection.InsertOne(c.Context(), otp)
	if err != nil {
		return nil, err
	}
	go func() {
		err := services.SendEmail(email, subject, body)
		if err != nil {
			fmt.Println(err)
		}
	}()
	response := map[string]interface{}{
		"success": true,
		"message": "OTP sent to your email",
	}
	return response, nil
}
