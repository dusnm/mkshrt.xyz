package mapping

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dusnm/mkshrt.xyz/pkg/models"
	"github.com/dusnm/mkshrt.xyz/pkg/utils/random"
	"time"
)

const (
	TableName       = "mappings"
	FieldUrl        = "url"
	FieldShortenKey = "shorten_key"
	FieldCreatedAt  = "created_at"
)

var _ Interface = &Repository{}

type (
	Interface interface {
		Fetch(ctx context.Context, searchKey string, searchValue string) (models.Mapping, error)
		Insert(ctx context.Context, url string) (models.Mapping, error)
	}

	Repository struct {
		db *sql.DB
	}
)

func New(db *sql.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Fetch(
	ctx context.Context,
	searchKey string,
	searchValue string,
) (models.Mapping, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	mapping := models.Mapping{}

	query := fmt.Sprintf(
		"SELECT %s, %s, %s FROM %s WHERE %s = ? LIMIT 1",
		FieldUrl,
		FieldShortenKey,
		FieldCreatedAt,
		TableName,
		searchKey,
	)

	stmt, err := r.db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return mapping, err
	}

	err = stmt.
		QueryRowContext(ctx, searchValue).
		Scan(
			&mapping.Url,
			&mapping.ShortenKey,
			&mapping.CreatedAt,
		)

	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return mapping, nil
		}

		return mapping, err
	}

	return mapping, nil
}

func (r Repository) Insert(ctx context.Context, url string) (models.Mapping, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	shortenKey, err := random.UniqueString()
	if err != nil {
		return models.Mapping{}, err
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s, %s) VALUES(?, ?)",
		TableName,
		FieldUrl,
		FieldShortenKey,
	)

	stmt, err := r.db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return models.Mapping{}, err
	}

	_, err = stmt.ExecContext(ctx, url, shortenKey)
	if err != nil {
		return models.Mapping{}, err
	}

	return models.Mapping{
		Url:        url,
		ShortenKey: shortenKey,
		CreatedAt:  time.Now(),
	}, nil
}
