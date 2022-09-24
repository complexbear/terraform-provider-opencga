package opencga

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-cty/cty"
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "This can be a user name or group id.",
			},
			"template": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateTemplate,
				Description:      "Preset permissions, can be one of: admin, analyst, view_only.",
			},
			"permissions": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comma separated list of OpenCGA permissions. Refer to OpenCGA docs for allowed values.",
			},
			"study": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the study that this ACL should be attached to.",
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
		"action": "SET",
		"study":  d.Get("study").(string),
	}

	// Template or permissions of the ACL
	template, template_ok := d.GetOk("template")
	permissions, permissions_ok := d.GetOk("permissions")
	if !template_ok && !permissions_ok {
		return diag.Errorf("Must provide either template or permissions")
	}
	if template_ok && permissions_ok {
		return diag.Errorf("Must provide either template or permissions but not both")
	}
	if template_ok {
		payload["template"] = template.(string)
	}
	if permissions_ok {
		payload["permissions"] = permissions.(string)
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
	return resourceStudyACLCreate(ctx, d, m)
}

func resourceStudyACLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("Pretending to delete but doing nothing....")
	return diags
}

func validateTemplate(v any, p cty.Path) diag.Diagnostics {
	// Check template is a valid name
	template := v.(string)
	supported_templates := []string{"admin", "analyst", "view_only"}
	for _, t := range supported_templates {
		if t == template {
			return nil
		}
	}
	return diag.Errorf("template must be one of %s, got: %s", supported_templates, template)
}
