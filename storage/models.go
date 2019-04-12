package storage

import "github.com/dohr-michael/go-libs/filters"

type Entity interface{}
type Entities interface{}

type Paged struct {
	Items Entities       `json:"items"`
	Total int64          `json:"total"`
	Query *filters.Query `json:"query"`
}
