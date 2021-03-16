package buildkite

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
	"github.com/shurcooL/graphql"
)

func resourcePipelineSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "A resource representing a pipeline schedule in Buildkite.",
		Schema: map[string]*schema.Schema{
			"pipeline_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cronline": {
				Type:     schema.TypeString,
				Required: true,
			},
			"env": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commit": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: createPipelineSchedule,
		Read:   readPipelineSchedule,
		Update: updatePipelineSchedule,
		Delete: deletePipelineSchedule,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// flattenMap is used to convert a map of string to string, to a slice of strings.
// Each element in the slice will have the format "key=value".
func flattenMap(m map[string]interface{}) []string {
	result := []string{}
	for k, v := range m {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}

// expandSlice takes a slice of string, with each element in the form of "key=value".
// It returns a map with the key as the key, and the value as the value.
// If an element contains an illegal format, an error is returned.
func expandSlice(s []string) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(s))
	for _, v := range s {
		segments := strings.Split(v, "=")
		if len(segments) != 2 {
			return result, fmt.Errorf("Illegal env entry - must be in the form of 'key=value': %s", v)
		}
		result[segments[0]] = segments[1]
	}
	return result, nil
}

func createPipelineSchedule(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)

	var env map[string]interface{}
	if v, ok := d.GetOk("env"); ok {
		env = v.(map[string]interface{})
	}

	ps := &client.PipelineSchedule{
		Branch:   graphql.String(d.Get("branch").(string)),
		Commit:   graphql.String(d.Get("commit").(string)),
		Cronline: graphql.String(d.Get("cronline").(string)),
		Enabled:  graphql.Boolean(d.Get("enabled").(bool)),
		Env:      flattenMap(env),
		Label:    graphql.String(d.Get("label").(string)),
		Message:  graphql.String(d.Get("message").(string)),
		Pipeline: struct {
			ID graphql.String
		}{
			ID: graphql.String(d.Get("pipeline_id").(string)),
		},
	}

	if err := bk.CreatePipelineSchedule(ps); err != nil {
		return err
	}
	d.SetId(string(ps.ID))
	return nil
}

func readPipelineSchedule(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	ps, err := bk.ReadPipelineSchedule(d.Id())
	if err != nil {
		return err
	}

	env := make([]string, 0)
	if ps.Env != nil {
		env = ps.Env
	}
	envMap, err := expandSlice(env)
	if err != nil {
		return err
	}

	d.Set("pipeline_id", ps.Pipeline.ID)
	d.Set("cronline", ps.Cronline)
	d.Set("env", envMap)
	d.Set("enabled", ps.Enabled)
	d.Set("message", ps.Message)
	d.Set("branch", ps.Branch)
	d.Set("commit", ps.Commit)
	d.Set("label", ps.Label)
	return nil
}

func updatePipelineSchedule(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)

	var env map[string]interface{}
	if v, ok := d.GetOk("env"); ok {
		env = v.(map[string]interface{})
	}

	ps := &client.PipelineSchedule{
		ID:       graphql.String(d.Id()),
		Branch:   graphql.String(d.Get("branch").(string)),
		Commit:   graphql.String(d.Get("commit").(string)),
		Cronline: graphql.String(d.Get("cronline").(string)),
		Enabled:  graphql.Boolean(d.Get("enabled").(bool)),
		Env:      flattenMap(env),
		Label:    graphql.String(d.Get("label").(string)),
		Message:  graphql.String(d.Get("message").(string)),
	}
	if err := bk.UpdatePipelineSchedule(ps); err != nil {
		return err
	}
	return nil
}

func deletePipelineSchedule(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	id := d.Id()
	if err := bk.DeletePipelineSchedule(&client.PipelineSchedule{ID: graphql.String(id)}); err != nil {
		return err
	}
	return nil
}
