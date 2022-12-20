package opencga

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func resourceStudy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStudyCreate,
		ReadContext:   resourceStudyRead,
		UpdateContext: resourceStudyUpdate,
		DeleteContext: resourceStudyDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The `id` of the project this study is associated with.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Study long name",
			},
			"alias": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Study alias name",
			},
			"description": &schema.Schema{
				Type:                  schema.TypeString,
				Required:              true,
				DiffSuppressFunc:      descriptionDiffSuppressFunc,
				DiffSuppressOnRefresh: true,
			},
			"check_description": &schema.Schema{
				Type:                  schema.TypeBool,
				Optional:              true,
				Default:               true,
				DiffSuppressFunc:      checkDescDiffSuppressFunc,
				DiffSuppressOnRefresh: true,
				Description:           "If true the description content will be checked against the state",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceStudyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	payload := map[string]interface{}{
		"name":        d.Get("name").(string),
		"alias":       d.Get("alias").(string),
		"description": d.Get("description").(string),
		"type":        "CASE_CONTROL",
	}

	if _, ok := d.GetOk("project"); !ok {
		return diag.Errorf("Must supply project id for study creation")
	}
	params := map[string]string{
		"projectId": d.Get("project").(string),
		"exclude":   "groups",
	}
	path := "studies/create"
	req, err := buildRequest(client, path, payload, params)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Call(req)
	if err != nil {
		return diag.FromErr(err)
	}
	var study Study
	err = mapstructure.Decode(resp.Results[0], &study)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(study.Id))
	resourceStudyRead(ctx, d, m)
	return diags
}

func resourceStudyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	path := fmt.Sprintf("studies/%s/info", d.Id())
	params := map[string]string{
		"include": "name,description,alias",
		"exclude": "groups",
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
		return diag.Errorf("Failed to find Study, got %d results", len(resp.Results))
	}
	var study Study
	err = mapstructure.Decode(resp.Results[0], &study)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", study.Name)
	d.Set("description", study.Description)
	d.Set("alias", study.Alias)
	return diags
}

func resourceStudyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceStudyRead(ctx, d, m)
}

func resourceStudyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("Pretending to delete but doing nothing....")
	return diags
}
