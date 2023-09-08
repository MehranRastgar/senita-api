package Controllers

import (
	"fmt"
	"senita-api/db"
	"senita-api/models"
	"senita-api/pkg/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AuthController represents the authentication controller
type AuthController struct {
	SecretKey string // JWT secret key
	DB        *gorm.DB
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(secretKey string) *AuthController {
	return &AuthController{
		SecretKey: secretKey,
	}
}

// Login handles user authentication and generates a JWT token
func (ac *AuthController) Login(ctx *fiber.Ctx) error {
	// Parse the request body to get the username and password
	loginRequest := struct {
		Username string `json:"user_name"`
		Password string `json:"password"`
	}{}
	if err := ctx.BodyParser(&loginRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Query the database to find the user by username
	var user models.User
	if err := db.DB.Where("user_name = ?", loginRequest.Username).First(&user).Error; err != nil {
		// User not found or error occurred
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf(`Invalid username %s`, loginRequest.Username)})
	}

	// Verify the password (you should use a secure password hashing library)
	if !utils.TestPassword([]byte(loginRequest.Password), user.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf(`Invalid username or password %s`, loginRequest.Username)})
	}

	// Generate a JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.UserName
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time (e.g., 24 hours)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(ac.SecretKey))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Return the token in the response
	return ctx.JSON(fiber.Map{"token": tokenString})
}

// VerifyToken verifies a JWT token
func (ac *AuthController) VerifyToken(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization")

	// Verify the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(ac.SecretKey), nil
	})

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(int)       // Extract user ID
		username := claims["username"].(string) // Extract username

		// You can use the user ID or username to perform further actions or authentication checks.

		return ctx.JSON(fiber.Map{"message": "Token is valid", "user_id": userID, "username": username})
	}

	return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
}

// RegisterRoutes registers the article routes.
func (ac *AuthController) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	auth := api.Group("/auth")

	auth.Post("/login", ac.Login)
	// Secure route that requires JWT authentication
	auth.Get("/secure", ac.VerifyToken, func(ctx *fiber.Ctx) error {
		return ctx.SendString("Secure route accessed successfully")
	})
}
