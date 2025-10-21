package shortener

import (
	"context"
)

// Store defines persistence behaviors for Link records.
type Store interface {
	Create(ctx context.Context, l Link) error
	Get(ctx context.Context, code string) (Link, bool, error)
	Delete(ctx context.Context, code string) error
	IncrementHit(ctx context.Context, code string) (Link, error)
	List(ctx context.Context) ([]Link, error)
}
