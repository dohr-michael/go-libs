package filters

import (
	"net/url"
	"strconv"
	"github.com/dohrm/go-rsql"
)

type Query struct {
	Filter Filter            `json:"filter"`
	Pager  Pager             `json:"pager"`
	Others map[string]string `json:"others"`
}

func ParseHttpValues(base url.Values) (*Query, error) {
	res := &Query{
		Pager: Pager{
			Offset: 0,
			Limit:  100,
		},
	}
	for k, v := range base {
		switch {
		case k == "page[limit]":
			limit, err := strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				return nil, err
			}
			res.Pager.Limit = limit
		case k == "page[offset]":
			offset, err := strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				return nil, err
			}
			res.Pager.Offset = offset
		case k == "_q":
			filter, err := rsql.Parse(v[0])
			if err != nil {
				return nil, err
			}
			res.Filter = filter
		}

	}
	return res, nil
}
