package opencga

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a single Project for use in other resources",
		ReadContext: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			// Filter values
			"id_filter": &schema.Schema{
				Description: "A project id to limit the search",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"name_filter": &schema.Schema{
				Description: "A project name to limit the search",
				Type:        schema.TypeString,
				Optional:    true,
			},
			// Computed values
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
			"scientific_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"taxonomy_code": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"assembly": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	var path string
	if v, ok := d.GetOk("id_filter"); ok {
		// Exact search based on the id
		path = fmt.Sprintf("projects/%d/info", v.(int))
	} else {
		// Wide search optionally filtered by name
		path = "projects/search"
	}

	params := map[string]string{
		"include": "name,description,alias,organism",
		"exclude": "studies",
	}
	if v, ok := d.GetOk("name_filter"); ok {
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
	if len(resp.Results) == 0 {
		return diag.Errorf("Project '%s' not found", params["name"])
	}

	var project Project
	err = mapstructure.Decode(resp.Results[0], &project)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(project.Id))

	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("alias", project.Alias)
	d.Set("id", project.Id)
	d.Set("scientific_name", project.Organism.ScientificName)
	d.Set("taxonomy_code", project.Organism.TaxonomyCode)
	d.Set("assembly", project.Organism.Assembly)
	return diags
}
