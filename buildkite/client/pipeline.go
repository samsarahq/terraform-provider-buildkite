package client

import (
	"context"
	"fmt"

	buildkiteRest "github.com/buildkite/go-buildkite/v2/buildkite"
	"github.com/shurcooL/graphql"
)

// Pipeline represents a pipeline in Buildkite.
type Pipeline = buildkiteRest.Pipeline

// GetPipelineID returns the gql ID for a given pipeline specified by its slug.
func (c *Client) GetPipelineID(slug string) (string, error) {
	var query struct {
		Pipeline struct {
			ID graphql.String `graphql:"id"`
		} `graphql:"pipeline(slug: $slug)"`
	}
	vars := map[string]interface{}{
		"slug": fmt.Sprintf("%s/%s", c.orgSlug, slug),
	}
	if err := c.gqlClient.Query(context.TODO(), &query, vars); err != nil {
		return "", err
	}
	return string(query.Pipeline.ID), nil
}

func (c *Client) CreatePipeline(pipeline *Pipeline) error {
	safeString := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}
	safeBool := func(b *bool) bool {
		if b == nil {
			return false
		}
		return *b
	}
	var provider buildkiteRest.ProviderSettings
	if pipeline.Provider != nil {
		provider = pipeline.Provider.Settings
	}
	payload := &buildkiteRest.CreatePipeline{
		// Steps are managed through the configuration field as a yaml string.
		Steps: nil,
		// Env is also part of the yaml steps.
		Env: nil,
		// Teams are associated through TeamPipelines rather than specifying them solely on
		// Pipeline create.
		TeamUuids: nil,

		Name:                            safeString(pipeline.Name),
		Repository:                      safeString(pipeline.Repository),
		DefaultBranch:                   safeString(pipeline.DefaultBranch),
		Description:                     safeString(pipeline.Description),
		BranchConfiguration:             safeString(pipeline.BranchConfiguration),
		SkipQueuedBranchBuilds:          safeBool(pipeline.SkipQueuedBranchBuilds),
		SkipQueuedBranchBuildsFilter:    safeString(pipeline.SkipQueuedBranchBuildsFilter),
		CancelRunningBranchBuilds:       safeBool(pipeline.CancelRunningBranchBuilds),
		CancelRunningBranchBuildsFilter: safeString(pipeline.CancelRunningBranchBuildsFilter),
		Configuration:                   pipeline.Configuration,

		ProviderSettings: provider,
	}
	p, _, err := c.restClient.Pipelines.Create(c.orgSlug, payload)
	if err != nil {
		return err
	}
	if p.ID == nil {
		return fmt.Errorf("nil ID for pipeline: %s", payload.Name)
	}
	pipeline.Slug = p.Slug
	return nil
}

func (c *Client) ReadPipeline(slug string) (*Pipeline, error) {
	p, _, err := c.restClient.Pipelines.Get(c.orgSlug, slug)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (c *Client) UpdatePipeline(pipeline *Pipeline) error {
	_, err := c.restClient.Pipelines.Update(c.orgSlug, pipeline)
	return err
}

func (c *Client) DeletePipeline(pipeline *Pipeline) error {
	_, err := c.restClient.Pipelines.Delete(c.orgSlug, *pipeline.Slug)
	return err
}
