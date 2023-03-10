package credentials

// Config ~/.docker/config.json file info
type Config struct {
	AuthConfigs       map[string]AuthConfig `json:"auths"`
	CredentialsStore  string                `json:"credsStore,omitempty"`
	CredentialHelpers map[string]string     `json:"credHelpers,omitempty"`
}

// AuthConfig contains authorization information for connecting to a Registry
type AuthConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Auth     string `json:"auth,omitempty"`

	// IdentityToken is used to authenticate the user and get
	// an access token for the registry.
	IdentityToken string `json:"identitytoken,omitempty"`

	// RegistryToken is a bearer token to be sent to a registry
	RegistryToken string `json:"registrytoken,omitempty"`
}
