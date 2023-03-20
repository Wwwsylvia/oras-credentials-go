package credentials

// GetStoreOptions is options for GetConfiguredStore.
type GetStoreOptions struct {
	// Path to the credential file
	CredentialsPath string
	// Disable saving credentials in plain text in configuration file.
	DisablePlainTextSave bool
}

// GetConfiguredStore returns a new store from the settings in the configuration
// file.
func GetConfiguredStore(configPath, serverAddress string, opts GetStoreOptions) Store {
	panic("not implemented") // TODO: Implement
}
