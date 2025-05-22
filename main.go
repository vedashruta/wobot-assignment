package main

import (
	"log"
	"server/api"
	"server/env"

	"github.com/gofiber/fiber/v2"
)

func main() {
	err := env.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	api.Configure(app)
	err = app.Listen(env.Port)
	if err != nil {
		log.Fatal(err)
	}
}
