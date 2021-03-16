package buildkite

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "A data source to reference users by email.",
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Read: getUser,
	}
}

func getUser(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	email := d.Get("email").(string)
	u, err := bk.GetUser(email)
	if err != nil {
		return err
	}
	d.SetId(string(u.ID))
	d.Set("uuid", string(u.UUID))
	d.Set("name", string(u.Name))
	return nil
}
