package models

import "time"

type (
	Mapping struct {
		Url        string
		ShortenKey string
		CreatedAt  time.Time
	}
)
