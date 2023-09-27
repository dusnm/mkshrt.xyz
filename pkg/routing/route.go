package routing

import (
	"github.com/gofiber/fiber/v2"
)

type (
	Route struct {
		Method   string
		Path     string
		Callback fiber.Handler
	}
)
