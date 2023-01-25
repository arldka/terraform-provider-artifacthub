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
}

func TestProviderConfigureFail(t *testing.T) {
	ctx := context.TODO()

	apiKey := os.Getenv("ARTIFACTHUB_API_KEY")
	apiKeySecret := os.Getenv("ARTIFACTHUB_API_KEY_SECRET")

	os.Unsetenv("ARTIFACTHUB_API_KEY")
	os.Unsetenv("ARTIFACTHUB_API_KEY_SECRET")

	p := Provider()
	rc := terraform.NewResourceConfigRaw(map[string]interface{}{})

	diags := p.Configure(ctx, rc)

	os.Setenv("ARTIFACTHUB_API_KEY", apiKey)
	os.Setenv("ARTIFACTHUB_API_KEY_SECRET", apiKeySecret)

	if !diags.HasError() {
		t.Fatal(diags)
	}

}
