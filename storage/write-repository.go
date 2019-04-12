package storage

import "context"

type WriteRepository interface {
	Create(entity interface{}, ctx context.Context) (string, Entity, error)
	Update(id string, toUpdate interface{}, ctx context.Context) (Entity, error)
	// Create or Update
	Save(id string, entity interface{}, ctx context.Context) (Entity, error)
	Delete(id string, ctx context.Context) error
}
