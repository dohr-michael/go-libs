package storage

import "context"

type WriteRepository interface {
	Create(entity interface{}, ctx context.Context) (ID, interface{}, error)
	Update(id ID, toUpdate interface{}, ctx context.Context) (interface{}, error)
	// Create or Update
	Save(id ID, entity interface{}, ctx context.Context) (interface{}, error)
	Delete(id ID, ctx context.Context) error
}
