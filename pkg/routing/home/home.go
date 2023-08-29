package home

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dusnm/mkshrt.xyz/pkg/config"
	"github.com/dusnm/mkshrt.xyz/pkg/models"
	"github.com/dusnm/mkshrt.xyz/pkg/repositories/mapping"
	"github.com/dusnm/mkshrt.xyz/pkg/routing"
	"github.com/dusnm/mkshrt.xyz/pkg/routing/home/data"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type (
	Handler struct {
		Config      *config.Config
		MappingRepo mapping.Interface
		Context     context.Context
	}
)

func (h Handler) Routes() []routing.Route {
	return []routing.Route{
		{
			Method:   http.MethodGet,
			Path:     "/",
			Callback: h.indexGet(),
		},
		{
			Method:   http.MethodPost,
			Path:     "/",
			Callback: h.indexPost(),
		},
		{
			Method:   http.MethodGet,
			Path:     "/:shortenKey<len(8)>",
			Callback: h.indexGetWithParam(),
		},
	}
}

func (h Handler) indexGet() routing.RouteCallback {
	return func(ctx *fiber.Ctx) error {
		return ctx.Render("views/home", fiber.Map{
			"Domain": ctx.Hostname(),
		})
	}
}

func (h Handler) indexPost() routing.RouteCallback {
	return func(ctx *fiber.Ctx) error {
		cntx, cancel := context.WithCancel(h.Context)
		defer cancel()

		domain := ctx.Hostname()
		d, err := data.New(ctx)
		if err != nil {
			return ctx.
				Status(http.StatusInternalServerError).
				Render("views/error", fiber.Map{})
		}

		if err = d.Validate(); err != nil {
			return ctx.
				Status(http.StatusUnprocessableEntity).
				Render("views/home", fiber.Map{
					"Domain": domain,
					"Url":    err.Error(),
				})
		}

		model, err := h.MappingRepo.Fetch(cntx, mapping.FieldUrl, d.Url)
		if err != nil {
			slog.ErrorCtx(cntx, err.Error())

			return ctx.
				Status(http.StatusInternalServerError).
				Render("views/error", fiber.Map{})
		}

		if model == (models.Mapping{}) {
			model, err = h.MappingRepo.Insert(cntx, d.Url)
			if err != nil {
				return ctx.
					Status(http.StatusInternalServerError).
					Render("views/error", fiber.Map{})
			}
		}

		return ctx.Render("views/home", fiber.Map{
			"Domain": ctx.Hostname(),
			"Url": fmt.Sprintf(
				"https://%s/%s",
				domain,
				model.ShortenKey,
			),
		})
	}
}

func (h Handler) indexGetWithParam() routing.RouteCallback {
	return func(ctx *fiber.Ctx) error {
		cntx, cancel := context.WithCancel(h.Context)
		defer cancel()

		shortenKey := ctx.Params("shortenKey", "")
		model, err := h.MappingRepo.Fetch(cntx, mapping.FieldShortenKey, shortenKey)
		if err != nil {
			slog.ErrorCtx(cntx, err.Error())

			return ctx.
				Status(http.StatusInternalServerError).
				Render("views/error", fiber.Map{})
		}

		if model == (models.Mapping{}) {
			return ctx.
				Status(http.StatusNotFound).
				Render("views/error-404", fiber.Map{})
		}

		return ctx.Redirect(model.Url, http.StatusFound)
	}
}
