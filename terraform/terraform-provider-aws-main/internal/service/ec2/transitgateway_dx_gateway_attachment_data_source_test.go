package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func testAccTransitGatewayDxGatewayAttachmentDataSource_TransitGatewayIdAndDxGatewayID(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	dataSourceName := "data.aws_ec2_transit_gateway_dx_gateway_attachment.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	dxGatewayResourceName := "aws_dx_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			testAccPreCheckTransitGateway(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayDxGatewayAttachmentDataSourceConfig_transit(rName, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "dx_gateway_id", dxGatewayResourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "transit_gateway_id", transitGatewayResourceName, "id"),
				),
			},
		},
	})
}

func testAccTransitGatewayDxGatewayAttachmentDataSource_filter(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	dataSourceName := "data.aws_ec2_transit_gateway_dx_gateway_attachment.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	dxGatewayResourceName := "aws_dx_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			testAccPreCheckTransitGateway(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayDxGatewayAttachmentDataSourceConfig_transitFilter(rName, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "dx_gateway_id", dxGatewayResourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "transit_gateway_id", transitGatewayResourceName, "id"),
				),
			},
		},
	})
}

func testAccTransitGatewayDxGatewayAttachmentDataSourceConfig_transit(rName string, rBgpAsn int) string {
	return fmt.Sprintf(`
resource "aws_dx_gateway" "test" {
  name            = %[1]q
  amazon_side_asn = "%[2]d"
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_dx_gateway_association" "test" {
  dx_gateway_id         = aws_dx_gateway.test.id
  associated_gateway_id = aws_ec2_transit_gateway.test.id

  allowed_prefixes = [
    "10.255.255.0/30",
    "10.255.255.8/30",
  ]
}

data "aws_ec2_transit_gateway_dx_gateway_attachment" "test" {
  transit_gateway_id = aws_dx_gateway_association.test.associated_gateway_id
  dx_gateway_id      = aws_dx_gateway_association.test.dx_gateway_id
}
`, rName, rBgpAsn)
}

func testAccTransitGatewayDxGatewayAttachmentDataSourceConfig_transitFilter(rName string, rBgpAsn int) string {
	return fmt.Sprintf(`
resource "aws_dx_gateway" "test" {
  name            = %[1]q
  amazon_side_asn = "%[2]d"
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_dx_gateway_association" "test" {
  dx_gateway_id         = aws_dx_gateway.test.id
  associated_gateway_id = aws_ec2_transit_gateway.test.id

  allowed_prefixes = [
    "10.255.255.0/30",
    "10.255.255.8/30",
  ]
}

data "aws_ec2_transit_gateway_dx_gateway_attachment" "test" {
  filter {
    name   = "resource-id"
    values = [aws_dx_gateway_association.test.dx_gateway_id]
  }
}
`, rName, rBgpAsn)
}
