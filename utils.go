package credentials

import (
	"context"

	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

func Login(ctx context.Context, store Store, registry remote.Registry, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}

func Logout(ctx context.Context, store Store, registryName string) error {
	panic("not implemented") // TODO: Implement
}

func Credential(store Store) func(context.Context, string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}

func loginReg(ctx context.Context, store Store, registry remote.Registry, cred auth.Credential) error {
	name := registry.Reference.Registry
	registry.Client = &auth.Client{
		Credential: auth.StaticCredential(name, cred),
	}

	if err := registry.Ping(ctx); err != nil {
		return err
	}

	return store.Put(ctx, name, cred)
}

func logoutReg(ctx context.Context, store Store, registryName string) error {
	return store.Delete(ctx, registryName)
}

func credentials(store Store) func(context.Context, string) (auth.Credential, error) {
	return store.Get
}
