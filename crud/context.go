package crud

import (
	"github.com/dohr-michael/go-libs/filters"
	"github.com/dohr-michael/go-libs/errors"
	"github.com/dohr-michael/go-libs/storage"
	"github.com/go-chi/render"
	"github.com/go-chi/chi"
	"context"
	"encoding/json"
	"net/http"
)

// Props Contexts
const CreateFactoryCtx = "crud.CreateFactoryCtx"
const UpdateFactoryCtx = "crud.UpdateFactoryCtx"
const ReadRepositoryCtx = "crud.ReadRepositoryCtx"
const WriteRepositoryCtx = "crud.WriteRepositoryCtx"

// Parse Ctx.
const FindManyRequestCtx = "crud.FindManyRequestCtx"
const CreatePayloadCtx = "crud.CreatePayloadCtx"
const UpdatePayloadCtx = "crud.UpdateRequestCtx"
const IdCtx = "crud.IdCtx"

func CreateFactory(ctx context.Context) Factory {
	return ctx.Value(CreateFactoryCtx).(Factory)
}

func UpdateFactory(ctx context.Context) Factory {
	return ctx.Value(UpdateFactoryCtx).(Factory)
}

func ReadRepository(ctx context.Context) storage.ReadRepository {
	return ctx.Value(ReadRepositoryCtx).(storage.ReadRepository)
}

func WriteRepository(ctx context.Context) storage.WriteRepository {
	return ctx.Value(WriteRepositoryCtx).(storage.WriteRepository)
}

func FindManyRequest(ctx context.Context) *filters.Query {
	return ctx.Value(FindManyRequestCtx).(*filters.Query)
}

func CreatePayload(ctx context.Context) interface{} {
	return ctx.Value(CreatePayloadCtx)
}

func UpdatePayload(ctx context.Context) interface{} {
	return ctx.Value(UpdatePayloadCtx)
}

func Id(ctx context.Context) string {
	return ctx.Value(IdCtx).(string)
}

func withBuilder(builder *RouterBuilder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), CreateFactoryCtx, builder.createFactory)
			ctx = context.WithValue(ctx, UpdateFactoryCtx, builder.updateFactory)
			ctx = context.WithValue(ctx, ReadRepositoryCtx, builder.readRepository)
			ctx = context.WithValue(ctx, WriteRepositoryCtx, builder.writeRepository)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query, err := filters.ParseHttpValues(r.URL.Query())
		if err != nil {
			render.Render(w, r, errors.InvalidRequestRenderer(err))
			return
		}
		ctx := context.WithValue(r.Context(), FindManyRequestCtx, query)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func readId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := string(chi.URLParam(r, "id"))
		if id == "" {
			render.Render(w, r, errors.NotFoundRenderer)
			return
		}
		ctx := context.WithValue(r.Context(), IdCtx, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseCreatePayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := CreateFactory(ctx)()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(req)
		if err != nil {
			render.Render(w, r, errors.InvalidRequestRenderer(err))
			return
		}
		// TODO validation.

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CreatePayloadCtx, req)))
	})
}

func parseUpdatePayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := UpdateFactory(ctx)()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(req)
		if err != nil {
			render.Render(w, r, errors.InvalidRequestRenderer(err))
			return
		}
		// TODO validation.

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UpdatePayloadCtx, req)))
	})
}
