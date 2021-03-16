package buildkite

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

func testAccPipelineScheduleConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "buildkite_pipeline_schedule" "tfAccTestSchedule" {
	pipeline_id = "${buildkite_pipeline.test.id}"
	cronline = "0 0 1 1 *"
	env = {
		KEY1 = "val1"
	}
	enabled = true
	message = "test message"
	branch = "master"
	commit = "HEAD"
	label = "%s"
}`, testAccPipelineConfig(name), name)
}

func testAccPipelineScheduleUpdated(name string) string {
	return fmt.Sprintf(`
%s

resource "buildkite_pipeline_schedule" "tfAccTestSchedule" {
	pipeline_id = "${buildkite_pipeline.test.id}"
	cronline = "0 0 1 2 *"
	env = {
		KEY1 = "VAL1"
		KEY2 = "val2"
	}
	enabled = false
	message = "test message"
	branch = "test_branch"
	commit = "HEAD"
	label = "%s"
}
`, testAccPipelineConfig(name), name)
}

func TestAccPipelineSchedule(t *testing.T) {
	rLabel := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccPipelineScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPipelineScheduleConfig(rLabel),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPipelineScheduleExists("buildkite_pipeline_schedule.tfAccTestSchedule"),
					resource.TestCheckResourceAttrSet("buildkite_pipeline_schedule.tfAccTestSchedule", "id"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "cronline", "0 0 1 1 *"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "env.%", "1"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "env.KEY1", "val1"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "enabled", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "message", "test message"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "branch", "master"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "commit", "HEAD"),
				),
			},
			{
				Config: testAccPipelineScheduleUpdated(rLabel),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPipelineScheduleExists("buildkite_pipeline_schedule.tfAccTestSchedule"),
					resource.TestCheckResourceAttrSet("buildkite_pipeline_schedule.tfAccTestSchedule", "id"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "cronline", "0 0 1 2 *"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "env.%", "2"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "env.KEY1", "VAL1"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "env.KEY2", "val2"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "enabled", "false"),
					resource.TestCheckResourceAttr("buildkite_pipeline_schedule.tfAccTestSchedule", "branch", "test_branch"),
				),
			},
		},
	})
}

func testAccPipelineScheduleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource %s not found", name)
		}
		id := rs.Primary.ID
		if _, err := cli.ReadPipelineSchedule(id); err != nil {
			return err
		}
		return nil
	}
}

func testAccPipelineScheduleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buildkite_pipeline_schedule" {
			continue
		}
		toDelete := &client.PipelineSchedule{ID: graphql.String(rs.Primary.ID)}
		if err := cli.DeletePipelineSchedule(toDelete); err != nil {
			if !strings.Contains(err.Error(), "No schedule found") {
				return err
			}
		}
	}
	return nil
}

func TestFlattenMap(t *testing.T) {
	testCases := []struct {
		description string
		input       map[string]interface{}
		expected    []string
	}{
		{
			description: "standard map",
			input: map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
			},
			expected: []string{"key1=val1", "key2=val2"},
		},
		{
			description: "empty map",
			input:       map[string]interface{}{},
			expected:    []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := flattenMap(tc.input)

			sort.Strings(tc.expected)
			sort.Strings(actual)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestExpandSlice(t *testing.T) {
	testCases := []struct {
		description string
		input       []string
		shouldError bool
		output      map[string]interface{}
	}{
		{
			description: "legal input",
			input:       []string{"key1=val1", "key2=val2", "key3=val3"},
			shouldError: false,
			output: map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			},
		},
		{
			description: "empty input",
			input:       []string{},
			shouldError: false,
			output:      map[string]interface{}{},
		},
		{
			description: "nil input",
			input:       nil,
			shouldError: false,
			output:      map[string]interface{}{},
		},
		{
			description: "illegal input",
			input:       []string{"key1=val1", "key2=val2", "no equals sign"},
			shouldError: true,
			output:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := expandSlice(tc.input)
			if tc.shouldError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.output, actual)
		})
	}
}
