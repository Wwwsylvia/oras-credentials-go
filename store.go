package credentials

import (
	"context"

	"oras.land/oras-go/v2/registry/remote/auth"
)

// Store is the interface that any credentials store must implement.
type Store interface {
	// Store saves credentials into the store
	Store(ctx context.Context, serverAddress string, cred auth.Credential) error
	// Erase removes credentials from the store for the given server
	Erase(ctx context.Context, serverAddress string) error
	// Get retrieves credentials from the store for the given server
	Get(ctx context.Context, serverAddress string) (auth.Credential, error)
}
