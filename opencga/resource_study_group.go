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

func resourceStudyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStudyGroupCreate,
		ReadContext:   resourceStudyGroupRead,
		UpdateContext: resourceStudyGroupUpdate,
		DeleteContext: resourceStudyGroupDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"study": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceStudyGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	payload := map[string]interface{}{
		"id":   d.Get("name").(string),
		"name": d.Get("name").(string),
	}
	path := fmt.Sprintf("studies/%s/groups/create", d.Get("study"))
	req, err := buildRequest(client, path, payload, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Call(req)
	if err != nil {
		return diag.FromErr(err)
	}
	var studyGroup StudyGroup
	err = mapstructure.Decode(resp.Results[0], &studyGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(0))
	resourceStudyGroupRead(ctx, d, m)
	return diags
}

func resourceStudyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	path := fmt.Sprintf("studies/%s/groups", d.Get("study"))
	params := map[string]string{
		"name": d.Get("name").(string),
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
		return diag.Errorf("Failed to find study group, got %d results", len(resp.Results))
	}
	var studyGroup []StudyGroup
	err = mapstructure.Decode(resp.Results, &studyGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", studyGroup[0].Name)
	return diags
}

func resourceStudyGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceStudyGroupRead(ctx, d, m)
}

func resourceStudyGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("Pretending to delete but doing nothing....")
	return diags
}
