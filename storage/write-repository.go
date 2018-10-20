package storage

import "context"

type WriteRepository interface {
	Create(entity interface{}, ctx context.Context) (string, interface{}, error)
	Update(id string, toUpdate interface{}, ctx context.Context) (interface{}, error)
	// Create or Update
	Save(id string, entity interface{}, ctx context.Context) (interface{}, error)
	Delete(id string, ctx context.Context) error
}
