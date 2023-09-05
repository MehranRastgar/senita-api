package Controllers

import (
	"senita-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CategoryController is the controller for managing categories.
type CategoryController struct {
    DB *gorm.DB
}

// NewCategoryController creates a new instance of CategoryController.
func NewCategoryController(database *gorm.DB) *CategoryController {
    return &CategoryController{DB: database}
}

// CreateCategory creates a new category.
func (cc *CategoryController) CreateCategory(ctx *fiber.Ctx) error {
    category := new(models.Category)
    if err := ctx.BodyParser(category); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    if err := cc.DB.Create(category).Error; err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return ctx.JSON(category)
}

// GetCategory retrieves a category by ID.
func (cc *CategoryController) GetCategory(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    var category models.Category

    if err := cc.DB.First(&category, id).Error; err != nil {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Category not found"})
    }

    return ctx.JSON(category)
}

// UpdateCategory updates an existing category by ID.
func (cc *CategoryController) UpdateCategory(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    var category models.Category

    if err := cc.DB.First(&category, id).Error; err != nil {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Category not found"})
    }

    if err := ctx.BodyParser(&category); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    if err := cc.DB.Save(&category).Error; err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return ctx.JSON(category)
}

// DeleteCategory deletes a category by ID.
func (cc *CategoryController) DeleteCategory(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    var category models.Category

    if err := cc.DB.First(&category, id).Error; err != nil {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Category not found"})
    }

    if err := cc.DB.Delete(&category).Error; err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return ctx.SendStatus(fiber.StatusNoContent)
}

// ListCategories retrieves a list of categories.
func (cc *CategoryController) ListCategories(ctx *fiber.Ctx) error {
    var categories []models.Category
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

    // Build the query with skip and limit
    query := cc.DB.Offset(skip).Limit(limit).Find(&categories)

    if query.Error != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": query.Error.Error()})
    }

    return ctx.JSON(categories)
}

// RegisterRoutes registers the category routes.
func (cc *CategoryController) RegisterRoutes(app *fiber.App) {
    api := app.Group("/api")
    categories := api.Group("/categories")

    categories.Post("/", cc.CreateCategory)
    categories.Get("/:id", cc.GetCategory)
    categories.Put("/:id", cc.UpdateCategory)
    categories.Delete("/:id", cc.DeleteCategory)
    categories.Get("/", cc.ListCategories)
}
