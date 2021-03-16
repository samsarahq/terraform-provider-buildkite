package buildkite

import (
	"strconv"

	buildkiteRest "github.com/buildkite/go-buildkite/v2/buildkite"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Description: "A resource representing a pipeline in Buildkite.",
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"steps": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"branch_configuration": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cancel_running_branch_builds": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cancel_running_branch_builds_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_branch": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"skip_queued_branch_builds": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"skip_queued_branch_builds_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_settings": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
		Create: createPipeline,
		Read:   readPipeline,
		Update: updatePipeline,
		Delete: deletePipeline,
		Importer: &schema.ResourceImporter{
			State: importPipeline,
		},
	}
}

func strPtr(s interface{}) *string {
	if s == nil {
		return strPtr("")
	}
	str, ok := s.(string)
	if !ok {
		return strPtr("")
	}
	return &str
}

func boolPtr(b interface{}) *bool {
	if b == nil {
		return boolPtr(false)
	}
	boo, ok := b.(bool)
	if !ok {
		return boolPtr(false)
	}
	return &boo
}

func pipelineFromSchema(d *schema.ResourceData) *client.Pipeline {
	boolStrPtr := func(b interface{}) *bool {
		if b == nil {
			return boolPtr(false)
		}
		s, ok := b.(string)
		if !ok {
			return boolPtr(false)
		}
		boo, _ := strconv.ParseBool(s)
		return &boo
	}
	// Always use GithubSettings as they are a superset of all settings.
	var provider *buildkiteRest.GitHubSettings
	if _, ok := d.GetOk("provider_settings"); ok {
		settings := d.Get("provider_settings").(map[string]interface{})
		provider = &buildkiteRest.GitHubSettings{
			TriggerMode:                             strPtr(settings["trigger_mode"]),
			BuildPullRequests:                       boolStrPtr(settings["build_pull_requests"]),
			PullRequestBranchFilterEnabled:          boolStrPtr(settings["pull_request_branch_filter_enabled"]),
			PullRequestBranchFilterConfiguration:    strPtr(settings["pull_request_branch_filter_configuration"]),
			SkipPullRequestBuildsForExistingCommits: boolStrPtr(settings["skip_pull_request_builds_for_existing_commits"]),
			BuildPullRequestForks:                   boolStrPtr(settings["build_pull_request_forks"]),
			PrefixPullRequestForkBranchNames:        boolStrPtr(settings["prefix_pull_request_fork_branch_names"]),
			BuildTags:                               boolStrPtr(settings["build_tags"]),
			PublishCommitStatus:                     boolStrPtr(settings["publish_commit_status"]),
			PublishCommitStatusPerStep:              boolStrPtr(settings["publish_commit_status_per_step"]),
			FilterEnabled:                           boolStrPtr(settings["filter_enabled"]),
			FilterCondition:                         strPtr(settings["filter_condition"]),
			SeparatePullRequestStatuses:             boolStrPtr(settings["separate_pull_request_statuses"]),
			PublishBlockedAsPending:                 boolStrPtr(settings["publish_blocked_as_pending"]),
		}
	}

	return &client.Pipeline{
		Name:                            strPtr(d.Get("name")),
		Slug:                            strPtr(d.Get("slug")),
		Repository:                      strPtr(d.Get("repository")),
		Steps:                           nil,
		Configuration:                   d.Get("steps").(string), // YAML steps specified here.
		DefaultBranch:                   strPtr(d.Get("default_branch")),
		Description:                     strPtr(d.Get("description")),
		BranchConfiguration:             strPtr(d.Get("branch_configuration")),
		SkipQueuedBranchBuilds:          boolPtr(d.Get("skip_queued_branch_builds")),
		SkipQueuedBranchBuildsFilter:    strPtr(d.Get("skip_queued_branch_builds_filter")),
		CancelRunningBranchBuilds:       boolPtr(d.Get("cancel_running_branch_builds")),
		CancelRunningBranchBuildsFilter: strPtr(d.Get("cancel_running_branch_builds_filter")),

		Provider: &buildkiteRest.Provider{
			Settings: provider,
		},
	}
}
func createPipeline(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	p := pipelineFromSchema(d)
	if err := bk.CreatePipeline(p); err != nil {
		return err
	}
	d.Set("slug", p.Slug)
	return readPipeline(d, m)
}

func readPipeline(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	slug := d.Get("slug").(string)
	p, err := bk.ReadPipeline(slug)
	if err != nil {
		return err
	}

	// Set the ID to the gql ID so it can be used by other resources.
	id, err := bk.GetPipelineID(slug)
	if err != nil {
		return err
	}

	d.SetId(id)
	// Terraform handles pointers gracefully.
	d.Set("repository", p.Repository)
	d.Set("steps", p.Configuration)
	d.Set("branch_configuration", p.BranchConfiguration)
	d.Set("cancel_running_branch_builds", p.CancelRunningBranchBuilds)
	d.Set("cancel_running_branch_builds_filter", p.CancelRunningBranchBuildsFilter)
	d.Set("default_branch", p.DefaultBranch)
	d.Set("description", p.Description)
	d.Set("skip_queued_branch_builds", p.SkipQueuedBranchBuilds)
	d.Set("skip_queued_branch_builds_filter", p.SkipQueuedBranchBuildsFilter)

	boolPtrToStr := func(b *bool) string {
		if b == nil {
			return ""
		}
		return strconv.FormatBool(*b)
	}
	safeString := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}
	settings := p.Provider.Settings.(*buildkiteRest.GitHubSettings)
	provider := map[string]interface{}{
		"trigger_mode":                                  safeString(settings.TriggerMode),
		"build_pull_requests":                           boolPtrToStr(settings.BuildPullRequests),
		"pull_request_branch_filter_enabled":            boolPtrToStr(settings.PullRequestBranchFilterEnabled),
		"pull_request_branch_filter_configuration":      safeString(settings.PullRequestBranchFilterConfiguration),
		"skip_pull_request_builds_for_existing_commits": boolPtrToStr(settings.SkipPullRequestBuildsForExistingCommits),
		"build_pull_request_forks":                      boolPtrToStr(settings.BuildPullRequestForks),
		"prefix_pull_request_fork_branch_names":         boolPtrToStr(settings.PrefixPullRequestForkBranchNames),
		"build_tags":                                    boolPtrToStr(settings.BuildTags),
		"publish_commit_status":                         boolPtrToStr(settings.PublishCommitStatus),
		"publish_commit_status_per_step":                boolPtrToStr(settings.PublishCommitStatusPerStep),
		"filter_enabled":                                boolPtrToStr(settings.FilterEnabled),
		"filter_condition":                              safeString(settings.FilterCondition),
		"separate_pull_request_statuses":                boolPtrToStr(settings.SeparatePullRequestStatuses),
		"publish_blocked_as_pending":                    boolPtrToStr(settings.PublishBlockedAsPending),
	}

	// Delete nil values.
	for k, v := range provider {
		val := v.(string)
		if val == "" {
			delete(provider, k)
		}
	}
	d.Set("provider_settings", provider)
	return nil
}

func updatePipeline(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	if err := bk.UpdatePipeline(pipelineFromSchema(d)); err != nil {
		return err
	}
	return readPipeline(d, m)
}

func deletePipeline(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	if err := bk.DeletePipeline(pipelineFromSchema(d)); err != nil {
		return err
	}
	return nil
}

func importPipeline(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	d.Set("name", d.Id())
	if err := readPipeline(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
