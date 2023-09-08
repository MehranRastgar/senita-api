package Controllers

import (
	Middleware "senita-api/middlewares"
	"senita-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserController is the controller for managing users.
type UserController struct {
	DB *gorm.DB
}

// NewUserController creates a new instance of UserController.
func NewUserController(database *gorm.DB) *UserController {
	return &UserController{DB: database}
}

// CreateUser creates a new user.
func (uc *UserController) CreateUser(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := uc.DB.Create(user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(user)
}

// GetUser retrieves a user by ID.
func (uc *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var user models.User

	if err := uc.DB.First(&user, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return ctx.JSON(user)
}

// UpdateUser updates an existing user by ID.
func (uc *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var user models.User

	if err := uc.DB.First(&user, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := uc.DB.Save(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(user)
}

// DeleteUser deletes a user by ID.
func (uc *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var user models.User

	if err := uc.DB.First(&user, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if err := uc.DB.Delete(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

// ListUsers retrieves a list of users.
func (uc *UserController) ListUsers(ctx *fiber.Ctx) error {
	var users []models.User
	var skip, limit int

	// Get the skip and limit values from the query string (if provided)
	if skipParam := ctx.Query("skip"); skipParam != "" {
		if parsedSkip, err := strconv.Atoi(skipParam); err == nil {
			skip = parsedSkip
		}
	}

	if limitParam := ctx.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil {
			limit = parsedLimit
		}
	}
	if skip == 0 {
		skip = 0
	}
	if limit == 0 {
		limit = 20
	}
	// Build the query with skip and limit
	query := uc.DB.Offset(skip).Limit(limit).Find(&users)

	if query.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": query.Error.Error()})
	}

	return ctx.JSON(users)
}

// RegisterRoutes registers the user routes.
func (uc *UserController) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	users := api.Group("/users")

	users.Post("/", uc.CreateUser)
	users.Get("/:id", uc.GetUser)
	users.Put("/:id", uc.UpdateUser)

	// These routes require authentication
	users.Use(Middleware.AuthMiddleware())

	users.Delete("/:id", uc.DeleteUser)
	users.Get("/", uc.ListUsers)
}
