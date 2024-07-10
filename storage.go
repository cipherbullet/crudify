package main

import (
	"context"
)

type StorageInterface interface {
	Connect(connStr string) error
	Disconnect() error

	All(ctx context.Context, dest interface{}) (error)
	Find(ctx context.Context, dest interface{}, id int) (error)
	Create(ctx context.Context, dest interface{}) (error)
	Delete(ctx context.Context, dest interface{}) (error)
}