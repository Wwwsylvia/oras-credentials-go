package credentials

import "oras.land/oras-go/v2/registry/remote/auth"

// Store is the interface that any credentials store must implement.
type Store interface {
	// Store saves credentials into the store
	Store(serverAddress string, cred auth.Credential) error
	// Erase removes credentials from the store for the given server
	Erase(serverAddress string) error
	// Get retrieves credentials from the store for the given server
	Get(serverAddress string) (auth.Credential, error)
}
