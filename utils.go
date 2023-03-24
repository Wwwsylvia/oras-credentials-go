package credentials

import (
	"context"

	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

func Login(ctx context.Context, store Store, registry remote.Registry, cred auth.Credential) error {
	name := registry.Reference.Registry
	registry.Client = &auth.Client{
		Credential: auth.StaticCredential(name, cred),
	}

	if err := registry.Ping(ctx); err != nil {
		return err
	}

	return store.Store(ctx, name, cred)
}

func Logout(ctx context.Context, store Store, registryName string) error {
	return store.Erase(ctx, registryName)
}

func Credentials(store Store) func(context.Context, string) (auth.Credential, error) {
	return store.Get
}
