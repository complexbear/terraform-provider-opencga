package opencga

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func resourceStudyACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStudyACLCreate,
		ReadContext:   resourceStudyACLRead,
		UpdateContext: resourceStudyACLUpdate,
		DeleteContext: resourceStudyACLDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"member": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"template": &schema.Schema{
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

func resourceStudyACLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	payload := map[string]interface{}{
		"action":   "SET",
		"template": d.Get("template").(string),
		"study":    d.Get("study").(string),
	}
	path := fmt.Sprintf("studies/acl/%s/update", d.Get("member"))
	req, err := buildRequest(client, path, payload, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Call(req)
	if err != nil {
		return diag.FromErr(err)
	}
	var studyACL StudyACL
	err = mapstructure.Decode(resp.Results[0], &studyACL)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(studyACL.Member)
	resourceStudyACLRead(ctx, d, m)
	return diags
}

func resourceStudyACLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	path := fmt.Sprintf("studies/%s/acl", d.Get("study"))
	params := map[string]string{
		"member": d.Get("member").(string),
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
		return diag.Errorf("Failed to find study acl, got %d results", len(resp.Results))
	}
	var studyACL []StudyACL
	err = mapstructure.Decode(resp.Results, &studyACL)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceStudyACLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceStudyACLRead(ctx, d, m)
}

func resourceStudyACLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("Pretending to delete but doing nothing....")
	return diags
}
