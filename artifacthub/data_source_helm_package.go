package artifacthub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHelmPackage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHelmPackageRead,
		Schema: map[string]*schema.Schema{
			"repo_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed Fields

			"package_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceHelmPackageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	version := d.Get("version").(string)
	repoName := d.Get("repo_name").(string)
	name := d.Get("name").(string)

	var req *http.Request
	var err error

	if version != "" {
		req, err = http.NewRequest("GET", fmt.Sprintf("https://artifacthub.io/api/v1/packages/helm/%s/%s/%s", repoName, name, version), nil)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		req, err = http.NewRequest("GET", fmt.Sprintf("https://artifacthub.io/api/v1/packages/helm/%s/%s", repoName, name), nil)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	req.Header.Add("accept", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	pkg := make(map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&pkg)
	if err != nil {
		return diag.FromErr(err)
	}

	if version == "" {
		if err := d.Set("version", pkg["version"]); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("package_id", pkg["package_id"]); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
