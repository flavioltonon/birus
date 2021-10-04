package repository

import (
	"context"
)

// Repository defines a repository
type Repository interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
}
