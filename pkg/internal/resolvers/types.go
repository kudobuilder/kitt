package resolvers

import (
	"context"

	"github.com/spf13/afero"
)

// Resolver retrieves an operator package into a file system.
type Resolver interface {
	Resolve(context.Context) (fs afero.Fs, remover Remover, err error)
}

// Remover cleans up temporary objects created by the resolver.
// To be called once the file system provided by the resolver is no longer used.
type Remover func() error
