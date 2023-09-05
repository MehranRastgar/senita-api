package Controllers

import (
	"senita-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ArticleController is the controller for managing articles.
type ArticleController struct {
	DB *gorm.DB
}

// NewArticleController creates a new instance of ArticleController.
func NewArticleController(database *gorm.DB) *ArticleController {
	return &ArticleController{DB: database}
}

// CreateArticle creates a new article.
func (ac *ArticleController) CreateArticle(ctx *fiber.Ctx) error {
	article := new(models.Article)
	if err := ctx.BodyParser(article); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := ac.DB.Create(article).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(article)
}

// GetArticle retrieves an article by ID.
func (ac *ArticleController) GetArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var article models.Article

	if err := ac.DB.First(&article, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Article not found"})
	}

	return ctx.JSON(article)
}

// UpdateArticle updates an existing article by ID.
func (ac *ArticleController) UpdateArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var article models.Article

	if err := ac.DB.First(&article, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Article not found"})
	}

	if err := ctx.BodyParser(&article); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := ac.DB.Save(&article).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(article)
}

// DeleteArticle deletes an article by ID.
func (ac *ArticleController) DeleteArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var article models.Article

	if err := ac.DB.First(&article, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Article not found"})
	}

	if err := ac.DB.Delete(&article).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

// ListArticles retrieves a list of articles.
// func (ac *ArticleController) ListArticles(ctx *fiber.Ctx) error {
// 	var articles []models.Article

// 	if err := ac.DB.Find(&articles).Error; err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}

//		return ctx.JSON(articles)
//	}
func (ac *ArticleController) ListArticles(ctx *fiber.Ctx) error {
	var articles []models.Article
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
	query := ac.DB.Offset(skip).Limit(limit).Find(&articles)

	if query.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": query.Error.Error()})
	}

	return ctx.JSON(articles)
}

// RegisterRoutes registers the article routes.
func (ac *ArticleController) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	articles := api.Group("/articles")

	articles.Post("/", ac.CreateArticle)
	articles.Get("/:id", ac.GetArticle)
	articles.Put("/:id", ac.UpdateArticle)
	articles.Delete("/:id", ac.DeleteArticle)
	articles.Get("/", ac.ListArticles)

}
