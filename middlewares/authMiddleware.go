package Middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthController struct {
	SecretKey string // JWT secret key
	DB        *gorm.DB
}

func ValidateToken() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Get the token from the "Authorization" header with the "Bearer" prefix
		authorizationHeader := ctx.Get("Authorization")
		if authorizationHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is missing"})
		}

		// Split the header value to get the token part
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		tokenString := headerParts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Check if the token is valid and not expired
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
			if time.Now().After(expirationTime) {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token has expired"})
			}

			// Store the user ID and username from the claims for later use if needed
			ctx.Locals("user_id", int(claims["user_id"].(float64)))
			ctx.Locals("username", claims["username"])

			// Continue to the next middleware or route handler
			return ctx.Next()
		}

		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}
}

// func AuthMiddleware() fiber.Handler {
// 	return func(ctx *fiber.Ctx) error {
// 		// Get the token from the Authorization header
// 		tokenString := ctx.Get("Authorization")

// 		if tokenString == "" {
// 			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized a"})
// 		}

// 		// Verify the token
// 		fmt.Println(tokenString + os.Getenv("JWT_SECRET"))
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			// Verify the signing method
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, jwt.ErrSignatureInvalid
// 			}
// 			return []byte(os.Getenv("JWT_SECRET")), nil
// 		})

// 		if err != nil {
// 			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized b"})
// 		}

// 		// Token is valid
// 		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			// Continue to the next middleware or route handler
// 			return ctx.Next()
// 		}

// 		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized c"})
// 	}
// }
