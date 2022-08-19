package opencga

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
Config is a struct to hold global provider info
    Username: the username to log in with
    BaseUrl: the OpenCGA REST url, eg https://opencgainternal.test.aws.gel.ac/opencga/webservices
    Token: this will be computed during login and stored for use in further API calls
*/
type ProviderConfig struct {
	Username string
	BaseUrl  string
	Token    string
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		// This schema represents the parameters required to configure the provider
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENCGA_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OPENCGA_PASSWORD", nil),
			},
			"base_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"opencga_project":     resourceProject(),
			"opencga_study":       resourceStudy(),
			"opencga_study_acl":   resourceStudyACL(),
			"opencga_variableset": resourceVariableSet(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"opencga_project":      dataSourceProject(),
			"opencga_projects":     dataSourceProjects(),
			"opencga_studies":      dataSourceStudies(),
			"opencga_variablesets": dataSourceVariableSets(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	base_url := d.Get("base_url").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (username != "") && (password != "") && (base_url != "") {
		client := newClient(base_url)
		err := client.Login(username, password)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return client, diags
	}

	return nil, diags
}
