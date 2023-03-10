package native

import (
	"credentials"

	"github.com/docker/docker-credential-helpers/client"
	"oras.land/oras-go/v2/registry/remote/auth"
)

const remoteCredentialsPrefix = "docker-credential-"

type Store struct {
	programFunc client.ProgramFunc
}

func GetStore(helperSuffix string) (credentials.Store, error) {
	return &Store{
		programFunc: client.NewShellProgramFunc(remoteCredentialsPrefix + helperSuffix),
	}, nil
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
