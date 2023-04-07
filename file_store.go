/*
Copyright The ORAS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package credentials

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/oras-project/oras-credentials-go/internal/ioutils"
	"oras.land/oras-go/v2/registry/remote/auth"
)

// FileStore implements a credentials store using the docker configuration file
// to keep the credentials in plain-text.
type FileStore struct {
	// DisableSave disable saving credentials in plain text.
	// If DisableSave is set to true, Put() will return ErrPlainTextSaveDisabled.
	DisableSave   bool
	configPath    string
	content       map[string]interface{}
	contentRWLock sync.RWMutex
}

const (
	configFieldAuthConfigs   = "auths"
	configFieldUsername      = "username"
	configFieldPassword      = "password"
	configFieldBasicAuth     = "auth"
	configFieldIdentityToken = "identitytoken"
	configFieldRegistryToken = "registrytoken"
)

var (
	// ErrInvalidConfigFormat is returned when the config format is invalid.
	ErrInvalidConfigFormat = errors.New("invalid config format")
	// ErrPlainTextSaveDisabled is returned by Put() when DisableSave is set
	// to true.
	ErrPlainTextSaveDisabled = errors.New("plain text save is disabled")
)

// authConfig contains authorization information for connecting to a Registry
type authConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Auth     string `json:"auth,omitempty"`

	// IdentityToken is used to authenticate the user and get
	// an access token for the registry.
	IdentityToken string `json:"identitytoken,omitempty"`
	// RegistryToken is a bearer token to be sent to a registry
	RegistryToken string `json:"registrytoken,omitempty"`
}

// NewFileStore creates a new file credentials store.
func NewFileStore(configPath string) (*FileStore, error) {
	fs := &FileStore{configPath: configPath}
	configFile, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// init content map if the content file does not exist
			fs.content = make(map[string]interface{})
			return fs, nil
		}
		return nil, fmt.Errorf("failed to open config file at %s: %w", configPath, err)
	}
	defer configFile.Close()

	fi, err := configFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat config file at %s: %w", configPath, err)
	}
	if fi.IsDir() {
		return nil, fmt.Errorf("%s: configPath cannot be a directory", configPath)
	}

	// decode config content if the config file exists
	if err := json.NewDecoder(configFile).Decode(&fs.content); err != nil {
		return nil, fmt.Errorf("failed to decode config file at %s: %w", configPath, err)
	}
	return fs, nil
}

// Get retrieves credentials from the store for the given server address.
func (fs *FileStore) Get(_ context.Context, serverAddress string) (auth.Credential, error) {
	fs.contentRWLock.RLock()
	defer fs.contentRWLock.RUnlock()

	authCfg, ok := fs.getAuthConfig(serverAddress)
	if !ok {
		return auth.EmptyCredential, nil
	}

	cred := auth.Credential{
		Username:     authCfg.Username,
		Password:     authCfg.Password,
		RefreshToken: authCfg.IdentityToken,
		AccessToken:  authCfg.RegistryToken,
	}
	if authCfg.Auth != "" {
		var err error
		// override username and password
		cred.Username, cred.Password, err = decodeAuth(authCfg.Auth)
		if err != nil {
			return auth.EmptyCredential, fmt.Errorf("failed to decode username and password: %w: %v", ErrInvalidConfigFormat, err)
		}
	}
	return cred, nil
}

// Put saves credentials into the store for the given server address.
// Returns ErrPlainTextSaveDisabled if fs.DisableSave is set to true.
func (fs *FileStore) Put(_ context.Context, serverAddress string, cred auth.Credential) error {
	if fs.DisableSave {
		return ErrPlainTextSaveDisabled
	}

	fs.contentRWLock.Lock()
	defer fs.contentRWLock.Unlock()

	fs.updateAuths(serverAddress, cred)
	return fs.saveFile()
}

// Delete removes credentials from the store for the given server address.
func (fs *FileStore) Delete(ctx context.Context, serverAddress string) error {
	fs.contentRWLock.Lock()
	defer fs.contentRWLock.Unlock()

	if _, err := os.Stat(fs.configPath); os.IsNotExist(err) {
		// no ops if the config file does not exist
		return nil
	}
	authsMap, ok := fs.content[configFieldAuthConfigs].(map[string]interface{})
	if !ok {
		// no ops
		return nil
	}
	if _, ok = authsMap[serverAddress]; !ok {
		// no ops
		return nil
	}

	// update data
	delete(authsMap, serverAddress)
	fs.content[configFieldAuthConfigs] = authsMap
	return fs.saveFile()
}

// getAuthConfig reads the config and returns authConfig for serverAddress.
func (fs *FileStore) getAuthConfig(serverAddress string) (authConfig, bool) {
	authsMap, ok := fs.content[configFieldAuthConfigs].(map[string]interface{})
	if !ok {
		return authConfig{}, false
	}
	authConfigObj, ok := authsMap[serverAddress].(map[string]interface{})
	if !ok {
		return authConfig{}, false
	}

	var authCfg authConfig
	for k, v := range authConfigObj {
		switch k {
		case configFieldUsername:
			authCfg.Username, _ = v.(string)
		case configFieldPassword:
			authCfg.Password, _ = v.(string)
		case configFieldBasicAuth:
			authCfg.Auth, _ = v.(string)
		case configFieldIdentityToken:
			authCfg.IdentityToken, _ = v.(string)
		case configFieldRegistryToken:
			authCfg.RegistryToken, _ = v.(string)
		}
	}
	return authCfg, true
}

// updateAuths updates the Auths field of fs.content based on cred.
func (fs *FileStore) updateAuths(serverAddress string, cred auth.Credential) {
	authsMap, ok := fs.content[configFieldAuthConfigs].(map[string]interface{})
	if !ok {
		authsMap = make(map[string]interface{})
	}
	authCfg, ok := authsMap[serverAddress].(map[string]interface{})
	if !ok {
		authCfg = make(map[string]interface{})
	}
	authCfg[configFieldBasicAuth] = encodeAuth(cred.Username, cred.Password)
	authCfg[configFieldUsername] = ""
	authCfg[configFieldPassword] = ""
	authCfg[configFieldIdentityToken] = cred.RefreshToken
	authCfg[configFieldRegistryToken] = cred.AccessToken

	// update data
	authsMap[serverAddress] = authCfg
	fs.content[configFieldAuthConfigs] = authsMap
}

// saveFile saves fs.content into fs.configPath.
func (fs *FileStore) saveFile() error {
	jsonData, err := json.MarshalIndent(fs.content, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	dir := filepath.Dir(fs.configPath)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("failed to make directory %s: %w", dir, err)
	}
	ingest, err := ioutils.Ingest(bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}

	// handle symlink
	targetPath := fs.configPath
	if link, err := os.Readlink(fs.configPath); err == nil {
		targetPath = link
	}
	// copy file with original ownership and permissions
	ioutils.CopyFilePermissions(targetPath, ingest)
	if err := os.Rename(ingest, targetPath); err != nil {
		// clean up the ingest file
		os.Remove(ingest)
		return fmt.Errorf("failed to save config file: %w", err)
	}
	return nil
}

// encodeAuth base64-encodes username and password into base64(username:password).
func encodeAuth(username, password string) string {
	if username == "" && password == "" {
		return ""
	}

	authStr := username + ":" + password
	msg := []byte(authStr)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(encoded, msg)
	return string(encoded)
}

// decodeAuth decodes a base64 encoded string and returns username and password.
func decodeAuth(authStr string) (username string, password string, err error) {
	if authStr == "" {
		return "", "", nil
	}

	decodedLen := base64.StdEncoding.DecodedLen(len(authStr))
	decoded := make([]byte, decodedLen)
	authByte := []byte(authStr)
	n, err := base64.StdEncoding.Decode(decoded, authByte)
	if err != nil {
		return "", "", err
	}
	if n > decodedLen {
		return "", "", errors.New("size mismatch")
	}
	arr := strings.SplitN(string(decoded), ":", 2)
	if len(arr) != 2 {
		return "", "", errors.New("auth does not conform username:password format")
	}
	password = strings.Trim(arr[1], "\x00")
	return arr[0], password, nil
}
