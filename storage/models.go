package storage

import "github.com/dohr-michael/go-libs/filters"

type Entity interface {
	Id() string
}
type Entities []Entity

type Paged struct {
	Items Entities       `json:"items"`
	Total int64          `json:"total"`
	Query *filters.Query `json:"query"`
}
