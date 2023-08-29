package mapping

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dusnm/mkshrt.xyz/pkg/models"
	"github.com/dusnm/mkshrt.xyz/pkg/utils/random"
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
		DeleteOldEntries(ctx context.Context) error
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
	if err != nil {
		return mapping, err
	}

	defer stmt.Close()

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

	// 6 bytes gives us 2^48 possible combinations
	// The approximate solution to the birthday problem
	// gives us a collision probability of ~ 2.13×10^-14
	// which is extremely unlikely, but the resulting string
	// contains only 8 characters after base64 encoding, which
	// I feel is an acceptable trade-off
	shortenKey, err := random.String(6)
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
	if err != nil {
		return models.Mapping{}, err
	}

	defer stmt.Close()

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

func (r Repository) DeleteOldEntries(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE DATEDIFF(CURDATE(), DATE(created_at)) > 30",
		TableName,
	)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
