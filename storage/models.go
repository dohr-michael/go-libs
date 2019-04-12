package storage

import "github.com/dohr-michael/go-libs/filters"

type Entity interface{}

type Paged struct {
	Items []Entity       `json:"items"`
	Total int64          `json:"total"`
	Query *filters.Query `json:"query"`
}
