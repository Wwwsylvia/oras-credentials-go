package credentials

import (
	"context"

	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

// oras

type orasStore struct {
	configPaths []string
}

func NewOrasStore(configPaths ...string) *orasStore {
	return &orasStore{configPaths: configPaths}
}

func (s *orasStore) Get(registry string) (auth.Credential, error) {
	for _, path := range s.configPaths {
		store := GetConfiguredStore(path, registry, GetStoreOptions{})
		cred, err := store.Get(registry)
		if err != nil {
			panic(err)
		}
		if cred != auth.EmptyCredential {
			return cred, nil
		}
	}
	return auth.EmptyCredential, nil
}

func (s *orasStore) Save(registry string, cred auth.Credential) error {
	store := GetConfiguredStore(s.configPaths[0], registry, GetStoreOptions{})
	return store.Store(registry, cred)
}

func OrasLogin() {
	orasStore := NewOrasStore("config")
	regName := "registry"
	client := auth.Client{
		Credential: func(_ context.Context, s string) (auth.Credential, error) {
			return orasStore.Get(regName)
		},
	}
	reg, err := remote.NewRegistry(regName)
	if err != nil {
		panic(err)
	}
	reg.Client = &client
	ctx := context.Background()
	if err := reg.Ping(ctx); err != nil {
		orasStore.Save(regName, auth.Credential{})
	}
}

// notation
func GetNotationStore(configPath, credPath, registry string) Store {
	return GetConfiguredStore(configPath, registry, GetStoreOptions{
		CredentialsPath: credPath,
	})
}

// helm
