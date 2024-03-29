package artifacthub

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

type Config struct {
	ApiKey       string
	ApiKeySecret string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARTIFACTHUB_API_KEY", nil),
			},
			"api_key_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARTIFACTHUB_API_KEY_SECRET", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"artifacthub_user_webhook": resourceUserWebhook(),
			"artifacthub_org_webhook":  resourceOrgWebhook(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"artifacthub_helm_package": dataSourceHelmPackage(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	apiKey := d.Get("api_key").(string)
	apiKeySecret := d.Get("api_key_secret").(string)

	if apiKey == "" || apiKeySecret == "" {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No Credentials were provided",
		})
	}

	config := Config{
		ApiKey:       apiKey,
		ApiKeySecret: apiKeySecret,
	}

	return &config, diags
}
