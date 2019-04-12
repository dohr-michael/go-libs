package storage

import (
	"context"
	"github.com/dohr-michael/go-libs/filters"
)

type ReadRepository interface {
	FetchOne(id string, ctx context.Context) (Entity, error)
	FetchMany(query *filters.Query, ctx context.Context) (*Paged, error)
}
