package networkmanager_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/networkmanager"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfnetworkmanager "github.com/hashicorp/terraform-provider-aws/internal/service/networkmanager"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccNetworkManagerTransitGatewayPeering_basic(t *testing.T) {
	var v networkmanager.TransitGatewayPeering
	resourceName := "aws_networkmanager_transit_gateway_peering.test"
	tgwResourceName := "aws_ec2_transit_gateway.test"
	testExternalProviders := map[string]resource.ExternalProvider{
		"awscc": {
			Source:            "hashicorp/awscc",
			VersionConstraint: "0.29.0",
		},
	}
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, networkmanager.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		ExternalProviders:        testExternalProviders,
		CheckDestroy:             testAccCheckTransitGatewayPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayPeeringConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTransitGatewayPeeringExists(resourceName, &v),
					acctest.MatchResourceAttrGlobalARN(resourceName, "arn", "networkmanager", regexp.MustCompile(`peering/.+`)),
					resource.TestCheckResourceAttrSet(resourceName, "core_network_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "core_network_id"),
					resource.TestCheckResourceAttr(resourceName, "edge_location", acctest.Region()),
					acctest.CheckResourceAttrAccountID(resourceName, "owner_account_id"),
					resource.TestCheckResourceAttrPair(resourceName, "resource_arn", tgwResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_arn", tgwResourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "transit_gateway_peering_attachment_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkManagerTransitGatewayPeering_disappears(t *testing.T) {
	var v networkmanager.TransitGatewayPeering
	resourceName := "aws_networkmanager_transit_gateway_peering.test"
	testExternalProviders := map[string]resource.ExternalProvider{
		"awscc": {
			Source:            "hashicorp/awscc",
			VersionConstraint: "0.29.0",
		},
	}
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, networkmanager.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		ExternalProviders:        testExternalProviders,
		CheckDestroy:             testAccCheckTransitGatewayPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayPeeringConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayPeeringExists(resourceName, &v),
					acctest.CheckResourceDisappears(acctest.Provider, tfnetworkmanager.ResourceTransitGatewayPeering(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccNetworkManagerTransitGatewayPeering_tags(t *testing.T) {
	var v networkmanager.TransitGatewayPeering
	resourceName := "aws_networkmanager_transit_gateway_peering.test"
	testExternalProviders := map[string]resource.ExternalProvider{
		"awscc": {
			Source:            "hashicorp/awscc",
			VersionConstraint: "0.29.0",
		},
	}
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, networkmanager.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		ExternalProviders:        testExternalProviders,
		CheckDestroy:             testAccCheckTransitGatewayPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayPeeringConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayPeeringExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTransitGatewayPeeringConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayPeeringExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccTransitGatewayPeeringConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayPeeringExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckTransitGatewayPeeringExists(n string, v *networkmanager.TransitGatewayPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Network Manager Transit Gateway Peering ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).NetworkManagerConn

		output, err := tfnetworkmanager.FindTransitGatewayPeeringByID(context.Background(), conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckTransitGatewayPeeringDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).NetworkManagerConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_networkmanager_transit_gateway_peering" {
			continue
		}

		_, err := tfnetworkmanager.FindTransitGatewayPeeringByID(context.Background(), conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Network Manager Transit Gateway Peering %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccTransitGatewayPeeringConfig_base(rName string) string {
	return fmt.Sprintf(`
data "aws_region" "current" {}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_policy_table" "test" {
  transit_gateway_id = aws_ec2_transit_gateway.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_networkmanager_global_network" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "awscc_networkmanager_core_network" "test" {
  global_network_id = aws_networkmanager_global_network.test.id
  policy_document   = jsonencode(jsondecode(data.aws_networkmanager_core_network_policy_document.test.json))
}

data "aws_networkmanager_core_network_policy_document" "test" {
  core_network_configuration {
    # Don't overlap with default TGW ASN: 64512.
    asn_ranges = ["65022-65534"]

    edge_locations {
      location = data.aws_region.current.name
    }
  }

  segments {
    name = "test"
  }
}
`, rName)
}

func testAccTransitGatewayPeeringConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccTransitGatewayPeeringConfig_base(rName), `
resource "aws_networkmanager_transit_gateway_peering" "test" {
  core_network_id     = awscc_networkmanager_core_network.test.id
  transit_gateway_arn = aws_ec2_transit_gateway.test.arn

  depends_on = [aws_ec2_transit_gateway_policy_table.test]
}
`)
}

func testAccTransitGatewayPeeringConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccTransitGatewayPeeringConfig_base(rName), fmt.Sprintf(`
resource "aws_networkmanager_transit_gateway_peering" "test" {
  core_network_id     = awscc_networkmanager_core_network.test.id
  transit_gateway_arn = aws_ec2_transit_gateway.test.arn

  tags = {
    %[1]q = %[2]q
  }

  depends_on = [aws_ec2_transit_gateway_policy_table.test]
}
`, tagKey1, tagValue1))
}

func testAccTransitGatewayPeeringConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccTransitGatewayPeeringConfig_base(rName), fmt.Sprintf(`
resource "aws_networkmanager_transit_gateway_peering" "test" {
  core_network_id     = awscc_networkmanager_core_network.test.id
  transit_gateway_arn = aws_ec2_transit_gateway.test.arn

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }

  depends_on = [aws_ec2_transit_gateway_policy_table.test]
}
`, tagKey1, tagValue1, tagKey2, tagValue2))
}
