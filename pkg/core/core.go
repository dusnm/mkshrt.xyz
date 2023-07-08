package core

import (
	"context"
	"fmt"
	"github.com/dusnm/mkshrt.xyz/pkg/config"
	"github.com/dusnm/mkshrt.xyz/pkg/container"
	"github.com/dusnm/mkshrt.xyz/pkg/routing"
	"github.com/dusnm/mkshrt.xyz/pkg/routing/home"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	"os"
)

type (
	Core struct {
		Application *fiber.App
		Container   *container.Container
		context     context.Context
	}
)

func New(c *container.Container, ctx context.Context) *Core {
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
		context:   ctx,
	}
}

func (c *Core) WireLogging() *Core {
	var level slog.Level
	var handler slog.Handler

	cfg := c.Container.GetConfig()

	err := level.UnmarshalText([]byte(cfg.Logging.Level))
	if err != nil {
		log.Fatal(err)
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	f, err := os.OpenFile(cfg.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	switch cfg.Logging.Format {
	case config.LoggingFormatText:
		handler = slog.NewTextHandler(f, opts)
	case config.LoggingFormatJson:
		handler = slog.NewJSONHandler(f, opts)
	default:
		log.Fatal("invalid logging format")
	}

	slog.SetDefault(slog.New(handler))

	return c
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
		Context:     c.context,
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
	c.Application.Hooks().OnListen(func() error {
		c.Container.GetPeriodicDeleteService().Work(c.context)

		return nil
	})

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

func (c *Core) Shutdown() error {
	return c.Application.Shutdown()
}
