package artifacthub

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Write a test for the resourceUserWebhookCreate function
// Path: artifacthub/resource_user_webhook.go

func TestAccArtifacthubResourceUserWebhookCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccArtifacthubResourceUserWebhookCreateSinglePackage(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("artifacthub_user_webhook.test", "name", "testSinglePackage"),
					resource.TestCheckResourceAttr("artifacthub_user_webhook.test", "description", "test"),
					resource.TestCheckResourceAttr("artifacthub_user_webhook.test", "url", "https://test.com"),
					resource.TestCheckResourceAttr("artifacthub_user_webhook.test", "packages.#", "1"),
				),
			},
		},
	})
}

func testAccArtifacthubResourceUserWebhookCreateSinglePackage() string {
	return `
	resource "artifacthub_user_webhook" "test_single_package" {
		name = "testSinglePackage"
		description = "test"
		url = "https://test.com"
		packages {
			package_id = "75ee6e00-b4d5-429e-9d82-33ab730081ff"
		}
		event_kinds = [0]
		active = false
	}
	`
}
