package credentials

import "oras.land/oras-go/v2/registry/remote/auth"

type CustomStore struct {
	Stores []Store // multiple file stores and native stores and custom stores?
}

type CustomStoreOptions struct {
	// TODO: Disable file store save?
}

func NewCustomStore(opts CustomStoreOptions) (Store, error) {
	return &CustomStore{}, nil
}

// Store saves credentials into the store
func (s *CustomStore) Store(serverAddress string, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}

// Erase removes credentials from the store for the given server
func (s *CustomStore) Erase(serverAddress string) error {
	panic("not implemented") // TODO: Implement
}

// Get retrieves credentials from the store for the given server
func (s *CustomStore) Get(serverAddress string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}
