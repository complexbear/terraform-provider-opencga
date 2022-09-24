package opencga

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		// 		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Long name for the project, spaces are allowed",
			},
			"alias": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				StateFunc:   aliasStateFunc,
				Description: "Project alias. Do not supply the `null@` prefix seen in created resources. This will be added by OpenCGA automatically",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scientific_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Usually this is `Homo Sapiens`",
			},
			"taxonomy_code": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Usually this is `9606`",
			},
			"assembly": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Reference genome assembly name. i.e. GRCh38",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	payload := map[string]interface{}{
		"name":        d.Get("name").(string),
		"alias":       d.Get("alias").(string),
		"description": d.Get("description").(string),
		"organism": map[string]interface{}{
			"scientificName": d.Get("scientific_name").(string),
			"taxonomyCode":   d.Get("taxonomy_code").(int),
			"assembly":       d.Get("assembly").(string),
		},
	}
	path := "projects/create"
	req, err := buildRequest(client, path, payload, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Call(req)
	if err != nil {
		return diag.FromErr(err)
	}
	var project Project
	err = mapstructure.Decode(resp.Results[0], &project)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(project.Id))
	resourceProjectRead(ctx, d, m)
	return diags
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	path := fmt.Sprintf("projects/%s/info", d.Id())
	params := map[string]string{
		"include": "name,description,alias,organism",
		"exclude": "studies",
	}
	req, err := buildRequest(client, path, nil, params)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Call(req)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(resp.Results) != 1 {
		return diag.Errorf("Failed to find project, got %d results", len(resp.Results))
	}
	var project Project
	err = mapstructure.Decode(resp.Results[0], &project)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("alias", project.Alias)
	d.Set("scientific_name", project.Organism.ScientificName)
	d.Set("taxonomy_code", project.Organism.TaxonomyCode)
	d.Set("assembly", project.Organism.Assembly)
	return diags
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("Pretending to delete but doing nothing....")
	return diags
}

func aliasStateFunc(v interface{}) string {
	s := v.(string)
	if strings.Contains(s, "null@") {
		return s
	} else {
		return "null@" + s
	}
}
