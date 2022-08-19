package opencga

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func dataSourceStudies() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to list several existing Studies for use in other resources",
		ReadContext: dataSourceStudiesRead,
		Schema: map[string]*schema.Schema{
			// Filter values
			"id_filter": &schema.Schema{
				Description: "A study id to limit the search",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"alias_filter": &schema.Schema{
				Description: "A study alias to limit the search",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"project": &schema.Schema{
				Description: "A project id or alias to limit the search",
				Type:        schema.TypeString,
				Optional:    true,
			},
			// Computed values
			"studies": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"alias": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceStudiesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	params := map[string]string{
		"include": "name,description,alias",
		"exclude": "groups",
	}

	var path string
	if v, ok := d.GetOk("id_filter"); ok {
		// Exact search based on the id
		path = fmt.Sprintf("studies/%d/info", v.(int))
	} else {
		// Wide search optionally filtered by name but requires project reference
		if v, ok := d.GetOk("project"); ok {
			path = "studies/search"
			params["project"] = v.(string)
		} else {
			return diag.Errorf("Must supply project id or alias for study search")
		}

	}

	if v, ok := d.GetOk("alias_filter"); ok {
		params["name"] = v.(string)
	}

	req, err := buildRequest(client, path, nil, params)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Call(req)
	if err != nil {
		return diag.FromErr(err)
	}

	studies := make([]Study, len(resp.Results))
	err = mapstructure.Decode(resp.Results, &studies)
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Printf("Studies: %+v", studies)
	d.SetId(computeStudiesDataSourceId(d))

	// Store the study info in the resource data
	if err := d.Set("studies", studies); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func computeStudiesDataSourceId(d *schema.ResourceData) string {
	// Create unique string representing the project search parameters
	var id strings.Builder
	if v, ok := d.GetOk("id_filter"); ok {
		id.WriteString(strconv.Itoa(v.(int)))
	}
	id.WriteRune('|')
	if v, ok := d.GetOk("alias_filter"); ok {
		id.WriteString(v.(string))
	}
	return id.String()
}
