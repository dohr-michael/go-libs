package filters

import (
	"github.com/dohrm/go-rsql"
	"net/url"
	"strconv"
)

type Query struct {
	Filter Filter            `json:"filter"`
	Pager  Pager             `json:"pager"`
	Others map[string]string `json:"others"`
}

type QueryConfig func(*Query) (*Query, error)

func NewQuery(configs ...QueryConfig) (*Query, error) {
	res := &Query{
		Pager: Pager{
			Offset: 0,
			Limit:  100,
		},
	}
	var err error = nil
	for _, fn := range configs {
		res, err = fn(res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func WithRsqlFilter(rsqlFilter string) QueryConfig {
	return func(q *Query) (*Query, error) {
		filter, err := rsql.Parse(rsqlFilter)
		if err != nil {
			return q, err
		}
		q.Filter = filter
		return q, nil
	}
}

func WithFilter(filter Filter) QueryConfig {
	return func(q *Query) (*Query, error) {
		q.Filter = filter
		return q, nil
	}
}

func WithLimit(limit int64) QueryConfig {
	return func(q *Query) (*Query, error) {
		q.Pager.Limit = limit
		return q, nil
	}
}

func WithOffset(offset int64) QueryConfig {
	return func(q *Query) (*Query, error) {
		q.Pager.Offset = offset
		return q, nil
	}
}

func ParseHttpValues(base url.Values) (*Query, error) {
	var configs []QueryConfig
	for k, v := range base {
		switch {
		case k == "page[limit]":
			limit, err := strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				return nil, err
			}
			configs = append(configs, WithLimit(limit))
		case k == "page[offset]":
			offset, err := strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				return nil, err
			}
			configs = append(configs, WithOffset(offset))
		case k == "_q":
			configs = append(configs, WithRsqlFilter(v[0]))
		}

	}
	return NewQuery(configs...)
}
