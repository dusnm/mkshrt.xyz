package periodic_delete

import (
	"context"
	"github.com/dusnm/mkshrt.xyz/pkg/repositories/mapping"
	"golang.org/x/exp/slog"
	"time"
)

type (
	Interface interface {
		Work(ctx context.Context)
	}

	Service struct {
		mappingRepo mapping.Interface
	}
)

func New(mappingRepo mapping.Interface) Service {
	return Service{
		mappingRepo: mappingRepo,
	}
}

func (s Service) Work(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)

	go func(
		ctx context.Context,
		ticker *time.Ticker,
		repo mapping.Interface,
	) {
		for {
			select {
			case <-ticker.C:
				if err := repo.DeleteOldEntries(ctx); err != nil {
					slog.ErrorCtx(ctx, err.Error())
				} else {
					slog.InfoCtx(ctx, "deleted old entries")
				}
			case <-ctx.Done():
				ticker.Stop()

				return
			}
		}
	}(ctx, ticker, s.mappingRepo)
}
