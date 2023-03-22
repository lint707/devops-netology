package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func testAccTransitGatewayVPNAttachmentDataSource_idAndVPNConnectionID(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	dataSourceName := "data.aws_ec2_transit_gateway_vpn_attachment.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	vpnConnectionResourceName := "aws_vpn_connection.test"

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
				Config: testAccTransitGatewayVPNAttachmentDataSourceConfig_idAndVPNConnectionID(rName, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "transit_gateway_id", transitGatewayResourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vpn_connection_id", vpnConnectionResourceName, "id"),
				),
			},
		},
	})
}

func testAccTransitGatewayVPNAttachmentDataSource_filter(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	dataSourceName := "data.aws_ec2_transit_gateway_vpn_attachment.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	vpnConnectionResourceName := "aws_vpn_connection.test"

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
				Config: testAccTransitGatewayVPNAttachmentDataSourceConfig_filter(rName, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "transit_gateway_id", transitGatewayResourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vpn_connection_id", vpnConnectionResourceName, "id"),
				),
			},
		},
	})
}

func testAccTransitGatewayVPNAttachmentDataSourceConfig_base(rName string, rBgpAsn int) string {
	return fmt.Sprintf(`
resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_customer_gateway" "test" {
  bgp_asn    = %[2]d
  ip_address = "178.0.0.1"
  type       = "ipsec.1"

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpn_connection" "test" {
  customer_gateway_id = aws_customer_gateway.test.id
  transit_gateway_id  = aws_ec2_transit_gateway.test.id
  type                = aws_customer_gateway.test.type

  tags = {
    Name = %[1]q
  }
}
`, rName, rBgpAsn)
}

func testAccTransitGatewayVPNAttachmentDataSourceConfig_idAndVPNConnectionID(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccTransitGatewayVPNAttachmentDataSourceConfig_base(rName, rBgpAsn), `
data "aws_ec2_transit_gateway_vpn_attachment" "test" {
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpn_connection_id  = aws_vpn_connection.test.id
}
`)
}

func testAccTransitGatewayVPNAttachmentDataSourceConfig_filter(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccTransitGatewayVPNAttachmentDataSourceConfig_base(rName, rBgpAsn), `
data "aws_ec2_transit_gateway_vpn_attachment" "test" {
  filter {
    name   = "resource-id"
    values = [aws_vpn_connection.test.id]
  }
}
`)
}
