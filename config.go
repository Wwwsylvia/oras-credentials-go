package credentials

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
