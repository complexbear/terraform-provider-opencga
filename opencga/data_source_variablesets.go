package opencga

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func dataSourceVariableSets() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to list several existing VariableSets for use in other resources",
		ReadContext: dataSourceVariableSetsRead,
		Schema: map[string]*schema.Schema{
			// Filter values
			"study": &schema.Schema{
				Description: "A study id or name to limit the search",
				Type:        schema.TypeString,
				Required:    true,
			},
			// Computed values
			"variable_sets": &schema.Schema{
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
						"unique": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"variables": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVariableSetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	path := "variableset/search"
	params := map[string]string{
		"study":   d.Get("study").(string),
		"exclude": "variables",
	}

	req, err := buildRequest(client, path, nil, params)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Call(req)
	if err != nil {
		return diag.FromErr(err)
	}

	variable_sets := make([]VariableSet, len(resp.Results))
	err = mapstructure.Decode(resp.Results, &variable_sets)
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Printf("VariableSets: %+v", variable_sets)

	// Store the variable set info in the resource data
	if err := d.Set("variable_sets", flattenVariableSets(variable_sets)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(computeVariableSetsDataSourceId(d))

	return diags
}

func computeVariableSetsDataSourceId(d *schema.ResourceData) string {
	// Create unique string representing the variable set search parameters
	var id strings.Builder
	if v, ok := d.GetOk("study"); ok {
		id.WriteString(v.(string))
	}
	return id.String()
}

func flattenVariableSets(vs []VariableSet) []interface{} {
	variable_sets := make([]interface{}, len(vs))
	for i, v := range vs {
		json_string, _ := json.Marshal(v.Variables)
		item := make(map[string]interface{})
		item["id"] = v.Id
		item["name"] = v.Name
		item["description"] = v.Description
		item["unique"] = v.Unique
		item["variables"] = string(json_string)
		variable_sets[i] = item
	}
	return variable_sets
}
