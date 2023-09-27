package home

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dusnm/mkshrt.xyz/pkg/utils/qr"

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

func (h Handler) indexGet() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Render("views/home", fiber.Map{
			"Domain": ctx.Hostname(),
		})
	}
}

func (h Handler) indexPost() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.SetUserContext(h.Context)

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

		model, err := h.MappingRepo.Fetch(ctx.UserContext(), mapping.FieldUrl, d.Url)
		if err != nil {
			slog.ErrorContext(ctx.UserContext(), err.Error())

			return ctx.
				Status(http.StatusInternalServerError).
				Render("views/error", fiber.Map{})
		}

		if model == (models.Mapping{}) {
			model, err = h.MappingRepo.Insert(ctx.UserContext(), d.Url)
			if err != nil {
				return ctx.
					Status(http.StatusInternalServerError).
					Render("views/error", fiber.Map{})
			}
		}

		url := fmt.Sprintf(
			"https://%s/%s",
			domain,
			model.ShortenKey,
		)

		qrB64, err := qr.GenerateWithBase64Encoding(url)
		if err != nil {
			slog.ErrorContext(ctx.UserContext(), err.Error())

			return ctx.
				Status(http.StatusInternalServerError).
				Render("views/error", fiber.Map{})
		}

		return ctx.Render("views/home", fiber.Map{
			"Domain": ctx.Hostname(),
			"Url":    url,
			"QRCode": qrB64,
		})
	}
}

func (h Handler) indexGetWithParam() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.SetUserContext(h.Context)

		shortenKey := ctx.Params("shortenKey", "")
		model, err := h.MappingRepo.Fetch(ctx.UserContext(), mapping.FieldShortenKey, shortenKey)
		if err != nil {
			slog.ErrorContext(ctx.UserContext(), err.Error())

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
