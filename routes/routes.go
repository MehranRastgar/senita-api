package Routes

import (
	"os"
	Controllers "senita-api/controllers"
	"senita-api/db"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// app.Post("/users/:userId/login", Controllers.Login)
	// app.Get("users/:userId/logout", Controllers.Logout)
	// app.Post("/users/userId/passcode", Controllers.Passcode)

	//user routes
	// app.Post("/users", Controllers.CreateUser)
	// app.Post("/users", Controllers.UpdateUser)
	// app.Get("/users", Controllers.UserList)
	// app.Get("/users/:userId", Controllers.UserDetails)
	// app.Put("/users/:userId", Controllers.EditUser)
	// app.Delete("/users/:userId", Controllers.DeleteUser)
	db := db.DB // Initialize your GORM database connection

	articleController := Controllers.NewArticleController(db)
	categoryController := Controllers.NewCategoryController(db)
	userController := Controllers.NewUserController(db)

	authController := Controllers.NewAuthController(os.Getenv("JWT_SECRET"))

	// Register authentication routes

	articleController.RegisterRoutes(app)
	categoryController.RegisterRoutes(app)
	userController.RegisterRoutes(app)
	authController.RegisterRoutes(app)
}
