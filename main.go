package main

import (
	"embed"
	"github.com/dusnm/mkshrt.xyz/pkg/container"
	"github.com/dusnm/mkshrt.xyz/pkg/core"
)

//go:embed views/*
var views embed.FS

//go:embed static/*
var staticAssets embed.FS

func main() {
	c := container.New(views, staticAssets)
	app := core.New(c)

	app.
		WireAppLevelMiddleware().
		WireRoutes().
		RegisterHooks().
		Listen()
}
