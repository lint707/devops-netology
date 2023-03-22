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

func TestAccVPCNetworkInsightsAnalysis_basic(t *testing.T) {
	resourceName := "aws_ec2_network_insights_analysis.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInsightsAnalysisDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsAnalysisConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkInsightsAnalysisExists(resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile(`network-insights-analysis/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "filter_in_arns.#", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "network_insights_path_id", "aws_ec2_network_insights_path.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "path_found", "true"),
					acctest.CheckResourceAttrRFC3339(resourceName, "start_date"),
					resource.TestCheckResourceAttr(resourceName, "status", "succeeded"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_completion", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_for_completion"},
			},
		},
	})
}

func TestAccVPCNetworkInsightsAnalysis_disappears(t *testing.T) {
	resourceName := "aws_ec2_network_insights_analysis.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInsightsAnalysisDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsAnalysisConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsAnalysisExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceNetworkInsightsAnalysis(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVPCNetworkInsightsAnalysis_tags(t *testing.T) {
	resourceName := "aws_ec2_network_insights_analysis.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInsightsAnalysisDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsAnalysisConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsAnalysisExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_for_completion"},
			},
			{
				Config: testAccVPCNetworkInsightsAnalysisConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsAnalysisExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccVPCNetworkInsightsAnalysisConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsAnalysisExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccVPCNetworkInsightsAnalysis_filterInARNs(t *testing.T) {
	resourceName := "aws_ec2_network_insights_analysis.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInsightsAnalysisDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsAnalysisConfig_filterInARNs(rName, "vpc-peering-connection/pcx-fakearn1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsAnalysisExists(resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "filter_in_arns.0", "ec2", regexp.MustCompile(`vpc-peering-connection/pcx-fakearn1$`)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_for_completion"},
			},
			{
				Config: testAccVPCNetworkInsightsAnalysisConfig_filterInARNs(rName, "vpc-peering-connection/pcx-fakearn2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsAnalysisExists(resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "filter_in_arns.0", "ec2", regexp.MustCompile(`vpc-peering-connection/pcx-fakearn2$`)),
				),
			},
		},
	})
}

func TestAccVPCNetworkInsightsAnalysis_waitForCompletion(t *testing.T) {
	resourceName := "aws_ec2_network_insights_analysis.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInsightsAnalysisDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsAnalysisConfig_waitForCompletion(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsAnalysisExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "wait_for_completion", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_for_completion"},
			},
			{
				Config: testAccVPCNetworkInsightsAnalysisConfig_waitForCompletion(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInsightsAnalysisExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "wait_for_completion", "true"),
				),
			},
		},
	})
}

func testAccCheckNetworkInsightsAnalysisExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 Network Insights Analysis ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		_, err := tfec2.FindNetworkInsightsAnalysisByID(context.Background(), conn, rs.Primary.ID)

		return err
	}
}

func testAccCheckNetworkInsightsAnalysisDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ec2_network_insights_analysis" {
			continue
		}

		_, err := tfec2.FindNetworkInsightsAnalysisByID(context.Background(), conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("EC2 Network Insights Analysis %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccVPCNetworkInsightsAnalysisConfig_base(rName string) string {
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
    Name = %[1]q
  }
}
`, rName))
}

func testAccVPCNetworkInsightsAnalysisConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccVPCNetworkInsightsAnalysisConfig_base(rName), `
resource "aws_ec2_network_insights_analysis" "test" {
  network_insights_path_id = aws_ec2_network_insights_path.test.id
}
`)
}

func testAccVPCNetworkInsightsAnalysisConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccVPCNetworkInsightsAnalysisConfig_base(rName), fmt.Sprintf(`
resource "aws_ec2_network_insights_analysis" "test" {
  network_insights_path_id = aws_ec2_network_insights_path.test.id

  tags = {
    %[1]q = %[2]q
  }
}
`, tagKey1, tagValue1))
}

func testAccVPCNetworkInsightsAnalysisConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccVPCNetworkInsightsAnalysisConfig_base(rName), fmt.Sprintf(`
resource "aws_ec2_network_insights_analysis" "test" {
  network_insights_path_id = aws_ec2_network_insights_path.test.id

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccVPCNetworkInsightsAnalysisConfig_filterInARNs(rName, arnSuffix string) string {
	return acctest.ConfigCompose(testAccVPCNetworkInsightsAnalysisConfig_base(rName), fmt.Sprintf(`
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
data "aws_partition" "current" {}

resource "aws_ec2_network_insights_analysis" "test" {
  network_insights_path_id = aws_ec2_network_insights_path.test.id
  filter_in_arns           = ["arn:${data.aws_partition.current.partition}:ec2:${data.aws_region.current.name}:${data.aws_caller_identity.current.id}:%[2]s"]

  tags = {
    Name = %[1]q
  }
}
`, rName, arnSuffix))
}

func testAccVPCNetworkInsightsAnalysisConfig_waitForCompletion(rName string, waitForCompletion bool) string {
	return acctest.ConfigCompose(testAccVPCNetworkInsightsAnalysisConfig_base(rName), fmt.Sprintf(`
resource "aws_ec2_network_insights_analysis" "test" {
  network_insights_path_id = aws_ec2_network_insights_path.test.id
  wait_for_completion      = %[2]t

  tags = {
    Name = %[1]q
  }
}
`, rName, waitForCompletion))
}
