package artifacthub

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccArtifacthubDataSourceHelmPackage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccArtifacthubDataSourceHelmPackageConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "name", "artifact-hub"),
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "repo_name", "artifact-hub"),
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "id", "75ee6e00-b4d5-429e-9d82-33ab730081ff"),
				),
			},
			{
				Config: testAccArtifacthubDataSourceHelmPackageConfigWithVersion(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "name", "artifact-hub"),
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "repo_name", "artifact-hub"),
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "id", "75ee6e00-b4d5-429e-9d82-33ab730081ff"),
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "version", "1.1.1"),
				),
			},
			{
				Config: testAccArtifacthubDataSourceHelmPackageConfigWithVersion(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "name", "artifact-hub"),
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "repo_name", "artifact-hub"),
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "id", "75ee6e00-b4d5-429e-9d82-33ab730081ff"),
					resource.TestCheckResourceAttr("data.artifacthub_helm_package.test", "version", "1.1.1"),
				),
			},
		},
	})
}

func TestAccArtifacthubDataSourceHelmPackageNotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccArtifacthubDataSourceHelmPackageConfigPackageDoesNotExist(),
				ExpectError: regexp.MustCompile(`.*`),
			},
		},
	})
}

func testAccArtifacthubDataSourceHelmPackageConfigBasic() string {
	return `
data "artifacthub_helm_package" "test" {
	repo_name = "artifact-hub"
	name      = "artifact-hub"
}
	`
}

func testAccArtifacthubDataSourceHelmPackageConfigWithVersion() string {
	return `
data "artifacthub_helm_package" "test" {
	repo_name = "artifact-hub"
	name      = "artifact-hub"
	version   = "1.1.1"
}
	`
}

func testAccArtifacthubDataSourceHelmPackageConfigPackageDoesNotExist() string {
	return `
data "artifacthub_helm_package" "test" {
	repo_name = "mypackagedoesnotexist"
	name      = "itverymuchisnotthere"
	version   = "1.1.1"
}
	`
}
