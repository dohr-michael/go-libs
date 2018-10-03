package crud

import (
	"github.com/dohr-michael/go-libs/storage"
	"net/http"
	"github.com/go-chi/chi"
	"strings"
	"context"
	"github.com/go-chi/render"
	"github.com/dohr-michael/go-libs/errors"
)

type Factory func() interface{}

type RouterBuilder struct {
	createFactory   Factory
	updateFactory   Factory
	readRepository  storage.ReadRepository
	writeRepository storage.WriteRepository
	fetchOne        http.HandlerFunc
	fetchMany       http.HandlerFunc
	create          http.HandlerFunc
	update          http.HandlerFunc
	delete          http.HandlerFunc
}

type Router func(prefix string, mux *chi.Mux) chi.Router

func NewRouterBuilder(read storage.ReadRepository, write storage.WriteRepository) *RouterBuilder {
	return &RouterBuilder{
		createFactory:   defaultOneFactory,
		updateFactory:   defaultOneFactory,
		fetchOne:        defaultFetchOne,
		fetchMany:       defaultFetchMany,
		create:          defaultCreate,
		update:          defaultUpdate,
		delete:          defaultDelete,
		readRepository:  read,
		writeRepository: write,
	}
}

func (h *RouterBuilder) WithCreateFactory(factory Factory) *RouterBuilder {
	h.createFactory = factory
	return h
}

func (h *RouterBuilder) WithUpdateFactory(factory Factory) *RouterBuilder {
	h.updateFactory = factory
	return h
}

func (h *RouterBuilder) WithFetchMany(fn http.HandlerFunc) *RouterBuilder {
	h.fetchMany = fn
	return h
}

func (h *RouterBuilder) WithFetchOne(fn http.HandlerFunc) *RouterBuilder {
	h.fetchOne = fn
	return h
}

func (h *RouterBuilder) WithCreate(fn http.HandlerFunc) *RouterBuilder {
	h.create = fn
	return h
}

func (h *RouterBuilder) WithUpdate(fn http.HandlerFunc) *RouterBuilder {
	h.update = fn
	return h
}

func (h *RouterBuilder) WithDelete(fn http.HandlerFunc) *RouterBuilder {
	h.delete = fn
	return h
}

func (h *RouterBuilder) Router(prefix string, mux *chi.Mux) chi.Router {
	return mux.Route("/"+strings.TrimPrefix(prefix, "/"), func(r chi.Router) {
		baseMiddleware := withBuilder(h)

		r.With(baseMiddleware, parseQuery).Get("/", h.fetchMany)
		r.With(baseMiddleware, readId).Get("/{id}", h.fetchOne)
		r.With(baseMiddleware, parseCreatePayload).Post("/", h.create)
		r.With(baseMiddleware, readId, parseUpdatePayload).Put("/{id}", h.update)
		r.With(baseMiddleware, readId).Delete("/{id}", h.delete)
	})
}

// Default implementations

func defaultOneFactory() interface{} {
	return &map[string]interface{}{}
}

func defaultFetchMany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := FindManyRequest(ctx)
	res, err := ReadRepository(ctx).FetchMany(query, ctx)
	if err != nil {
		render.Render(w, r, errors.InternalServerErrorRenderer(err))
		return
	}
	render.JSON(w, r, res)
}

func defaultFetchOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := Id(ctx)
	res, err := ReadRepository(ctx).FetchOne(id, ctx)
	if err != nil {
		render.Render(w, r, errors.ToRenderer(err))
		return
	}
	render.JSON(w, r, res)

}

func defaultCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload := CreatePayload(ctx)
	id, res, err := WriteRepository(ctx).Create(payload, ctx)
	if err != nil {
		render.Render(w, r, errors.ToRenderer(err))
		return
	}
	render.JSON(w, r.WithContext(context.WithValue(ctx, render.StatusCtxKey, http.StatusCreated)), &Created{Id: id, Item: res})
}

func defaultUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload := UpdatePayload(ctx)
	id := Id(ctx)
	res, err := WriteRepository(ctx).Update(id, payload, ctx)
	if err != nil {
		render.Render(w, r, errors.ToRenderer(err))
		return
	}
	render.JSON(w, r, res)
}

func defaultDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := Id(ctx)
	err := WriteRepository(ctx).Delete(id, ctx)
	if err != nil {
		render.Render(w, r, errors.ToRenderer(err))
		return
	}
	render.JSON(w, r, &Deleted{Id: id})
}
