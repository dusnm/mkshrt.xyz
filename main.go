package main

import (
	"context"
	"embed"
	"github.com/dusnm/mkshrt.xyz/pkg/container"
	"github.com/dusnm/mkshrt.xyz/pkg/core"
	"os"
	"os/signal"
	"syscall"
)

//go:embed views/*
var views embed.FS

//go:embed static/*
var staticAssets embed.FS

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	c := container.New(views, staticAssets)
	app := core.New(c, ctx)

	go func(cancel context.CancelFunc, sig chan os.Signal) {
		<-sig

		cancel()
	}(cancel, sig)

	app.
		WireAppLevelMiddleware().
		WireRoutes().
		RegisterHooks().
		Listen()
}
