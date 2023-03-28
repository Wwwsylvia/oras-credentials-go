package credentials

import (
	"context"

	"oras.land/oras-go/v2/registry/remote/auth"
)

// Store is the interface that any credentials store must implement.
type Store interface {
	// Put saves credentials into the store
	Put(ctx context.Context, serverAddress string, cred auth.Credential) error
	// Delete removes credentials from the store for the given server
	Delete(ctx context.Context, serverAddress string) error
	// Get retrieves credentials from the store for the given server
	Get(ctx context.Context, serverAddress string) (auth.Credential, error)
}
