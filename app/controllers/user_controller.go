package controllers

import (
	"fmt"

	"github.com/bearts/go-fiber/app/models"
	"github.com/bearts/go-fiber/app/services"
	"github.com/bearts/go-fiber/platform/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")
var validate = validator.New()

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"` // this is not returned to the user
}

func CreateResponseUser(user models.User) User {
	return User{
		ID:    user.Id.String(),
		Email: user.Email,
	}
}

func SignUp(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	email := user.Email
	password := user.Password
	// password is byte array
	if email == "" {
		return c.Status(400).JSON("Email is required")
	}
	if len(password) < 6 {
		return c.Status(400).JSON("Password must be atleast 6 characters")
	}

	// check if email already exists
	var existingUser models.User
	userCollection.FindOne(c.Context(), bson.M{"email": email}).Decode(&existingUser)
	if existingUser.Email != "" {
		return c.Status(400).JSON("Email already exists")
	}
	fmt.Println("user", existingUser)
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(400).JSON(validationErr.Error())
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Email:    user.Email,
		Password: string(bs),
	}
	_, err = userCollection.InsertOne(c.Context(), newUser)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	response := CreateResponseUser(newUser)
	return c.Status(200).JSON(response)
}

func GetUsers(c *fiber.Ctx) error {
	// get all users
	var users []models.User
	cursor, err := userCollection.Find(c.Context(), models.User{})
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if err = cursor.All(c.Context(), &users); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	var responseUsers []User
	for _, user := range users {
		responseUsers = append(responseUsers, CreateResponseUser(user))
	}
	return c.Status(200).JSON(responseUsers)
}

func CurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	var currentUser models.User
	err := userCollection.FindOne(c.Context(), models.User{Email: email}).Decode(&currentUser)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	response := CreateResponseUser(currentUser)
	return c.Status(200).JSON(response)
}

func Login(c *fiber.Ctx) error {

	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	email := user.Email
	password := user.Password
	if email == "" {
		return c.Status(400).JSON("Email is required")
	}
	if len(password) < 6 {
		return c.Status(400).JSON("Password must be atleast 6 characters")
	}
	err := userCollection.FindOne(c.Context(), models.User{Email: email}).Decode(&user)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// password is incorrect, throw error
		// you can write the error message to the response writer
		return c.Status(400).JSON("Invalid credentials")
	}
	token, err := services.CreateJWTTokenUser(user)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	response := map[string]interface{}{
		"success": true,
		"token":   token,
	}
	return c.Status(200).JSON(response)

}
