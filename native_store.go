package credentials

import (
	"github.com/docker/docker-credential-helpers/client"
	"oras.land/oras-go/v2/registry/remote/auth"
)

const remoteCredentialsPrefix = "docker-credential-"

// nativeStore implements a credentials store
// using native keychain to keep credentials secure.
type nativeStore struct {
	programFunc client.ProgramFunc
	fs          *fileStore
}

type NativeStoreOptions struct {
	DisablePlainTextSave bool
}

// NewNativeStore creates a new native store that
// uses a remote helper program to manage credentials.
func NewNativeStore(configPath, helperSuffix string, opts NativeStoreOptions) Store {
	return &nativeStore{
		fs:          &fileStore{configPath: configPath},
		programFunc: client.NewShellProgramFunc(remoteCredentialsPrefix + helperSuffix),
	}
}

// Store saves credentials into the store
func (ns *nativeStore) Store(serverAddress string, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}

// Erase removes credentials from the store for the given server
func (ns *nativeStore) Erase(serverAddress string) error {
	panic("not implemented") // TODO: Implement
}

// Get retrieves credentials from the store for the given server
func (ns *nativeStore) Get(serverAddress string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}
