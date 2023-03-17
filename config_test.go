package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"oras.land/oras-go/v2/registry/remote/auth"
)

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
