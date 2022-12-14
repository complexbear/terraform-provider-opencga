package opencga

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func resourceFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFileCreate,
		ReadContext:   resourceFileRead,
		UpdateContext: resourceFileUpdate,
		DeleteContext: resourceFileDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"study": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The `id` of the study this file is associated with.",
			},
			"uri": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "File absolute path (URI), e.g. /genomes/sample/A00001.cram",
			},
			"path": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validatePathFunc,
				Description:      "Directory path, this does not have to be the absolute path if a root is configured. e.g. sample/, /genomes/sample",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Reference implementation:
	// https://gitlab.com/genomicsengland/bertha/core-bio-pipeline/-/blob/master/code/bertha-catalog/src/bertha/catalog/catalog.py

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	// OpenCGA has a race condition when linking files that share paths
	// Use a mutex to ensure sequential execution
	client.Mutex.Lock()
	defer client.Mutex.Unlock()

	payload := map[string]interface{}{
		"description":  "",
		"relatedFiles": []string{},
		"uri":          d.Get("uri").(string),
		"path":         d.Get("path").(string),
	}

	if _, ok := d.GetOk("study"); !ok {
		return diag.Errorf("Must supply study id for file linking")
	}
	params := map[string]string{
		"study":        d.Get("study").(string),
		"type":         "FILE",
		"parents":      "true",
		"createFolder": "false",
	}
	path := "files/link"
	req, err := buildRequest(client, path, payload, params)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Call(req)
	if err != nil {
		return diag.FromErr(err)
	}
	var file File
	err = mapstructure.Decode(resp.Results[0], &file)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(file.Id))
	resourceFileRead(ctx, d, m)
	return diags
}

func resourceFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	path := fmt.Sprintf("files/%s/info", d.Id())
	params := map[string]string{}
	req, err := buildRequest(client, path, nil, params)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Call(req)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(resp.Results) != 1 {
		return diag.Errorf("Failed to find File, got %d results", len(resp.Results))
	}
	var file File
	err = mapstructure.Decode(resp.Results[0], &file)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", file.Name)
	// Remove "file://"" that is added by OpenCGA in the response
	d.Set("uri", strings.Replace(file.Uri, "file://", "", 1))
	// Remove the file from the path, OpenCGA adds this in the response
	d.Set("path", filepath.Dir(file.Path)+"/")
	return diags
}

func resourceFileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// No updates supported
	return resourceFileRead(ctx, d, m)
}

func resourceFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("Pretending to delete but doing nothing....")
	return diags
}

func validatePathFunc(value interface{}, p cty.Path) diag.Diagnostics {
	path := value.(string)
	if strings.HasSuffix(path, "/") {
		return nil
	}
	return diag.Errorf("path parameter must be a dir path and end in /")
}
