package storage

import (
	"context"
	"github.com/dohr-michael/go-libs/filters"
)

type ID string

type Paged struct {
	Items []interface{}  `json:"items"`
	Total int64          `json:"total"`
	Query *filters.Query `json:"query"`
}

type ReadRepository interface {
	FetchOne(id ID, ctx context.Context) (interface{}, error)
	FetchMany(query *filters.Query, ctx context.Context) (*Paged, error)
}
