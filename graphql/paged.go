package graphql

import (
	"errors"
	"fmt"
	"github.com/dohr-michael/go-libs/filters"
	"github.com/dohr-michael/go-libs/storage"
	"github.com/graphql-go/graphql"
)

func PagedQuery(prefix string, baseType *graphql.Object, repoContextKey string) *graphql.Field {
	return &graphql.Field{
		Type: PagedType(prefix, baseType),
		Args: graphql.FieldConfigArgument{
			"query": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"limit": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"offset": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			repo, ok := p.Context.Value(repoContextKey).(storage.ReadRepository)
			if !ok {
				return nil, errors.New(fmt.Sprintf("%s is not set", repoContextKey))
			}
			var conf []filters.QueryConfig
			if v, ok := p.Args["query"]; ok {
				conf = append(conf, filters.WithRsqlFilter(v.(string)))
			}
			if v, ok := p.Args["limit"]; ok {
				conf = append(conf, filters.WithLimit(int64(v.(int))))
			}
			if v, ok := p.Args["offset"]; ok {
				conf = append(conf, filters.WithOffset(int64(v.(int))))
			}

			query, err := filters.NewQuery(conf...)
			if err != nil {
				return nil, err
			}
			return repo.FetchMany(query, p.Context)
		},
	}
}

func PagedType(prefix string, baseType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        prefix + "Paged",
		Description: "",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(baseType))),
				Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
					if paged, ok := p.Source.(*storage.Paged); ok {
						return paged.Items, nil
					}
					return nil, nil
				},
			},
			"total": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
				Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
					if paged, ok := p.Source.(*storage.Paged); ok {
						return paged.Total, nil
					}
					return nil, nil
				},
			},
			"limit": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
				Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
					if paged, ok := p.Source.(*storage.Paged); ok && paged.Query != nil {
						return paged.Query.Pager.Limit, nil
					}
					return nil, nil
				},
			},
			"offset": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
				Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
					if paged, ok := p.Source.(*storage.Paged); ok && paged.Query != nil {
						return paged.Query.Pager.Offset, nil
					}
					return nil, nil
				},
			},
		},
	})
}
