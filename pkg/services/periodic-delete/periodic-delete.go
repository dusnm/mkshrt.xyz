package periodic_delete

import (
	"context"
	"fmt"
	"github.com/dusnm/mkshrt.xyz/pkg/repositories/mapping"
	"os"
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
	go func(ctx context.Context, ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				if err := s.mappingRepo.DeleteOldEntries(ctx); err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx, ticker)
}
