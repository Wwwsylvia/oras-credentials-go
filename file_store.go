package credentials

import (
	"context"

	"oras.land/oras-go/v2/registry/remote/auth"
)

// FileStore implements a credentials store using
// the docker configuration file to keep the credentials in plain text.
type FileStore struct {
	configPath  string
	DisableSave bool
}

// NewFileStore creates a new file credentials store.
func NewFileStore(configPath string) Store {
	return &FileStore{
		configPath: configPath,
	}
}

// Store saves credentials into the store
func (fs *FileStore) Store(ctx context.Context, serverAddress string, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}

// Erase removes credentials from the store for the given server
func (fs *FileStore) Erase(ctx context.Context, serverAddress string) error {
	panic("not implemented") // TODO: Implement
}

// Get retrieves credentials from the store for the given server
func (fs *FileStore) Get(ctx context.Context, serverAddress string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}