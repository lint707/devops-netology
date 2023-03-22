package lightsail_test

import (
	"context"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
)

// Lightsail Domains can only be created in specific regions.

// testAccLightsailDomainRegion is the chosen Lightsail Domains testing region
//
// Cached to prevent issues should multiple regions become available.
var testAccLightsailDomainRegion string

// testAccProviderLightsailDomain is the Lightsail Domains provider instance
//
// This Provider can be used in testing code for API calls without requiring
// the use of saving and referencing specific ProviderFactories instances.
//
// testAccPreCheckDomain(t) must be called before using this provider instance.
var testAccProviderLightsailDomain *schema.Provider

// testAccProviderLightsailDomainConfigure ensures the provider is only configured once
var testAccProviderLightsailDomainConfigure sync.Once

// testAccPreCheckDomain verifies AWS credentials and that Lightsail Domains is supported
func testAccPreCheckDomain(t *testing.T) {
	acctest.PreCheckPartitionHasService(lightsail.EndpointsID, t)

	region := testAccGetDomainRegion()

	if region == "" {
		t.Skip("Lightsail Domains not available in this AWS Partition")
	}

	// Since we are outside the scope of the Terraform configuration we must
	// call Configure() to properly initialize the provider configuration.
	testAccProviderLightsailDomainConfigure.Do(func() {
		ctx := context.Background()
		var err error
		testAccProviderLightsailDomain, err = provider.New(ctx)

		if err != nil {
			t.Fatal(err)
		}

		config := map[string]interface{}{
			"region": region,
		}

		diags := testAccProviderLightsailDomain.Configure(ctx, terraform.NewResourceConfigRaw(config))

		if diags != nil && diags.HasError() {
			for _, d := range diags {
				if d.Severity == diag.Error {
					t.Fatalf("error configuring Lightsail Domains provider: %s", d.Summary)
				}
			}
		}
	})
}

// testAccDomainRegionProviderConfig is the Terraform provider configuration for Lightsail Domains region testing
//
// Testing Lightsail Domains assumes no other provider configurations
// are necessary and overwrites the "aws" provider configuration.
func testAccDomainRegionProviderConfig() string {
	return acctest.ConfigRegionalProvider(testAccGetDomainRegion())
}

// testAccGetDomainRegion returns the Lightsail Domains region for testing
func testAccGetDomainRegion() string {
	if testAccLightsailDomainRegion != "" {
		return testAccLightsailDomainRegion
	}

	// AWS Commercial: https://lightsail.aws.amazon.com/ls/docs/en_us/articles/lightsail-how-to-create-dns-entry
	// AWS GovCloud (US) - service not supported: https://docs.aws.amazon.com/govcloud-us/latest/UserGuide/using-services.html
	// AWS China - service not supported: https://www.amazonaws.cn/en/about-aws/regional-product-services/
	switch acctest.Partition() {
	case endpoints.AwsPartitionID:
		testAccLightsailDomainRegion = endpoints.UsEast1RegionID
	}

	return testAccLightsailDomainRegion
}
