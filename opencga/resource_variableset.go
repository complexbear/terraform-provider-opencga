package opencga

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/mapstructure"
)

var required_variable_attrs = []string{"allowedValues", "description", "multiValue", "name", "required", "title", "type"}

func resourceVariableSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVariableSetCreate,
		ReadContext:   resourceVariableSetRead,
		UpdateContext: resourceVariableSetUpdate,
		DeleteContext: resourceVariableSetDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"study": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"unique": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"variables": &schema.Schema{
				Type:                  schema.TypeString,
				Required:              true,
				ValidateFunc:          validation.StringIsJSON,
				DiffSuppressFunc:      variableDiffSuppressFunc,
				DiffSuppressOnRefresh: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVariableSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	// Convert variables json string into json data struct
	var variables_json []interface{}
	json_data := []byte(d.Get("variables").(string))
	err := json.Unmarshal(json_data, &variables_json)
	if err != nil {
		return diag.Errorf("Unable to convert variable string to json")
	}

	payload := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"unique":      d.Get("unique").(bool),
		"variables":   variables_json,
	}
	params := map[string]string{
		"study": d.Get("study").(string),
	}
	path := "variableset/create"
	req, err := buildRequest(client, path, payload, params)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Call(req)
	if err != nil {
		return diag.FromErr(err)
	}
	var variable_set VariableSet
	err = mapstructure.Decode(resp.Results[0], &variable_set)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(variable_set.Id))
	resourceVariableSetRead(ctx, d, m)
	return diags
}

func resourceVariableSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := m.(*APIClient)

	path := fmt.Sprintf("variableset/%s/info", d.Id())
	params := map[string]string{
		"study": d.Get("study").(string),
	}
	req, err := buildRequest(client, path, nil, params)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Call(req)
	if err != nil {
		return diag.Errorf("Failed to decode result: %s", err)
	}
	if len(resp.Results) != 1 {
		return diag.Errorf("Failed to find VariableSets, got %d results", len(resp.Results))
	}

	var variable_set VariableSet
	err = mapstructure.Decode(resp.Results[0], &variable_set)
	if err != nil {
		return diag.FromErr(err)
	}

	// convert variables json data struct to string for schema
	variables_string, err := json.Marshal(variable_set.Variables)
	if err != nil {
		fmt.Println("Failed to marshall variable data")
	}

	d.Set("name", variable_set.Name)
	d.Set("description", variable_set.Description)
	d.Set("unique", variable_set.Unique)
	d.Set("variables", string(variables_string))
	return diags
}

func resourceVariableSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceVariableSetRead(ctx, d, m)
}

func resourceVariableSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	fmt.Println("Pretending to delete but doing nothing....")
	return diags
}

func variableDiffSuppressFunc(k, oldValue, newValue string, d *schema.ResourceData) bool {
	// Compare entries in old and new lists to check equivalence

	// Convert json strings to list of maps
	var oldAttrMap []map[string]interface{}
	json.Unmarshal([]byte(oldValue), &oldAttrMap)
	var newAttrMap []map[string]interface{}
	json.Unmarshal([]byte(newValue), &newAttrMap)

	// Must have equal item counts
	if len(oldAttrMap) != len(newAttrMap) {
		fmt.Printf(
			"Mismatched variable set counts. Old:%d, New:%d",
			len(oldAttrMap),
			len(newAttrMap),
		)
		return false
	}

	// Check names of items match
	oldNames := make([]string, len(oldAttrMap))
	newNames := make([]string, len(newAttrMap))
	for i := 0; i < len(oldAttrMap); i++ {
		oldNames[i] = oldAttrMap[i]["name"].(string)
		newNames[i] = newAttrMap[i]["name"].(string)
	}
	sort.Strings(oldNames)
	sort.Strings(newNames)
	for i, n := range oldNames {
		if n != newNames[i] {
			fmt.Printf(
				"Mismatched variable set names. Old:%v, New:%v",
				oldNames,
				newNames,
			)
			return false
		}
	}

	// TODO - check further the other attributes of the variable setss
	return true
}
