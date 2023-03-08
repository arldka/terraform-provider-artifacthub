package artifacthub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrgWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrgWebhookCreate,
		ReadContext:   resourceOrgWebhookRead,
		UpdateContext: resourceOrgWebhookUpdate,
		DeleteContext: resourceOrgWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the webhook. ",
				Required:    true,
			},
			"org_name": {
				Type:        schema.TypeString,
				Description: "Name of the organization. ",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description for the webhook. ",
				Optional:    true,
				Default:     "",
			},
			"url": {
				Type:        schema.TypeString,
				Description: "Webhook target url. ",
				Required:    true,
			},
			"secret": {
				Type:        schema.TypeString,
				Description: "Webhook secret for basic authentication. ",
				Optional:    true,
				Default:     "",
			},
			"content_type": {
				Type:        schema.TypeString,
				Description: "Content Type of the webhook request body. ",
				Optional:    true,
				Default:     "application/json",
			},
			"template": {
				Type:        schema.TypeString,
				Description: "Template of the webhook request body. ",
				Optional:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Status of the webhook. ",
				Required:    true,
			},
			"event_kinds": {
				Type:        schema.TypeList,
				Description: "Event Kinds of the webhook. `0` for new package release, `1` for security alerts, `2` for Repository tracking errors, `4` for repository scanning errors. ",
				Required:    true,
				Elem:        schema.TypeInt,
			},
			"packages": {
				Type:        schema.TypeList,
				Description: "Packages list",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"package_id": {
							Type:        schema.TypeString,
							Description: "package_id (can be recovered from a data source). ",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceOrgWebhookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var req *http.Request
	var err error

	webhookData := map[string]any{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"url":          d.Get("url").(string),
		"secret":       d.Get("secret").(string),
		"content_type": d.Get("content_type").(string),
		"active":       d.Get("active").(bool),
		"event_kinds":  d.Get("event_kinds").([]int),
		"packages":     d.Get("packages").([]map[string]string),
	}

	tpl := d.Get("template").(string)

	if tpl != "" {
		webhookData["template"] = tpl
	}

	data, _ := json.Marshal(webhookData)

	req, err = http.NewRequest("POST", fmt.Sprintf("https://artifacthub.io/api/v1/webhooks/org/%s", d.Get("org_name").(string)), bytes.NewBuffer(data))

	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("X-API-KEY-ID", d.Get("api_key").(string))
	req.Header.Add("X-API-KEY-SECRET", d.Get("api_key_secret").(string))

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	resp := make(map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return diag.FromErr(err)
	}

	if r.Status != "201" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Org Webhook",
			Detail:   resp["message"].(string),
		})
		return diags
	}

	verifResp := make([]map[string]interface{}, 0)

	verifReq, err := http.NewRequest("GET", fmt.Sprintf("https://artifacthub.io/api/v1/webhooks/org/%s", d.Get("org_name").(string)), nil)

	if err != nil {
		return diag.FromErr(err)
	}

	verifReq.Header.Add("X-API-KEY-ID", d.Get("api_key").(string))
	verifReq.Header.Add("X-API-KEY-SECRET", d.Get("api_key_secret").(string))

	verifR, err := client.Do(verifReq)
	if err != nil {
		return diag.FromErr(err)
	}
	defer verifR.Body.Close()

	err = json.NewDecoder(verifR.Body).Decode(&verifResp)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, wb := range verifResp {
		if wb["name"].(string) != d.Get("string").(string) {
			continue
		} else if wb["url"].(string) != d.Get("url").(string) {
			continue
		} else if wb["secret"].(string) != d.Get("secret").(string) {
			continue
		} else {
			d.SetId(wb["webhook_id"].(string))
			break
		}
	}

	if d.Get("id").(string) == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Org Webhook",
		})
		return diags
	}

	return diags
}

func resourceOrgWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	resp := make([]map[string]interface{}, 0)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://artifacthub.io/api/v1/webhooks/org/%s", d.Get("org_name").(string)), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("X-API-KEY-ID", d.Get("api_key").(string))
	req.Header.Add("X-API-KEY-SECRET", d.Get("api_key_secret").(string))

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	webhookData := make(map[string]any, 0)

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, wb := range resp {
		if wb["package_id"].(string) == d.Get("id").(string) {
			webhookData = wb
			break
		}
		if wb["name"].(string) != d.Get("string").(string) {
			continue
		} else if d.Get("url").(string) != "" && wb["url"].(string) != d.Get("url").(string) {
			continue
		} else if d.Get("secret").(string) != "" && wb["secret"].(string) != d.Get("secret").(string) {
			continue
		} else {
			webhookData = wb
			break
		}
	}

	if webhookData["webhook_id"].(string) == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to retrieve Org Webhook",
		})
	}

	if err = d.Set("name", webhookData["name"].(string)); err != nil {
		return diag.FromErr(err)
	} else if err = d.Set("description", webhookData["description"].(string)); err != nil {
		return diag.FromErr(err)
	} else if err = d.Set("url", webhookData["url"].(string)); err != nil {
		return diag.FromErr(err)
	} else if err = d.Set("secret", webhookData["secret"].(string)); err != nil {
		return diag.FromErr(err)
	} else if err = d.Set("content_type", webhookData["content_type"].(string)); err != nil {
		return diag.FromErr(err)
	} else if err = d.Set("active", webhookData["active"].(int)); err != nil {
		return diag.FromErr(err)
	} else if err = d.Set("event_kinds", webhookData["event_kinds"].([]int)); err != nil {
		return diag.FromErr(err)
	} else if err = d.Set("packages", webhookData["packages"].([]map[string]string)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(webhookData["webhook_id"].(string))

	return diags
}

func resourceOrgWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var req *http.Request
	var err error

	webhookData := map[string]any{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"url":          d.Get("url").(string),
		"secret":       d.Get("secret").(string),
		"content_type": d.Get("content_type").(string),
		"active":       d.Get("active").(bool),
		"event_kinds":  d.Get("event_kinds").([]int),
		"packages":     d.Get("packages").([]map[string]string),
	}

	tpl := d.Get("template").(string)

	if tpl != "" {
		webhookData["template"] = tpl
	}

	data, _ := json.Marshal(webhookData)

	req, err = http.NewRequest("POST", fmt.Sprintf("https://artifacthub.io/api/v1/webhooks/org/%s/%s", d.Get("org_name").(string), d.Get("id").(string)), bytes.NewBuffer(data))

	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("X-API-KEY-ID", d.Get("api_key").(string))
	req.Header.Add("X-API-KEY-SECRET", d.Get("api_key_secret").(string))

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	resp := make(map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return diag.FromErr(err)
	}

	if r.Status != "204" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Org Webhook",
			Detail:   fmt.Sprintf("status: %s message : %s", r.Status, resp["message"].(string)),
		})
		return diags
	}

	return diags
}

func resourceOrgWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var req *http.Request
	var err error

	req, err = http.NewRequest("DELETE", fmt.Sprintf("https://artifacthub.io/api/v1/webhooks/org/%s/%s", d.Get("org_name").(string), d.Get("id").(string)), nil)

	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("X-API-KEY-ID", d.Get("api_key").(string))
	req.Header.Add("X-API-KEY-SECRET", d.Get("api_key_secret").(string))

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	if r.Status != "204" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "The User Webhook could not be deleted",
		})
		return diags
	}

	d.SetId("")

	return diags
}
