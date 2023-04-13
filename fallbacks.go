package credentials

import (
	"context"

	"oras.land/oras-go/v2/registry/remote/auth"
)

type StoreWithFallbacks struct {
	stores []Store
}

// NewStoreWithFallbacks returns a new store based on the given stores.
// The second and the subsequent stores will be used as fallbacks for the first store.
func NewStoreWithFallbacks(store Store, fallbacks ...Store) Store {
	return &StoreWithFallbacks{
		stores: append([]Store{store}, fallbacks...),
	}
}

// Get retrieves credentials from the store for the given server address.
func (swf *StoreWithFallbacks) Get(ctx context.Context, serverAddress string) (auth.Credential, error) {
	for _, store := range swf.stores {
		cred, err := store.Get(ctx, serverAddress)
		if err != nil {
			return auth.EmptyCredential, err
		}
		if cred != auth.EmptyCredential {
			return cred, nil
		}
	}
	return auth.EmptyCredential, nil
}

// Put saves credentials into the store for the given server address.
func (swf *StoreWithFallbacks) Put(ctx context.Context, serverAddress string, cred auth.Credential) error {
	return swf.stores[0].Put(ctx, serverAddress, cred)
}

// Delete removes credentials from the store for the given server address.
func (swf *StoreWithFallbacks) Delete(ctx context.Context, serverAddress string) error {
	return swf.stores[0].Delete(ctx, serverAddress)
}
