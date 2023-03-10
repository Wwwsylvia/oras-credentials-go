package docker

import (
	"oras.land/oras-go/v2/registry/remote/auth"
)

type ConfigStore struct {
	configs []*ConfigFile
}

func GetCredentialStore(configPaths []string) (*ConfigStore, error) {
	panic("not implemented") // TODO: Implement
}

// Store saves credentials into the store
func (cs *ConfigStore) Store(serverAddress string, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}

// Erase removes credentials from the store for the given server
func (cs *ConfigStore) Erase(serverAddress string) error {
	panic("not implemented") // TODO: Implement
}

// Get retrieves credentials from the store for the given server
func (cs *ConfigStore) Get(serverAddress string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}
