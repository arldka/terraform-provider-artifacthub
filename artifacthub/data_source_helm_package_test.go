package artifacthub

import (
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
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

func TestDataSourceHelmPackageRead(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name: "basic",
			input: map[string]interface{}{
				"repo_name": "artifact-hub",
				"name":      "artifact-hub",
			},
			expected: map[string]interface{}{
				"version": "1.1.1",
			},
			wantErr: false,
		},
		{
			name: "with_version",
			input: map[string]interface{}{
				"repo_name": "artifact-hub",
				"name":      "artifact-hub",
				"version":   "1.1.1",
			},
			expected: map[string]interface{}{
				"version": "1.1.1",
			},
			wantErr: false,
		},
		{
			name: "not_found",
			input: map[string]interface{}{
				"repo_name": "artifact-hub",
				"name":      "artifact-hub-not-found",
			},
			expected: map[string]interface{}{},
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := dataSourceHelmPackage()
			ctx := context.Background()
			rd := schema.TestResourceDataRaw(t, d.Schema, tc.input)
			m := &Config{os.Getenv("ARTIFACTHUB_API_KEY"), os.Getenv("ARTIFACTHUB_API_KEY_SECRET")}
			err := dataSourceHelmPackageRead(ctx, rd, m)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			for k, v := range tc.expected {
				assert.Equal(t, v, rd.Get(k))
			}
		})
	}
}
