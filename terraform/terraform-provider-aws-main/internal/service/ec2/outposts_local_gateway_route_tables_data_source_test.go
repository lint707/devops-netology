package ec2_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccEC2OutpostsLocalGatewayRouteTablesDataSource_basic(t *testing.T) {
	dataSourceName := "data.aws_ec2_local_gateway_route_tables.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOutpostsLocalGatewayRouteTablesDataSourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrGreaterThanValue(dataSourceName, "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccEC2OutpostsLocalGatewayRouteTablesDataSource_filter(t *testing.T) {
	dataSourceName := "data.aws_ec2_local_gateway_route_tables.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOutpostsLocalGatewayRouteTablesDataSourceConfig_filter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "1"),
				),
			},
		},
	})
}

func testAccOutpostsLocalGatewayRouteTablesDataSourceConfig_basic() string {
	return `
data "aws_ec2_local_gateway_route_tables" "test" {}
`
}

func testAccOutpostsLocalGatewayRouteTablesDataSourceConfig_filter() string {
	return `
data "aws_ec2_local_gateway_route_tables" "all" {}

data "aws_ec2_local_gateway_route_tables" "test" {
  filter {
    name   = "local-gateway-route-table-id"
    values = [tolist(data.aws_ec2_local_gateway_route_tables.all.ids)[0]]
  }
}
`
}
