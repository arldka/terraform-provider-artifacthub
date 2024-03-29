package artifacthub

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"artifacthub": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestProviderConfigure(t *testing.T) {
	ctx := context.TODO()

	apiKey := os.Getenv("ARTIFACTHUB_API_KEY")
	apiKeySecret := os.Getenv("ARTIFACTHUB_API_KEY_SECRET")

	os.Setenv("ARTIFACTHUB_API_KEY", uuid.New().String())
	os.Setenv("ARTIFACTHUB_API_KEY_SECRET", uuid.New().String())

	p := Provider()
	rc := terraform.NewResourceConfigRaw(map[string]interface{}{})

	diags := p.Configure(ctx, rc)

	os.Unsetenv("ARTIFACTHUB_API_KEY")
	os.Unsetenv("ARTIFACTHUB_API_KEY_SECRET")

	if diags.HasError() {
		t.Fatal(diags)
	}

	os.Setenv("ARTIFACTHUB_API_KEY_SECRET", apiKeySecret)
	os.Setenv("ARTIFACTHUB_API_KEY", apiKey)
}

func TestProviderConfigureFail(t *testing.T) {
	ctx := context.TODO()

	apiKey := os.Getenv("ARTIFACTHUB_API_KEY")
	apiKeySecret := os.Getenv("ARTIFACTHUB_API_KEY_SECRET")

	os.Unsetenv("ARTIFACTHUB_API_KEY_SECRET")

	p := Provider()
	rc := terraform.NewResourceConfigRaw(map[string]interface{}{})

	diags := p.Configure(ctx, rc)

	if !diags.HasError() {
		t.Fatal(diags)
	}

	os.Unsetenv("ARTIFACTHUB_API_KEY")
	os.Setenv("ARTIFACTHUB_API_KEY_SECRET", apiKeySecret)

	diags = p.Configure(ctx, rc)

	os.Setenv("ARTIFACTHUB_API_KEY", apiKey)

	if !diags.HasError() {
		t.Fatal(diags)
	}

}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ARTIFACTHUB_API_KEY"); v == "" {
		t.Fatal("ARTIFACTHUB_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ARTIFACTHUB_API_KEY_SECRET"); v == "" {
		t.Fatal("ARTIFACTHUB_API_KEY_SECRET must be set for acceptance tests")
	}
}
