package elasticsearch_test

import (
	"fmt"

	awspolicy "github.com/hashicorp/awspolicyequivalence"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckPolicyMatch(resource, attr, expectedPolicy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		given, ok := rs.Primary.Attributes[attr]
		if !ok {
			return fmt.Errorf("Attribute %q not found for %q", attr, resource)
		}

		areEquivalent, err := awspolicy.PoliciesAreEquivalent(given, expectedPolicy)
		if err != nil {
			return fmt.Errorf("Comparing AWS Policies failed: %s", err)
		}

		if !areEquivalent {
			return fmt.Errorf("AWS policies differ.\nGiven: %s\nExpected: %s", given, expectedPolicy)
		}

		return nil
	}
}
