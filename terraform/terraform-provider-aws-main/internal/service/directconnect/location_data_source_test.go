package directconnect_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccDirectConnectLocationDataSource_basic(t *testing.T) {
	dsResourceName := "data.aws_dx_location.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, directconnect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLocationDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsResourceName, "available_macsec_port_speeds.#"),
					resource.TestCheckResourceAttrSet(dsResourceName, "available_port_speeds.#"),
					resource.TestCheckResourceAttrSet(dsResourceName, "available_providers.#"),
					resource.TestCheckResourceAttrSet(dsResourceName, "location_code"),
					resource.TestCheckResourceAttrSet(dsResourceName, "location_name"),
				),
			},
		},
	})
}

const testAccLocationDataSourceConfig_basic = `
data "aws_dx_locations" "test" {}

locals {
  location_codes = tolist(data.aws_dx_locations.test.location_codes)
}

data "aws_dx_location" "test" {
  location_code = local.location_codes[length(local.location_codes) - 1]
}
`
