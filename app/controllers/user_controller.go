package controllers

import (
	"time"

	"github.com/bearts/go-fiber/app/models"
	"github.com/bearts/go-fiber/app/services"
	"github.com/bearts/go-fiber/platform/database"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"` // this is not returned to the user
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateResponseUser(user models.User) User {
	return User{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func SignUp(c *fiber.Ctx) error {
	var user models.User
	err := c.BodyParser(&user)
	if err != nil {
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
	var count int64
	database.Database.Db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return c.Status(400).JSON("User already exists")
	}
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	newUser := models.User{
		Email:    user.Email,
		Password: string(bs),
	}
	database.Database.Db.Create(&newUser)
	responseUser := CreateResponseUser(newUser)
	return c.Status(200).JSON(responseUser)
}

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}
	database.Database.Db.Find(&users)
	responseUsers := []User{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

func CurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	var currentUser models.User
	database.Database.Db.Where("email = ?", email).First(&currentUser)
	responseUser := CreateResponseUser(currentUser)
	return c.Status(200).JSON(responseUser)
}

func Login(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	email := user.Email
	password := user.Password

	if err := database.Database.Db.Where("email = ?", email).First(&user).Error; err != nil {
		// user does not exist, throw error
		// you can write the error message to the response writer
		return c.Status(400).JSON("Invalid credentials")
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
