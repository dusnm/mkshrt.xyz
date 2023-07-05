package routing

import (
	"github.com/gofiber/fiber/v2"
)

type (
	RouteCallback func(ctx *fiber.Ctx) error

	Route struct {
		Method   string
		Path     string
		Callback RouteCallback
	}
)
