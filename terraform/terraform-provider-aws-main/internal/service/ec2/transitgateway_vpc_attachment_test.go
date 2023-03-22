package ec2_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func testAccTransitGatewayVPCAttachment_basic(t *testing.T) {
	var transitGatewayVpcAttachment1 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	vpcResourceName := "aws_vpc.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment1),
					resource.TestCheckResourceAttr(resourceName, "dns_support", ec2.DnsSupportValueEnable),
					resource.TestCheckResourceAttr(resourceName, "ipv6_support", ec2.Ipv6SupportValueDisable),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_association", "true"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_propagation", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_id", transitGatewayResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", vpcResourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_owner_id"),
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

func testAccTransitGatewayVPCAttachment_disappears(t *testing.T) {
	var transitGatewayVpcAttachment1 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment1),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceTransitGatewayVPCAttachment(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccTransitGatewayVPCAttachment_ApplianceModeSupport(t *testing.T) {
	var transitGatewayVpcAttachment1, transitGatewayVpcAttachment2 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_applianceModeSupport(rName, ec2.ApplianceModeSupportValueDisable),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment1),
					resource.TestCheckResourceAttr(resourceName, "appliance_mode_support", ec2.ApplianceModeSupportValueDisable),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_applianceModeSupport(rName, ec2.ApplianceModeSupportValueEnable),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment2),
					testAccCheckTransitGatewayVPCAttachmentNotRecreated(&transitGatewayVpcAttachment1, &transitGatewayVpcAttachment2),
					resource.TestCheckResourceAttr(resourceName, "appliance_mode_support", ec2.ApplianceModeSupportValueEnable),
				),
			},
		},
	})
}

func testAccTransitGatewayVPCAttachment_DNSSupport(t *testing.T) {
	var transitGatewayVpcAttachment1, transitGatewayVpcAttachment2 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_dnsSupport(rName, ec2.DnsSupportValueDisable),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment1),
					resource.TestCheckResourceAttr(resourceName, "dns_support", ec2.DnsSupportValueDisable),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_dnsSupport(rName, ec2.DnsSupportValueEnable),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment2),
					testAccCheckTransitGatewayVPCAttachmentNotRecreated(&transitGatewayVpcAttachment1, &transitGatewayVpcAttachment2),
					resource.TestCheckResourceAttr(resourceName, "dns_support", ec2.DnsSupportValueEnable),
				),
			},
		},
	})
}

func testAccTransitGatewayVPCAttachment_IPv6Support(t *testing.T) {
	var transitGatewayVpcAttachment1, transitGatewayVpcAttachment2 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_ipv6Support(rName, ec2.Ipv6SupportValueEnable),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment1),
					resource.TestCheckResourceAttr(resourceName, "ipv6_support", ec2.Ipv6SupportValueEnable),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_ipv6Support(rName, ec2.Ipv6SupportValueDisable),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment2),
					testAccCheckTransitGatewayVPCAttachmentNotRecreated(&transitGatewayVpcAttachment1, &transitGatewayVpcAttachment2),
					resource.TestCheckResourceAttr(resourceName, "ipv6_support", ec2.Ipv6SupportValueDisable),
				),
			},
		},
	})
}

func testAccTransitGatewayVPCAttachment_SharedTransitGateway(t *testing.T) {
	var transitGatewayVpcAttachment1 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckAlternateAccount(t)
			testAccPreCheckTransitGateway(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(t),
		CheckDestroy:             testAccCheckTransitGatewayVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_sharedTransitGateway(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment1),
				),
			},
			{
				Config:            testAccTransitGatewayVPCAttachmentConfig_sharedTransitGateway(rName),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTransitGatewayVPCAttachment_SubnetIDs(t *testing.T) {
	var transitGatewayVpcAttachment1, transitGatewayVpcAttachment2, transitGatewayVpcAttachment3 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_subnetIds2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment1),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_subnetIds1(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment2),
					testAccCheckTransitGatewayVPCAttachmentNotRecreated(&transitGatewayVpcAttachment1, &transitGatewayVpcAttachment2),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "1"),
				),
			},
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_subnetIds2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment3),
					testAccCheckTransitGatewayVPCAttachmentNotRecreated(&transitGatewayVpcAttachment2, &transitGatewayVpcAttachment3),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "2"),
				),
			},
		},
	})
}

func testAccTransitGatewayVPCAttachment_Tags(t *testing.T) {
	var transitGatewayVpcAttachment1, transitGatewayVpcAttachment2, transitGatewayVpcAttachment3 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment1),
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
				Config: testAccTransitGatewayVPCAttachmentConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment2),
					testAccCheckTransitGatewayVPCAttachmentNotRecreated(&transitGatewayVpcAttachment1, &transitGatewayVpcAttachment2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment3),
					testAccCheckTransitGatewayVPCAttachmentNotRecreated(&transitGatewayVpcAttachment2, &transitGatewayVpcAttachment3),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccTransitGatewayVPCAttachment_TransitGatewayDefaultRouteTableAssociationAndPropagationDisabled(t *testing.T) {
	var transitGateway1 ec2.TransitGateway
	var transitGatewayVpcAttachment1 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_defaultRouteTableAssociationAndPropagationDisabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(transitGatewayResourceName, &transitGateway1),
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment1),
					testAccCheckTransitGatewayAssociationDefaultRouteTableAttachmentNotAssociated(&transitGateway1, &transitGatewayVpcAttachment1),
					testAccCheckTransitGatewayPropagationDefaultRouteTableAttachmentNotPropagated(&transitGateway1, &transitGatewayVpcAttachment1),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_association", "false"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_propagation", "false"),
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

func testAccTransitGatewayVPCAttachment_TransitGatewayDefaultRouteTableAssociation(t *testing.T) {
	var transitGateway1, transitGateway2, transitGateway3 ec2.TransitGateway
	var transitGatewayVpcAttachment1, transitGatewayVpcAttachment2, transitGatewayVpcAttachment3 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_defaultRouteTableAssociation(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(transitGatewayResourceName, &transitGateway1),
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment1),
					testAccCheckTransitGatewayAssociationDefaultRouteTableAttachmentNotAssociated(&transitGateway1, &transitGatewayVpcAttachment1),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_association", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_defaultRouteTableAssociation(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(transitGatewayResourceName, &transitGateway2),
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment2),
					testAccCheckTransitGatewayVPCAttachmentNotRecreated(&transitGatewayVpcAttachment1, &transitGatewayVpcAttachment2),
					testAccCheckTransitGatewayAssociationDefaultRouteTableAttachmentAssociated(&transitGateway2, &transitGatewayVpcAttachment2),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_association", "true"),
				),
			},
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_defaultRouteTableAssociation(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(transitGatewayResourceName, &transitGateway3),
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment3),
					testAccCheckTransitGatewayVPCAttachmentNotRecreated(&transitGatewayVpcAttachment2, &transitGatewayVpcAttachment3),
					testAccCheckTransitGatewayAssociationDefaultRouteTableAttachmentNotAssociated(&transitGateway3, &transitGatewayVpcAttachment3),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_association", "false"),
				),
			},
		},
	})
}

func testAccTransitGatewayVPCAttachment_TransitGatewayDefaultRouteTablePropagation(t *testing.T) {
	var transitGateway1, transitGateway2, transitGateway3 ec2.TransitGateway
	var transitGatewayVpcAttachment1, transitGatewayVpcAttachment2, transitGatewayVpcAttachment3 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_defaultRouteTablePropagation(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(transitGatewayResourceName, &transitGateway1),
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment1),
					testAccCheckTransitGatewayPropagationDefaultRouteTableAttachmentNotPropagated(&transitGateway1, &transitGatewayVpcAttachment1),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_propagation", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_defaultRouteTablePropagation(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(transitGatewayResourceName, &transitGateway2),
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment2),
					testAccCheckTransitGatewayVPCAttachmentNotRecreated(&transitGatewayVpcAttachment1, &transitGatewayVpcAttachment2),
					testAccCheckTransitGatewayPropagationDefaultRouteTableAttachmentPropagated(&transitGateway2, &transitGatewayVpcAttachment2),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_propagation", "true"),
				),
			},
			{
				Config: testAccTransitGatewayVPCAttachmentConfig_defaultRouteTablePropagation(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(transitGatewayResourceName, &transitGateway3),
					testAccCheckTransitGatewayVPCAttachmentExists(resourceName, &transitGatewayVpcAttachment3),
					testAccCheckTransitGatewayVPCAttachmentNotRecreated(&transitGatewayVpcAttachment2, &transitGatewayVpcAttachment3),
					testAccCheckTransitGatewayPropagationDefaultRouteTableAttachmentNotPropagated(&transitGateway3, &transitGatewayVpcAttachment3),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_propagation", "false"),
				),
			},
		},
	})
}

func testAccCheckTransitGatewayVPCAttachmentExists(n string, v *ec2.TransitGatewayVpcAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 Transit Gateway VPC Attachment ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		output, err := tfec2.FindTransitGatewayVPCAttachmentByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckTransitGatewayVPCAttachmentDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ec2_transit_gateway_vpc_attachment" {
			continue
		}

		_, err := tfec2.FindTransitGatewayVPCAttachmentByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("EC2 Transit Gateway VPC Attachment %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckTransitGatewayVPCAttachmentNotRecreated(i, j *ec2.TransitGatewayVpcAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(i.TransitGatewayAttachmentId) != aws.StringValue(j.TransitGatewayAttachmentId) {
			return errors.New("EC2 Transit Gateway VPC Attachment was recreated")
		}

		return nil
	}
}

func testAccTransitGatewayVPCAttachmentConfig_base(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccTransitGatewayVPCAttachmentConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccTransitGatewayVPCAttachmentConfig_base(rName), `
resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id
}
`)
}

func testAccTransitGatewayVPCAttachmentConfig_applianceModeSupport(rName, appModeSupport string) string {
	return acctest.ConfigCompose(testAccTransitGatewayVPCAttachmentConfig_base(rName), fmt.Sprintf(`
resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  appliance_mode_support = %[2]q
  subnet_ids             = [aws_subnet.test.id]
  transit_gateway_id     = aws_ec2_transit_gateway.test.id
  vpc_id                 = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName, appModeSupport))
}

func testAccTransitGatewayVPCAttachmentConfig_dnsSupport(rName, dnsSupport string) string {
	return acctest.ConfigCompose(testAccTransitGatewayVPCAttachmentConfig_base(rName), fmt.Sprintf(`
resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  dns_support        = %[2]q
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName, dnsSupport))
}

func testAccTransitGatewayVPCAttachmentConfig_ipv6Support(rName, ipv6Support string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  assign_generated_ipv6_cidr_block = true
  cidr_block                       = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  ipv6_cidr_block   = cidrsubnet(aws_vpc.test.ipv6_cidr_block, 8, 1)
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  ipv6_support       = %[2]q
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName, ipv6Support))
}

func testAccTransitGatewayVPCAttachmentConfig_sharedTransitGateway(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigAlternateAccountProvider(),
		acctest.ConfigAvailableAZsNoOptInDefaultExclude(),
		fmt.Sprintf(`
data "aws_organizations_organization" "test" {}

resource "aws_ec2_transit_gateway" "test" {
  provider = "awsalternate"

  tags = {
    Name = %[1]q
  }
}

resource "aws_ram_resource_share" "test" {
  provider = "awsalternate"

  name = %[1]q
}

resource "aws_ram_resource_association" "test" {
  provider = "awsalternate"

  resource_arn       = aws_ec2_transit_gateway.test.arn
  resource_share_arn = aws_ram_resource_share.test.id
}

resource "aws_ram_principal_association" "test" {
  provider = "awsalternate"

  principal          = data.aws_organizations_organization.test.arn
  resource_share_arn = aws_ram_resource_share.test.id
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  depends_on = [aws_ram_principal_association.test, aws_ram_resource_association.test]

  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccTransitGatewayVPCAttachmentConfig_subnetIds1(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  count = "2"

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = "10.0.${count.index}.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = [aws_subnet.test[0].id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccTransitGatewayVPCAttachmentConfig_subnetIds2(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  count = "2"

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = "10.0.${count.index}.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = aws_subnet.test[*].id
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccTransitGatewayVPCAttachmentConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccTransitGatewayVPCAttachmentConfig_base(rName), fmt.Sprintf(`
resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    %[1]q = %[2]q
  }
}
`, tagKey1, tagValue1))
}

func testAccTransitGatewayVPCAttachmentConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccTransitGatewayVPCAttachmentConfig_base(rName), fmt.Sprintf(`
resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccTransitGatewayVPCAttachmentConfig_defaultRouteTableAssociationAndPropagationDisabled(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  default_route_table_association = "disable"
  default_route_table_propagation = "disable"

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids                                      = [aws_subnet.test.id]
  transit_gateway_default_route_table_association = false
  transit_gateway_default_route_table_propagation = false
  transit_gateway_id                              = aws_ec2_transit_gateway.test.id
  vpc_id                                          = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccTransitGatewayVPCAttachmentConfig_defaultRouteTableAssociation(rName string, transitGatewayDefaultRouteTableAssociation bool) string {
	return acctest.ConfigCompose(testAccTransitGatewayVPCAttachmentConfig_base(rName), fmt.Sprintf(`
resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids                                      = [aws_subnet.test.id]
  transit_gateway_default_route_table_association = %[2]t
  transit_gateway_id                              = aws_ec2_transit_gateway.test.id
  vpc_id                                          = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName, transitGatewayDefaultRouteTableAssociation))
}

func testAccTransitGatewayVPCAttachmentConfig_defaultRouteTablePropagation(rName string, transitGatewayDefaultRouteTablePropagation bool) string {
	return acctest.ConfigCompose(testAccTransitGatewayVPCAttachmentConfig_base(rName), fmt.Sprintf(`
resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids                                      = [aws_subnet.test.id]
  transit_gateway_default_route_table_propagation = %[2]t
  transit_gateway_id                              = aws_ec2_transit_gateway.test.id
  vpc_id                                          = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName, transitGatewayDefaultRouteTablePropagation))
}
