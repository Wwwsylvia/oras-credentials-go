package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"oras.land/oras-go/v2/registry/remote/auth"
)

type config struct {
	path string
	data map[string]interface{}
}

func (c *config) Get(hostName string) auth.Credential {
	f, err := os.Open(c.path)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	jsonObj, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(jsonObj, &c.data)
	if err != nil {
		fmt.Println(err)
	}

	authItem := c.data["auths"].(map[string]interface{})
	regItem := authItem[hostName].(map[string]interface{})
	fmt.Println("registry:", hostName, " item:", regItem)
	fmt.Println("registry:", hostName, "auth:", regItem["auth"])

	cred := auth.Credential{}
	if token, ok := regItem["identitytoken"]; ok {
		tokenStr := token.(string)
		cred.RefreshToken = tokenStr

	}
	return cred
}

func (c *config) Save(hostName string, cred auth.Credential) {
	authItem := c.data["auths"].(map[string]interface{})
	regItem := authItem[hostName].(map[string]interface{})
	regItem["auth"] = cred.Username + ":" + cred.Password

	authItem[hostName] = regItem
	c.data["auths"] = authItem

	updatedData, err := json.MarshalIndent(c.data, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(c.path, updatedData, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGet(t *testing.T) {
	config := &config{path: "config.test.json"}
	config.Get("myacr.azurecr.io")
	config.Get("myregistry.example.com")

	config.Save("myregistry.example.com", auth.Credential{Username: "abc", Password: "123"})
}

func TestFileStore_Get(t *testing.T) {
	ctx := context.Background()
	fs, err := NewFileStore("testdata/config.test.json")
	if err != nil {
		t.Fatal("NewFileStore() error =", err)
	}

	tests := []struct {
		name          string
		serverAddress string
		want          auth.Credential
		wantErr       bool
	}{
		{
			name:          "username + password",
			serverAddress: "registry1.example.com",
			want: auth.Credential{
				Username: "username",
				Password: "password",
			},
		},
		{
			name:          "identity token",
			serverAddress: "registry2.example.com",
			want: auth.Credential{
				RefreshToken: "identity_token",
			},
		},
		{
			name:          "registry token",
			serverAddress: "registry3.example.com",
			want: auth.Credential{
				AccessToken: "registry_token",
			},
		},
		{
			name:          "empty",
			serverAddress: "registry4.example.com",
			want:          auth.EmptyCredential,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fs.Get(ctx, tt.serverAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileStore.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: how to test?
func TestFileStore_Put(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.test.json")

	fs, err := NewFileStore(configPath)
	if err != nil {
		panic(err)
	}

	reg := "test0331.test.com"
	cred := auth.Credential{
		Username:     "username",
		Password:     "password",
		RefreshToken: "refresh_token",
		AccessToken:  "access_token",
	}

	ctx := context.Background()
	if err := fs.Put(ctx, reg, cred); err != nil {
		panic(err)
	}
}

func TestFileStore_Delete(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.test.json")

	fs, err := NewFileStore(configPath)
	if err != nil {
		panic(err)
	}

	reg1 := "test1.test.com"
	cred1 := auth.Credential{
		Username:     "username",
		Password:     "password",
		RefreshToken: "refresh_token",
		AccessToken:  "access_token",
	}

	ctx := context.Background()
	if err := fs.Put(ctx, reg1, cred1); err != nil {
		panic(err)
	}

	reg2 := "test2.test.com"
	cred2 := auth.Credential{
		Username:     "username2",
		Password:     "password2",
		RefreshToken: "refresh_token2",
		AccessToken:  "access_token2",
	}
	// TODO: WRONG! reg1 is overwritten
	if err := fs.Put(ctx, reg2, cred2); err != nil {
		panic(err)
	}

	if err := fs.Delete(ctx, reg1); err != nil {
		panic(err)
	}
}