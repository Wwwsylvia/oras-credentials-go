package credentials

import (
	"oras.land/oras-go/v2/registry/remote/auth"
)

// fileStore implements a credentials store using
// the docker configuration file to keep the credentials in plain text.
type fileStore struct {
	configPath string
}

func NewFileStore(configPath string) (Store, error) {
	return &fileStore{}, nil
}

// Store saves credentials into the store
func (fs *fileStore) Store(serverAddress string, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}

// Erase removes credentials from the store for the given server
func (fs *fileStore) Erase(serverAddress string) error {
	panic("not implemented") // TODO: Implement
}

// Get retrieves credentials from the store for the given server
func (fs *fileStore) Get(serverAddress string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}
