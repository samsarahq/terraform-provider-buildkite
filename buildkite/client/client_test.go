package client

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

const (
	orgEnvVar   = "BUILDKITE_ORGANIZATION_SLUG"
	tokenEnvVar = "BUILDKITE_TOKEN"
	userEnvVar  = "BUILDKITE_USER_EMAIL"
)

var (
	cli       *Client
	userEmail string
	userID    string
)

func init() {
	requiredVars := []string{orgEnvVar, tokenEnvVar, userEnvVar}
	for _, v := range requiredVars {
		if val := os.Getenv(v); val == "" {
			panic(fmt.Sprintf("Required env var %s is not set", v))
		}
	}

	c, err := NewClient(os.Getenv(orgEnvVar), os.Getenv(tokenEnvVar))
	if err != nil {
		panic("Couldn't create client")
	}
	cli = c
	userEmail = os.Getenv(userEnvVar)
	u, err := cli.GetUser(userEmail)
	if err != nil {
		panic("Couldn't get user")
	}
	userID = string(u.ID)
}

func TestCheckAuth(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{
			Transport: &tokenTransport{token: "wontwork"},
		},
	}
	if err := c.CheckAuth(); err == nil {
		t.Error("Invalid token still passed auth")
	}

	token := os.Getenv("BUILDKITE_TOKEN")
	c = &Client{
		httpClient: &http.Client{
			Transport: &tokenTransport{token: token},
		},
	}
	if err := c.CheckAuth(); err != nil {
		t.Errorf("Auth should have passed but failed, a valid token must be set at BUILDKITE_TOKEN")
	}
}
