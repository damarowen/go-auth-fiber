package main

import (
	"auth-go-fiber/database"
	"auth-go-fiber/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/asaskevich/govalidator"
)



func main() {
	govalidator.SetFieldsRequiredByDefault(true)

	database.Connect()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	routes.Setup(app)

	app.Listen(":3000")
}
