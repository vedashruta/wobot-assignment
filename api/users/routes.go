package users

import (
	"server/env"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

func Route(router fiber.Router) {
	users := router.Group("/users")
	users.Post("/register", timeout.NewWithContext(register, env.Timeout))
	users.Post("/login", timeout.NewWithContext(login, env.Timeout))
}
