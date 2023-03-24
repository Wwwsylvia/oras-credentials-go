package credentials

import (
	"context"

	"github.com/docker/docker-credential-helpers/client"
	"oras.land/oras-go/v2/registry/remote/auth"
)

const remoteCredentialsPrefix = "docker-credential-"

// NativeStore implements a credentials store
// using native keychain to keep credentials secure.
type NativeStore struct {
	programFunc client.ProgramFunc
}

// NewNativeStore creates a new native store that uses a remote helper program to
// manage credentials.
func NewNativeStore(helperSuffix string) Store {
	return &NativeStore{
		programFunc: client.NewShellProgramFunc(remoteCredentialsPrefix + helperSuffix),
	}
}

// Store saves credentials into the store
func (ns *NativeStore) Store(ctx context.Context, serverAddress string, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}

// Erase removes credentials from the store for the given server
func (ns *NativeStore) Erase(ctx context.Context, serverAddress string) error {
	panic("not implemented") // TODO: Implement
}

// Get retrieves credentials from the store for the given server
func (ns *NativeStore) Get(ctx context.Context, serverAddress string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}
