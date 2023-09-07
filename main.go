package main

import (
	"fmt"
	"log"
	"os"
	Controllers "senita-api/controllers"
	"senita-api/db"

	Routes "senita-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	if os.Getenv("ENV") == "PRODUCTION" {
		fmt.Println("is PRODUCTION env")

	}
	if os.Getenv("ENV") == "LOCAL" {
		fmt.Println("is local env")
	}

	//Start PostgreSQL database
	//Example: db.GetDB() - More info in the models folder
	db.Init()

	//Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
	//Example: db.GetRedis().Set(KEY, VALUE, at.Sub(now)).Err()
	db.InitRedis(1)

	fmt.Println("creating structure.")

	fmt.Println("app running")

	app := fiber.New()
	app.Use(cors.New())
	app.Get("/testApi", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Go fiber first api created",
		})
	})

	app.Get("/autoMigrateWithToken", func(ctx *fiber.Ctx) error {

		token := ctx.Query("token")
		if token == os.Getenv("MIGRATION_TOKEN") {
			db.AutoMigrate()
			return ctx.Status(200).JSON(fiber.Map{
				"success": true,
				"message": "Go fiber first api created",
			})
		}
		return ctx.Status(200).JSON(fiber.Map{
			"success": false,
			"message": "Token not true",
		})
	})
	// Create a new app instance
	db := db.DB // Initialize your GORM database connection

	articleController := Controllers.NewArticleController(db)
	categoryController := Controllers.NewCategoryController(db)
	userController := Controllers.NewUserController(db)

	Routes.Setup(app)
	articleController.RegisterRoutes(app)
	categoryController.RegisterRoutes(app)
	userController.RegisterRoutes(app)

	app.Listen(":" + os.Getenv("PORT"))

}
