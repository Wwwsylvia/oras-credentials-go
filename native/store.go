package native

import (
	"github.com/docker/docker-credential-helpers/client"
	"oras.land/oras-go/v2/registry/remote/auth"
)

type Store struct {
	programFunc client.ProgramFunc
}

func GetStore() *Store {
	return &Store{}
}

// Store saves credentials into the store
func (s *Store) Store(serverAddress string, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}

// Erase removes credentials from the store for the given server
func (s *Store) Erase(serverAddress string) error {
	panic("not implemented") // TODO: Implement
}

// Get retrieves credentials from the store for the given server
func (s *Store) Get(serverAddress string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}
