package graphql

import (
	"errors"
	"fmt"
	lerrors "github.com/dohr-michael/go-libs/errors"
	"github.com/dohr-michael/go-libs/storage"
	"github.com/graphql-go/graphql"
)

func ById(fieldId string, baseType *graphql.Object, repoContextKey string) *graphql.Field {
	return &graphql.Field{
		Type: baseType,
		Args: graphql.FieldConfigArgument{
			fieldId: &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			repo, ok := p.Context.Value(repoContextKey).(storage.ReadRepository)
			if !ok {
				return nil, errors.New(fmt.Sprintf("%s is not set", repoContextKey))
			}
			res, err := repo.FetchOne(p.Args[fieldId].(string), p.Context)
			if err == lerrors.NotFoundError {
				return nil, nil
			} else if err != nil {
				return nil, err
			}
			return res, nil
		},
	}
}
