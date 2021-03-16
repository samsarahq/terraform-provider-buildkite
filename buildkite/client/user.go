package client

import (
	"context"
	"errors"

	"github.com/shurcooL/graphql"
)

// User represents a Buildkite user.
type User struct {
	ID    graphql.String
	Name  graphql.String
	Email graphql.String
	UUID  graphql.String
}

// GetUser returns user from the GraphQL API by email.
func (c *Client) GetUser(email string) (*User, error) {
	var query struct {
		Organization struct {
			Members struct {
				Edges []struct {
					Node struct {
						User User
					}
				}
			} `graphql:"members(first: 1, email: $email)"`
		} `graphql:"organization(slug: $slug)"`
	}
	vars := map[string]interface{}{
		"email": graphql.String(email),
		"slug":  c.orgSlug,
	}
	err := c.gqlClient.Query(context.TODO(), &query, vars)
	if err != nil {
		return nil, err
	}
	edges := query.Organization.Members.Edges
	if len(edges) != 1 {
		return nil, errors.New("expected exactly 1 result")
	}
	return &edges[0].Node.User, nil
}
