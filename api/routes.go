package api

import (
	"server/api/users"
	"server/middlewares/jwt"

	"github.com/gofiber/fiber/v2"
)

func Configure(app *fiber.App) {
	app.Use(jwt.Auth)
	// app.Group("/api", jwt.Auth)
	users.Route(app)
}
