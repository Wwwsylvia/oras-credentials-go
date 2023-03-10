package file

import (
	"credentials"

	"oras.land/oras-go/v2/registry/remote/auth"
)

type Store struct {
	config credentials.AuthConfig
}

func GetStore(configPath string) (credentials.Store, error) {
	return &Store{}, nil
}

// Store saves credentials into the store
func (cs *Store) Store(serverAddress string, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}

// Erase removes credentials from the store for the given server
func (cs *Store) Erase(serverAddress string) error {
	panic("not implemented") // TODO: Implement
}

// Get retrieves credentials from the store for the given server
func (cs *Store) Get(serverAddress string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}
