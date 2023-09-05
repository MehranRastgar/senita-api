package Controllers

import (
	"log"
	"senita-api/db"
	models "senita-api/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("registeration error in post request %v", err)
	}

	if data["user_name"] == "" || data["password"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "user user_name and password is required",
			"error":   map[string]interface{}{},
		})
	}
	//passCode := strconv.Itoa(rand.Intn(1000000))
	//fmt.Println("passCode:::", passCode

	user := models.User{
		UserName:  data["user_name"],
		Password:  data["password"],
		Email:     data["email"],
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	db.DB.Create(&user)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})

}

func EditUser(c *fiber.Ctx) error {
	return nil

}

func UpdateUser(c *fiber.Ctx) error {
	return nil

}

// Cashiers struct with two values
type Users struct {
	Id       uint   `json:"cashierId"`
	UserName string `json:"user_name"`
}

func UserList(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var user []Users
	db.DB.Select("*").Limit(limit).Offset(skip).Find(&user).Count(&count)
	metaMap := map[string]interface{}{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}
	usersData := map[string]interface{}{
		"users": user,
		"meta":  metaMap,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    usersData,
	})

}
func UserDetails(c *fiber.Ctx) error {
	return nil

}
func DeleteUser(c *fiber.Ctx) error {
	return nil

}
