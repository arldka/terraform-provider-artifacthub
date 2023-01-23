package artifacthub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "application/json",
			},
			"template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"event_kinds": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     schema.TypeInt,
			},
			"packages": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"package_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"webhook_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

	req, err = http.NewRequest("POST", fmt.Sprintf("https://artifacthub.io/api/v1/webhooks/user"), bytes.NewBuffer(data))

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
			Summary:  "Unable to create User Webhook",
			Detail:   resp["message"].(string),
		})
		return diags
	}

	d.Set("package_id", resp[""])

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceUserWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var req *http.Request
	var err error

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceUserWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var req *http.Request
	var err error

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceUserWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var req *http.Request
	var err error

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
