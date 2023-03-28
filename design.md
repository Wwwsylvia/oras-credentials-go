# Design of `oras-credential-go`

## Goals

1. The library should be able to read credentials from Docker configuration file and external credentials store such as the native keychain of the operating system.
2. The library should be able to save credentials in Docker configuration file and external credentials store.
3. The library should be able to remove the specified credentials from external credential stores or Docker configuration file.
4. The library should be easy to integrate with Go projects that uses `oras-go v2`.
5. The library should be able to replace the following credential modules:
   - `oras`: https://github.com/oras-project/oras/tree/main/internal/credential
   - `oras-go` v1: https://github.com/oras-project/oras-go/tree/v1/pkg/auth/docker
   - `notation`: https://github.com/notaryproject/notation/tree/main/pkg/auth

Notes: The following is an example of Docker configuration file.

```json
{
    "auths": {
        "registry1.example.com": {
            "auth": "base64_encode(username:password)",
        },
        "registry2.example.com": {
            "identitytoken": "identity_token"
        },
        "registry3.example.com": {
            "registrytoken": "registry_token"
        },
        "registry4.example.com": {}
    },
    "credsStore": "desktop",
    "credHelpers": {
        "registry.example.com": "registryhelper",
        "awesomereg.example.org": "hip-star",
        "unicorn.example.io": "vcbait"
    }
}
```

## Non-Goals

1. The library will not support configuration formats other than the Docker configuration file.
2. The library will not handle encryption of credentials.

## Challenges

1. Depending on the version of the Docker CLI installed on the target machine, the format of the Docker configuration file may be different. The library should ensure that no config field is lost when saving credentials to the Docker configuration file.
2. Some clients like [Notation](https://github.com/notaryproject/notation) need to read plain-text credentials from Docker configuration file, but do not want to save credentials in plain-text in the configuration file.

## Proposal

### Solution to challenges

1. To ensure that no config field is lost when saving credentials to the Docker configuration file, we can first unmarshal the json file into a json object instead of a fixed struct when parsing the configuration file. And then we can make some changes to the `auths` field of the json object, and marshal the updated json object back to the file. That way we can keep all the unknown fields in the configuration file.
2. Provide an option to allow users to disable saving credentials in plain-text in configuration files. If the option is set, reading credentials from configuration files will be allowed but saving will result in no operation.

### Interfaces

We can define a basic interface for reading, saving and removing credentials as follows.

```go
package credentials

// Store is the interface that any credentials store must implement.
type Store interface {
	// Store saves credentials into the store
	Store(ctx context.Context, serverAddress string, cred auth.Credential) error
	// Erase removes credentials from the store for the given server
	Erase(ctx context.Context, serverAddress string) error
	// Get retrieves credentials from the store for the given server
	Get(ctx context.Context, serverAddress string) (auth.Credential, error)
}
```

The `auth.Credentials` refers to [`Credential`](https://pkg.go.dev/oras.land/oras-go/v2@v2.0.2/registry/remote/auth#Credential) defined in the `auth` package of `oras-go v2`.

### File Store

Based on the interface, we can further implement a `FileStore` for managing credentials stored in the Docker configuration file.

```go
package credentials

// FileStore implements a credentials store using the docker configuration file
// to keep the credentials in plain-text.
type FileStore struct {
	configPath  string
	DisableSave bool
}

// NewFileStore creates a new file credentials store.
func NewFileStore(configPath string) Store {
	return &FileStore{
		configPath: configPath,
	}
}

// Store saves credentials into the store
func (fs *FileStore) Store(ctx context.Context, serverAddress string, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}

// Erase removes credentials from the store for the given server
func (fs *FileStore) Erase(ctx context.Context, serverAddress string) error {
	panic("not implemented") // TODO: Implement
}

// Get retrieves credentials from the store for the given server
func (fs *FileStore) Get(ctx context.Context, serverAddress string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}
```

### Native Store

Besides, we can also implement a `NativeStore` for managing credentials using a native [credential store](https://docs.docker.com/engine/reference/commandline/login/#credentials-store) or [credential helpers](https://docs.docker.com/engine/reference/commandline/login/#credential-helpers).

```go
package credentials

const remoteCredentialsPrefix = "docker-credential-"

// NativeStore implements a credentials store
// using native keychain to keep credentials secure.
type NativeStore struct {
	programFunc client.ProgramFunc
}

// NewNativeStore creates a new native store that uses a remote helper program to
// manage credentials.
func NewNativeStore(helperSuffix string) Store {
	return &NativeStore{
		programFunc: client.NewShellProgramFunc(remoteCredentialsPrefix + helperSuffix),
	}
}

// Store saves credentials into the store
func (ns *NativeStore) Store(ctx context.Context, serverAddress string, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}

// Erase removes credentials from the store for the given server
func (ns *NativeStore) Erase(ctx context.Context, serverAddress string) error {
	panic("not implemented") // TODO: Implement
}

// Get retrieves credentials from the store for the given server
func (ns *NativeStore) Get(ctx context.Context, serverAddress string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}
```

The `client.ProgramFunc` refers to the [`ProgramFunc`](https://pkg.go.dev/github.com/docker/docker-credential-helpers@v0.7.0/client#ProgramFunc) defined in the package `client` of `docker-credential-helper`.

### Utility Methods

We can provide some common utility methods for convenience. The method names can be determined later.

#### NewStore()

This method is to return a new credential store based on the settings in the configuration file.  
The method should look for the credential store for a given server address in the order of credential helper, credential store and configuration file.  
The method should provide an option `PlainTextSave` to allow users to specify whether to save credentials in plain-text. When the native store is not available, if the option is set to false (default value), calling `NewStore().Save()` will result in no operation; if the option is set to true, calling `NewStore().Save()` will save the credential in plain-text in the configuration file.


```go
package credentials

// StoreOptions provides options for NewStore.
type StoreOptions struct {
	// PlainTextSave allows saving credentials in plain-text in configuration file.
	PlainTextSave bool
}

// NewStore returns a new store from the settings in the configuration
// file.
func NewStore(configPath, serverAddress string, opts StoreOptions) Store {
	panic("not implemented") // TODO: Implement
}
```

#### NewNStore()

This method is to return a new store which will search credentials from the files specified by the given config paths in order.  
There might be a better name for this method.

```go
package credentials

// NewNStore returns a new store which will search credentials from the files
// specified by configPaths in order.
func NewNStore(configPaths []string, opts StoreOptions) Store {
	panic("not implemented") // TODO: Implement
}
```

#### Login()

This method is to log a registry in.

```go
package credentials

func Login(ctx context.Context, store Store, registry remote.Registry, cred auth.Credential) error {
	panic("not implemented") // TODO: Implement
}
```

#### Logout()

This method is to log a registry out.

```go
package credentials

func Logout(ctx context.Context, store Store, registryName string) error {
	panic("not implemented") // TODO: Implement
}
```

#### Credential()

This method is to return a `Credential` function that can be used by [`auth.Client`](https://pkg.go.dev/oras.land/oras-go/v2@v2.0.2/registry/remote/auth#Client) of `oras-go v2`.

```go
package credentials

func Credential(store Store) func(context.Context, string) (auth.Credential, error) {
	panic("not implemented") // TODO: Implement
}
```

## References

- Project proposal: https://github.com/oras-project/oras-go/discussions/413
- #2
- [Docker Credentails Store](https://docs.docker.com/engine/reference/commandline/login/#credentials-store)
- [`docker/cli/config`](https://github.com/docker/cli/tree/master/cli/config)