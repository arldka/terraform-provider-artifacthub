package artifacthub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserWebhookCreate,
		ReadContext:   resourceUserWebhookRead,
		UpdateContext: resourceUserWebhookUpdate,
		DeleteContext: resourceUserWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the webhook. ",
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
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"packages": {
				Type:        schema.TypeList,
				Description: "Packages list",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserWebhookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var req *http.Request
	var err error

	webhookData := map[string]any{
		"name": d.Get("name").(string),
		//	"description":  d.Get("description").(string),
		"url": d.Get("url").(string),
		//	"secret":       d.Get("secret").(string),
		//	"content_type": d.Get("content_type").(string),
		//	"template":     d.Get("template").(string),
		"active": d.Get("active").(bool),
	}

	// Defines the webhook event kinds
	rawEventKinds := d.Get("event_kinds").([]interface{})
	var eventKinds []int
	for _, kind := range rawEventKinds {
		k := kind.(int)
		eventKinds = append(eventKinds, k)
	}
	webhookData["event_kinds"] = eventKinds

	// Defines the webhook packages
	rawPackages := d.Get("packages").([]interface{})
	var packages []map[string]string
	for _, pkg := range rawPackages {
		p := map[string]string{
			"package_id": pkg.(string),
		}
		packages = append(packages, p)
	}
	webhookData["packages"] = packages

	// Defines the webhook description
	if d.Get("description").(string) != "" {
		webhookData["description"] = d.Get("description").(string)
	}

	// Defines the webhook secret
	if d.Get("secret").(string) != "" {
		webhookData["secret"] = d.Get("secret").(string)
	}

	// Defines the webhook content type
	if d.Get("content_type").(string) != "" {
		webhookData["content_type"] = d.Get("content_type").(string)
	}

	// Defines the webhook template
	if d.Get("template").(string) != "" {
		webhookData["template"] = d.Get("template").(string)
	}

	data, _ := json.Marshal(webhookData)

	req, err = http.NewRequest("POST", "https://artifacthub.io/api/v1/webhooks/user", bytes.NewBuffer(data))

	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("X-API-KEY-ID", m.(*Config).ApiKey)
	req.Header.Add("X-API-KEY-SECRET", m.(*Config).ApiKeySecret)

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	defer r.Body.Close()

	resp := make(map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&resp)

	if r.StatusCode != 201 && err != nil {
		return diag.FromErr(err)
	}

	if r.StatusCode != 201 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create User Webhook",
			Detail:   resp["message"].(string),
		})
		return diags
	}

	verifResp := make([]map[string]interface{}, 0)

	verifReq, err := http.NewRequest("GET", "https://artifacthub.io/api/v1/webhooks/user", nil)

	if err != nil {
		return diag.FromErr(err)
	}

	verifReq.Header.Add("X-API-KEY-ID", m.(*Config).ApiKey)
	verifReq.Header.Add("X-API-KEY-SECRET", m.(*Config).ApiKeySecret)

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
		if wb["name"].(string) != d.Get("name").(string) {
			continue
		} else {
			d.SetId(wb["webhook_id"].(string))
			break
		}
	}

	if d.Id() == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create User Webhook",
		})
		return diags
	}

	return diags
}

func resourceUserWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	resp := make([]map[string]interface{}, 0)

	req, err := http.NewRequest("GET", "https://artifacthub.io/api/v1/webhooks/user", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("X-API-KEY-ID", m.(*Config).ApiKey)
	req.Header.Add("X-API-KEY-SECRET", m.(*Config).ApiKeySecret)

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
		if wb["webhook_id"].(string) == d.Id() {
			webhookData = wb
			break
		}
		if wb["name"].(string) != d.Get("string").(string) {
			continue
		} else {
			webhookData = wb
			break
		}
	}

	if webhookData["webhook_id"].(string) == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to retrieve User Webhook",
		})
		return diags
	}

	// Check if each key of webookData is non null before setting it

	if err = d.Set("name", webhookData["name"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("description", webhookData["description"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("url", webhookData["url"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("active", webhookData["active"].(bool)); err != nil {
		return diag.FromErr(err)
	}

	// Defines the webhook event kinds
	rawEventKinds := webhookData["event_kinds"].([]interface{})
	var eventKinds []float64
	for _, kind := range rawEventKinds {
		k := kind.(float64)
		eventKinds = append(eventKinds, k)
	}
	if err = d.Set("event_kinds", eventKinds); err != nil {
		return diag.FromErr(err)
	}

	// Defines the webhook packages
	rawPackages := webhookData["packages"].([]interface{})
	var packages []string
	for _, pkg := range rawPackages {
		pm := pkg.(map[string]interface{})
		p := pm["package_id"].(string)
		packages = append(packages, p)
	}
	if err = d.Set("packages", packages); err != nil {
		return diag.FromErr(err)
	}

	if webhookData["template"] != nil {
		if err = d.Set("template", webhookData["template"].(string)); err != nil {
			return diag.FromErr(err)
		}
	}
	if webhookData["secret"] != nil {
		if err = d.Set("secret", webhookData["secret"].(string)); err != nil {
			return diag.FromErr(err)
		}
	}
	if webhookData["content_type"] != nil {
		if err = d.Set("content_type", webhookData["content_type"].(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(webhookData["webhook_id"].(string))

	return diags
}

func resourceUserWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	fmt.Fprintln(os.Stdout, "Start of Update")

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

	req, err = http.NewRequest("POST", fmt.Sprintf("https://artifacthub.io/api/v1/webhooks/user/"+d.Get("id").(string)), bytes.NewBuffer(data))

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
			Summary:  "Unable to update User Webhook",
			Detail:   fmt.Sprintf("status: %s message : %s", r.Status, resp["message"].(string)),
		})
		return diags
	}

	return diags
}

func resourceUserWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var req *http.Request
	var err error

	req, err = http.NewRequest("DELETE", fmt.Sprintf("https://artifacthub.io/api/v1/webhooks/user/"+d.Id()), nil)

	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("X-API-KEY-ID", m.(*Config).ApiKey)
	req.Header.Add("X-API-KEY-SECRET", m.(*Config).ApiKeySecret)

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	if r.StatusCode != 204 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "The User Webhook could not be deleted",
		})
		return diags
	}

	d.SetId("")

	return diags
}
