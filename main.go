package main

import (
	"context"
	"embed"
	"github.com/dusnm/mkshrt.xyz/pkg/container"
	"github.com/dusnm/mkshrt.xyz/pkg/core"
	"os"
	"os/signal"
	"sync"
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

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		<-sig
		cancel()

		_ = app.Shutdown()
	}(&wg)

	app.
		WireAppLevelMiddleware().
		WireRoutes().
		RegisterHooks().
		Listen()

	wg.Wait()
}
