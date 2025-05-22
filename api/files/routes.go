package files

import (
	"server/env"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

func Route(router fiber.Router) {
	router.Post("/upload", timeout.NewWithContext(upload, env.Timeout))
	router.Get("/files", timeout.NewWithContext(files, env.Timeout))
	storage := router.Group("/storage")
	storage.Get("/remaining", timeout.NewWithContext(remaining, env.Timeout))
	storage.Get("/fetch/:id", timeout.NewWithContext(fetch, env.Timeout))
}
