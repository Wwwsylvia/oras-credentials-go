package credentials

type GetStoreOptions struct {
	NativeStoreOptions
	CredentialsPath string
}

// GetConfiguredStore returns a new store from the settings in the configuration
// file.
func GetConfiguredStore(configPath, serverAddress string, opts GetStoreOptions) (Store, error) {
	panic("not implemented") // TODO: Implement
}
