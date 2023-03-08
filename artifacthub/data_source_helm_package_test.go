package artifacthub

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccArtifacthubDataSourceHelmPackage_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccArtifacthubDataSourceHelmPackageConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "name", "artifact-hub"),
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "repo_name", "artifact-hub"),
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "id", "75ee6e00-b4d5-429e-9d82-33ab730081ff"),
				),
			},
		},
	})

}

func testAccArtifacthubDataSourceHelmPackageConfig_basic() string {
	return `
data "artifacthub_helm_package" "test" {
	repo_name = "artifact-hub"
	name      = "artifact-hub"
}
	`
}
