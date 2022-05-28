package main

import (
	"go-refresh/internals"
	"go-refresh/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
)

func main() {
	app := fiber.New()
	err := internals.AppInternals.Connect()
	if err != nil {
		panic(err)
	}

	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(logger.New())

	routers.Router.SetupUserRouter(app)
	app.Listen(":3000")
}
