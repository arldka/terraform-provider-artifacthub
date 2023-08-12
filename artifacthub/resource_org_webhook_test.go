package artifacthub

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// Write a test for the resourceOrgWebhookCreate function
// Path: artifacthub/resource_org_webhook.go

func TestAccArtifacthubResourceOrgWebhookCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccArtifacthubResourceOrgWebhookCreateSinglePackage(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("artifacthub_org_webhook.test_single_package", "name", "testSinglePackage"),
					resource.TestCheckResourceAttr("artifacthub_org_webhook.test_single_package", "org_name", "terraform-provider-artifacthub"),
					resource.TestCheckResourceAttr("artifacthub_org_webhook.test_single_package", "description", "test"),
					resource.TestCheckResourceAttr("artifacthub_org_webhook.test_single_package", "url", "https://test.com"),
				),
			},
		},
	})
}

func testAccArtifacthubResourceOrgWebhookCreateSinglePackage() string {
	return `
	resource "artifacthub_org_webhook" "test_single_package" {
		name = "testSinglePackage"
		org_name = "terraform-provider-artifacthub"
		description = "test"
		url = "https://test.com"
		packages = ["75ee6e00-b4d5-429e-9d82-33ab730081ff"]
		event_kinds = [0]
		active = false
	}
	`
}

func TestResourceOrgWebhookCreateReadDelete(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
		wantErr  bool
		action   string
	}{
		{
			name: "create_single_webhook",
			input: map[string]interface{}{
				"name":        "testSinglePackageUnit",
				"org_name":    "terraform-provider-artifacthub",
				"description": "yes",
				"url":         "https://google.com",
				"packages":    []interface{}{"75ee6e00-b4d5-429e-9d82-33ab730081ff"},
				"event_kinds": []interface{}{0},
				"active":      "false",
			},
			expected: map[string]interface{}{
				"name":        "testSinglePackageUnit",
				"org_name":    "terraform-provider-artifacthub",
				"description": "yes",
				"url":         "https://google.com",
				"packages":    []interface{}{"75ee6e00-b4d5-429e-9d82-33ab730081ff"},
				"event_kinds": []interface{}{0},
				"active":      false,
			},
			wantErr: false,
			action:  "create",
		},
		{
			name: "update_single_webhook",
			input: map[string]interface{}{
				"name":        "testSinglePackageUnit",
				"org_name":    "terraform-provider-artifacthub",
				"description": "no",
				"url":         "https://test.com",
				"packages":    []interface{}{"75ee6e00-b4d5-429e-9d82-33ab730081ff"},
				"event_kinds": []interface{}{0},
				"active":      "true",
			},
			expected: map[string]interface{}{
				"name":        "testSinglePackageUnit",
				"org_name":    "terraform-provider-artifacthub",
				"description": "no",
				"url":         "https://test.com",
				"packages":    []interface{}{"75ee6e00-b4d5-429e-9d82-33ab730081ff"},
				"event_kinds": []interface{}{0},
				"active":      true,
			},
			wantErr: false,
			action:  "update",
		},
		{
			name: "delete_single_webhook",
			input: map[string]interface{}{
				"name":     "testSinglePackageUnit",
				"org_name": "terraform-provider-artifacthub",
			},
			expected: map[string]interface{}{},
			wantErr:  false,
			action:   "delete",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := resourceOrgWebhook()
			ctx := context.Background()
			rr := schema.TestResourceDataRaw(t, r.Schema, tc.input)
			m := &Config{os.Getenv("ARTIFACTHUB_API_KEY"), os.Getenv("ARTIFACTHUB_API_KEY_SECRET")}
			if tc.action == "create" {
				err := resourceOrgWebhookCreate(ctx, rr, m)
				if tc.wantErr {
					fmt.Fprintln(os.Stdout, err)
					assert.NotEqual(t, err, diag.Diagnostics(diag.Diagnostics(nil)))
				} else {
					assert.Equal(t, err, diag.Diagnostics(diag.Diagnostics(nil)))
				}
			} else if tc.action == "update" {
				// Read before update to set the ID
				err := resourceOrgWebhookRead(ctx, rr, m)
				if tc.wantErr {
					fmt.Fprintln(os.Stdout, err)
					assert.NotEqual(t, err, diag.Diagnostics(diag.Diagnostics(nil)))
				} else {
					assert.Equal(t, err, diag.Diagnostics(diag.Diagnostics(nil)))
				}

				// Preserve the ID for the update
				id := rr.Id()

				// Reset the resource data to the updated values as it has been modified by the read
				rr = schema.TestResourceDataRaw(t, r.Schema, tc.input)

				// Set the ID to the value retrieved by the read
				rr.SetId(id)

				err = resourceOrgWebhookUpdate(ctx, rr, m)
				if tc.wantErr {
					fmt.Fprintln(os.Stdout, err)
					assert.NotEqual(t, err, diag.Diagnostics(diag.Diagnostics(nil)))
				} else {
					assert.Equal(t, err, diag.Diagnostics(diag.Diagnostics(nil)))
				}
			} else if tc.action == "delete" {
				// Read before update to set the ID
				err := resourceOrgWebhookRead(ctx, rr, m)
				if tc.wantErr {
					fmt.Fprintln(os.Stdout, err)
					assert.NotEqual(t, err, diag.Diagnostics(diag.Diagnostics(nil)))
				} else {
					assert.Equal(t, err, diag.Diagnostics(diag.Diagnostics(nil)))
				}
				err = resourceOrgWebhookDelete(ctx, rr, m)
				if tc.wantErr {
					fmt.Fprintln(os.Stdout, err)
					assert.NotEqual(t, err, diag.Diagnostics(diag.Diagnostics(nil)))
				} else {
					assert.Equal(t, err, diag.Diagnostics(diag.Diagnostics(nil)))
				}
			}
			for k, v := range tc.expected {
				assert.Equal(t, v, rr.Get(k))
			}
		})
	}
}
