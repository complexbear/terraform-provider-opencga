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

func dataSourceProjects() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to list several existing Projects for use in other resources",
		ReadContext: dataSourceProjectsRead,
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
			"projects": &schema.Schema{
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
				},
			},
		},
	}
}

func dataSourceProjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	projects := make([]Project, len(resp.Results))
	err = mapstructure.Decode(resp.Results, &projects)
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Printf("projects: %+v", projects)
	d.SetId(computeProjectsDataSourceId(d))

	// Store the project info in the resource data
	if err := d.Set("projects", flattenProjects(projects)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func computeProjectsDataSourceId(d *schema.ResourceData) string {
	// Create unique string representing the project search parameters
	var id strings.Builder
	if v, ok := d.GetOk("id_filter"); ok {
		id.WriteString(strconv.Itoa(v.(int)))
	}
	id.WriteRune('|')
	if v, ok := d.GetOk("name_filter"); ok {
		id.WriteString(v.(string))
	}
	return id.String()
}

func flattenProjects(projects []Project) []interface{} {
	result := make([]interface{}, len(projects))
	for i, p := range projects {
		r := make(map[string]interface{}, 0)
		r["id"] = p.Id
		r["name"] = p.Name
		r["description"] = p.Description
		r["alias"] = p.Alias
		r["scientific_name"] = p.Organism.ScientificName
		r["taxonomy_code"] = p.Organism.TaxonomyCode
		r["assembly"] = p.Organism.Assembly
		result[i] = r
	}
	return result
}
