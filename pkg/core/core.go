package core

import (
	"fmt"
	"github.com/dusnm/mkshrt.xyz/pkg/container"
	"github.com/dusnm/mkshrt.xyz/pkg/routing"
	"github.com/dusnm/mkshrt.xyz/pkg/routing/home"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"log"
	"net/http"
)

type (
	Core struct {
		Application *fiber.App
		Container   *container.Container
	}
)

func New(c *container.Container) *Core {
	cfg := c.GetConfig()

	return &Core{
		Application: fiber.New(fiber.Config{
			Views:                   html.NewFileSystem(http.FS(c.Views), ".html"),
			ViewsLayout:             "views/layouts/main",
			EnableTrustedProxyCheck: cfg.Application.IsBehindProxy,
			TrustedProxies:          cfg.Application.TrustedProxies,
			PassLocalsToViews:       true,
		}),
		Container: c,
	}
}

func (c *Core) WireAppLevelMiddleware() *Core {
	c.Application.Use("/static", filesystem.New(
		filesystem.Config{
			Root:       http.FS(c.Container.Assets),
			PathPrefix: "static",
			Browse:     false,
		},
	))

	return c
}

func (c *Core) WireRoutes() *Core {
	homeHandler := home.Handler{
		Config:      c.Container.GetConfig(),
		MappingRepo: c.Container.GetMappingRepository(),
	}

	routes := make([]routing.Route, 0, 10)
	routes = append(
		routes,
		homeHandler.Routes()...,
	)

	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			c.Application.Get(route.Path, route.Callback)
		case http.MethodPost:
			c.Application.Post(route.Path, route.Callback)
		case http.MethodPut:
			c.Application.Put(route.Path, route.Callback)
		case http.MethodDelete:
			c.Application.Delete(route.Path, route.Callback)
		}
	}

	return c
}

func (c *Core) RegisterHooks() *Core {
	c.Application.Hooks().OnShutdown(func() error {
		if err := c.Container.Close(); err != nil {
			return err
		}

		return nil
	})

	return c
}

func (c *Core) Listen() {
	cfg := c.Container.GetConfig()

	host := cfg.Application.Host
	port := cfg.Application.Port
	socket := fmt.Sprintf("%s:%d", host, port)

	if err := c.Application.Listen(socket); err != nil {
		log.Fatal(err)
	}
}
