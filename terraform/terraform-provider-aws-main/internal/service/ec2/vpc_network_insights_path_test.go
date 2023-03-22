package ec2_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccVPCNetworkInsightsPath_basic(t *testing.T) {
	resourceName := "aws_ec2_network_insights_path.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInsightsPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsPathConfig_basic(rName, "tcp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkInsightsPathExists(resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile(`network-insights-path/.+$`)),
					resource.TestCheckResourceAttrPair(resourceName, "destination", "aws_network_interface.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "destination_ip", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttrPair(resourceName, "source", "aws_network_interface.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "source_ip", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccVPCNetworkInsightsPath_disappears(t *testing.T) {
	resourceName := "aws_ec2_network_insights_path.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInsightsPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsPathConfig_basic(rName, "udp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsPathExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceNetworkInsightsPath(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVPCNetworkInsightsPath_tags(t *testing.T) {
	resourceName := "aws_ec2_network_insights_path.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInsightsPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsPathConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsPathExists(resourceName),
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
				Config: testAccVPCNetworkInsightsPathConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsPathExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccVPCNetworkInsightsPathConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsPathExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccVPCNetworkInsightsPath_sourceIP(t *testing.T) {
	resourceName := "aws_ec2_network_insights_path.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInsightsPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsPathConfig_sourceIP(rName, "1.1.1.1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsPathExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "source_ip", "1.1.1.1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVPCNetworkInsightsPathConfig_sourceIP(rName, "8.8.8.8"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsPathExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "source_ip", "8.8.8.8"),
				),
			},
		},
	})
}

func TestAccVPCNetworkInsightsPath_destinationIP(t *testing.T) {
	resourceName := "aws_ec2_network_insights_path.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInsightsPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsPathConfig_destinationIP(rName, "1.1.1.1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsPathExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_ip", "1.1.1.1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVPCNetworkInsightsPathConfig_destinationIP(rName, "8.8.8.8"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsPathExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_ip", "8.8.8.8"),
				),
			},
		},
	})
}

func TestAccVPCNetworkInsightsPath_destinationPort(t *testing.T) {
	resourceName := "aws_ec2_network_insights_path.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInsightsPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsPathConfig_destinationPort(rName, 80),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsPathExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_port", "80"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVPCNetworkInsightsPathConfig_destinationPort(rName, 443),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsPathExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_port", "443"),
				),
			},
		},
	})
}

func testAccCheckNetworkInsightsPathExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 Network Insights Path ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		_, err := tfec2.FindNetworkInsightsPathByID(context.Background(), conn, rs.Primary.ID)

		return err
	}
}

func testAccCheckNetworkInsightsPathDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ec2_network_insights_path" {
			continue
		}

		_, err := tfec2.FindNetworkInsightsPathByID(context.Background(), conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("EC2 Network Insights Path %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccVPCNetworkInsightsPathConfig_basic(rName, protocol string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_network_interface" "test" {
  count = 2

  subnet_id = aws_subnet.test[0].id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_network_insights_path" "test" {
  source      = aws_network_interface.test[0].id
  destination = aws_network_interface.test[1].id
  protocol    = %[2]q
}
`, rName, protocol))
}

func testAccVPCNetworkInsightsPathConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_network_interface" "test" {
  count = 2

  subnet_id = aws_subnet.test[0].id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_network_insights_path" "test" {
  source      = aws_network_interface.test[0].id
  destination = aws_network_interface.test[1].id
  protocol    = "tcp"

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccVPCNetworkInsightsPathConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_network_interface" "test" {
  count = 2

  subnet_id = aws_subnet.test[0].id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_network_insights_path" "test" {
  source      = aws_network_interface.test[0].id
  destination = aws_network_interface.test[1].id
  protocol    = "tcp"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccVPCNetworkInsightsPathConfig_sourceIP(rName, sourceIP string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_interface" "test" {
  subnet_id = aws_subnet.test[0].id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_network_insights_path" "test" {
  source      = aws_internet_gateway.test.id
  destination = aws_network_interface.test.id
  protocol    = "tcp"
  source_ip   = %[2]q

  tags = {
    Name = %[1]q
  }
}
`, rName, sourceIP))
}

func testAccVPCNetworkInsightsPathConfig_destinationIP(rName, destinationIP string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_interface" "test" {
  subnet_id = aws_subnet.test[0].id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_network_insights_path" "test" {
  source         = aws_network_interface.test.id
  destination    = aws_internet_gateway.test.id
  protocol       = "tcp"
  destination_ip = %[2]q

  tags = {
    Name = %[1]q
  }
}
`, rName, destinationIP))
}

func testAccVPCNetworkInsightsPathConfig_destinationPort(rName string, destinationPort int) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_network_interface" "test" {
  count = 2

  subnet_id = aws_subnet.test[0].id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_network_insights_path" "test" {
  source           = aws_network_interface.test[0].id
  destination      = aws_network_interface.test[1].id
  protocol         = "tcp"
  destination_port = %[2]d

  tags = {
    Name = %[1]q
  }
}
`, rName, destinationPort))
}
