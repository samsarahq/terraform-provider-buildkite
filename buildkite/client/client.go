package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	buildkiteRest "github.com/buildkite/go-buildkite/v2/buildkite"
	"github.com/shurcooL/graphql"
)

// Base URLs for Buildkite API
const (
	RESTBaseURL = "https://api.buildkite.com/v2"
	GQLBaseURL  = "https://graphql.buildkite.com/v1"
)

// Client encapsulates the REST and GQL client for a given org.
type Client struct {
	// orgSlug is the slug of the org
	orgSlug string
	// orgID is the gql ID for the org
	orgID string

	httpClient *http.Client
	restClient *buildkiteRest.Client
	gqlClient  *graphql.Client
}

// We need this transport to add the authorization headers to our requests.
// We can't use `buildkite.NewTokenConfig` because that transport does not work
// for GQL requests.
type tokenTransport struct {
	token string
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return http.DefaultTransport.RoundTrip(req)
}

// NewClient returns a new buildkite client based on the given org and token.
// It will return an error if the token is invalid.
func NewClient(org, token string) (*Client, error) {
	httpClient := &http.Client{
		Transport: &tokenTransport{token: token},
	}
	restCli := buildkiteRest.NewClient(httpClient)
	gqlCli := graphql.NewClient(GQLBaseURL, httpClient)
	c := &Client{
		orgSlug:    org,
		httpClient: httpClient,
		restClient: restCli,
		gqlClient:  gqlCli,
	}
	if err := c.CheckAuth(); err != nil {
		return nil, fmt.Errorf("checking auth: %w", err)
	}

	// get org id
	var query struct {
		Organization struct {
			ID string `graphql:"id"`
		} `graphql:"organization(slug: $slug)"`
	}
	vars := map[string]interface{}{
		"slug": org,
	}
	err := c.gqlClient.Query(context.TODO(), &query, vars)
	if err != nil {
		return nil, fmt.Errorf("getting org id: %w", err)
	}
	c.orgID = query.Organization.ID
	return c, nil
}

// CheckAuth validates the client's token against the access token endpoint and
// returns whether the token appears valid based on this request.
func (c *Client) CheckAuth() error {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/access-token", RESTBaseURL), nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("non 200 status")
	}
	return nil
}
