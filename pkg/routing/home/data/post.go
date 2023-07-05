package data

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/url"
)

type (
	Data struct {
		Url string `form:"url"`
	}
)

func New(ctx *fiber.Ctx) (*Data, error) {
	d := new(Data)
	if err := ctx.BodyParser(d); err != nil {
		return &Data{}, err
	}

	return d, nil
}

func (d *Data) Validate() error {
	errorMessage := "invalid url"
	if d.Url == "" {
		return errors.New(errorMessage)
	}

	uri, err := url.Parse(d.Url)
	if err != nil {
		return errors.New(errorMessage)
	}

	if uri.Scheme == "" || uri.Host == "" {
		return errors.New(errorMessage)
	}

	if len(d.Url) > 1000 {
		return errors.New("only urls of less than a 1000 characters can be shortened")
	}

	return nil
}
