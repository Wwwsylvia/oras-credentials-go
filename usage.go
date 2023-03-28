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

func NewOrasStore(configPaths ...string) *NStore {
	return &NStore{configPaths: configPaths}
}

func (s *orasStore) Get(ctx context.Context, registry string) (auth.Credential, error) {
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

func (s *orasStore) Save(ctx context.Context, registry string, cred auth.Credential) error {
	store := NewStore(s.configPaths[0], registry, StoreOptions{})
	return store.Put(ctx, registry, cred)
}

func OrasLogin() {
	orasStore := NewOrasStore("config")
	regName := "registry"

	client := auth.Client{
		Credential: func(ctx context.Context, s string) (auth.Credential, error) {
			return orasStore.Get(ctx, regName)
		},
	}

	ctx := context.Background()
	reg, err := remote.NewRegistry(regName)
	if err != nil {
		panic(err)
	}
	reg.Client = &client
	if err := reg.Ping(ctx); err != nil {
		orasStore.Save(ctx, regName, auth.Credential{})
	}
}

// notation
func GetNotationStore(configPath, credPath, registry string) Store {
	return NewStore(configPath, registry, StoreOptions{})
}

func login(registry, username, password, configPath string) error {
	ctx := context.Background()
	reg, err := remote.NewRegistry(registry)
	if err != nil {
		return err
	}
	cred := auth.Credential{
		Username: username,
		Password: password,
	}
	reg.Client = &auth.Client{
		Credential: auth.StaticCredential(registry, cred),
	}

	if err := reg.Ping(ctx); err != nil {
		return err
	}
	credStore := NewStore(configPath, registry, StoreOptions{
		PlainTextSave: true,
	})
	return credStore.Put(ctx, registry, cred)
}

func authenticate(registry, configPath string) error {
	credStore := NewStore(configPath, registry, StoreOptions{})
	reg, err := remote.NewRegistry(registry)
	if err != nil {
		return err
	}
	reg.Client = &auth.Client{
		Credential: func(ctx context.Context, s string) (auth.Credential, error) {
			return credStore.Get(ctx, registry)
		},
	}
	// do something with reg
	return nil
}

func logout(registry, configPath string) error {
	ctx := context.Background()
	credStore := NewStore(configPath, registry, StoreOptions{})
	return credStore.Delete(ctx, registry)
}

// helm
