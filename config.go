package credentials

import (
	"context"

	"oras.land/oras-go/v2/registry/remote/auth"
)

// StoreOptions provides options for GetConfiguredStore.
type StoreOptions struct {
	// PlainTextSave allows saving credentials in plain text in configuration file.
	PlainTextSave bool
}

// NewStore returns a new store from the settings in the configuration
// file.
func NewStore(configPath, serverAddress string, opts StoreOptions) Store {
	panic("not implemented") // TODO: Implement
}

// NewNStore returns a new store which will search credentials from the files
// specified by configPaths in order.
func NewNStore(configPaths []string, opts StoreOptions) Store {
	panic("not implemented") // TODO: Implement
}

func newNtore(configPaths ...string) *NStore {
	return &NStore{configPaths: configPaths}
}

type NStore struct {
	configPaths []string
}

func (s *NStore) Get(ctx context.Context, registry string) (auth.Credential, error) {
	for _, path := range s.configPaths {
		store := NewStore(path, registry, StoreOptions{})
		cred, err := store.Get(ctx, registry)
		if err != nil {
			panic(err)
		}
		if cred != auth.EmptyCredential {
			return cred, nil
		}
	}
	return auth.EmptyCredential, nil
}

func (s *NStore) Save(ctx context.Context, registry string, cred auth.Credential) error {
	store := NewStore(s.configPaths[0], registry, StoreOptions{})
	return store.Store(ctx, registry, cred)
}
