# Design

## Goals

1. The library should be able to read credentials from Docker configuration file and external credentials stores such as the native keychain of the operating system.
2. The library should be able to save credentials in Docker configuration file and external credential stores. Besides, it should provide options to disable saving credentials in plain text.
3. The library should be able to remove the specified credentials from external credential stores or Docker configuration file.
4. The library should be easy to integrate with Go projects that uses `oras-go v2`.
5. The library should be able to replace the following credential modules:
   - `oras`: https://github.com/oras-project/oras/tree/main/internal/credential
   - `oras-go` v1: https://github.com/oras-project/oras-go/tree/v1/pkg/auth/docker
   - `notation`: https://github.com/notaryproject/notation/tree/main/pkg/auth

Notes: An example of Docker configuration file is shown below.

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
2. The credential helper configuration (`credsStore` and `credHelpers`) and the credential itself (`auths`) may be maintained in different files. (Some CLIs like notation have its own configuration file)

## Proposal

### Solution to challenges

1. To be able to update the `auths` field without touching other fields in the configuration file, we can unmarshal the json file into a json object instead of a fixed struct. After we update the `auths` field, we can save the updated json object back to fhe file.
2. We should allow users to pass the path to the credential helper configuration file and the path to the auth configuration file separately.

### Interface

We can define a basic interface as below for reading, saving and removing credentials.

```go
// Store is the interface that any credentials store must implement.
type Store interface {
    // Store saves credentials into the store
    Store(serverAddress string, cred auth.Credential) error
    // Erase removes credentials from the store for the given server
    Erase(serverAddress string) error
    // Get retrieves credentials from the store for the given server
    Get(serverAddress string) (auth.Credential, error)
}
```

The `auth.Credentials` refers to `Credential` defined in the `auth` package in `oras-go v2`.

Based on the interface, we can further implement a `FileStore` for managing credentials stored in the Docker configuration file.

In addition, we can also implement a `NativeStore` for managing credentials using native credential helpers like `pass`, `wincred`, `osxkeychain`, etc.

### Utility Methods

### Use cases

## References

- [Docker Credentails Store](https://docs.docker.com/engine/reference/commandline/login/#credentials-store)