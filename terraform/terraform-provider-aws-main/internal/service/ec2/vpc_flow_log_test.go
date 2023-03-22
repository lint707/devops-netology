package ec2_test

import (
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

func TestAccVPCFlowLog_vpcID(t *testing.T) {
	var flowLog ec2.FlowLog
	cloudwatchLogGroupResourceName := "aws_cloudwatch_log_group.test"
	iamRoleResourceName := "aws_iam_role.test"
	resourceName := "aws_flow_log.test"
	vpcResourceName := "aws_vpc.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_id(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile(`vpc-flow-log/fl-.+`)),
					resource.TestCheckResourceAttrPair(resourceName, "iam_role_arn", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "log_destination", ""),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", "cloud-watch-logs"),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_name", cloudwatchLogGroupResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "max_aggregation_interval", "600"),
					resource.TestCheckResourceAttr(resourceName, "traffic_type", "ALL"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", vpcResourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccVPCFlowLogConfig_destinationTypeCloudWatchLogs(rName),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccVPCFlowLog_logFormat(t *testing.T) {
	var flowLog ec2.FlowLog
	resourceName := "aws_flow_log.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	logFormat := "${version} ${vpc-id} ${subnet-id}"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_format(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					resource.TestCheckResourceAttr(resourceName, "log_format", logFormat),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccVPCFlowLogConfig_destinationTypeCloudWatchLogs(rName),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccVPCFlowLog_subnetID(t *testing.T) {
	var flowLog ec2.FlowLog
	cloudwatchLogGroupResourceName := "aws_cloudwatch_log_group.test"
	iamRoleResourceName := "aws_iam_role.test"
	resourceName := "aws_flow_log.test"
	subnetResourceName := "aws_subnet.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_subnetID(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					resource.TestCheckResourceAttrPair(resourceName, "iam_role_arn", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "log_destination", ""),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", "cloud-watch-logs"),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_name", cloudwatchLogGroupResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "max_aggregation_interval", "600"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", subnetResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "traffic_type", "ALL"),
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

func TestAccVPCFlowLog_transitGatewayID(t *testing.T) {
	var flowLog ec2.FlowLog
	cloudwatchLogGroupResourceName := "aws_cloudwatch_log_group.test"
	iamRoleResourceName := "aws_iam_role.test"
	resourceName := "aws_flow_log.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_transitGatewayId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile(`vpc-flow-log/fl-.+`)),
					resource.TestCheckResourceAttrPair(resourceName, "iam_role_arn", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "log_destination", ""),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", "cloud-watch-logs"),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_name", cloudwatchLogGroupResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "max_aggregation_interval", "60"),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_id", transitGatewayResourceName, "id"),
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

func TestAccVPCFlowLog_transitGatewayAttachmentID(t *testing.T) {
	var flowLog ec2.FlowLog
	cloudwatchLogGroupResourceName := "aws_cloudwatch_log_group.test"
	iamRoleResourceName := "aws_iam_role.test"
	resourceName := "aws_flow_log.test"
	transitGatewayAttachmentResourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_transitGatewayAttachmentId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile(`vpc-flow-log/fl-.+`)),
					resource.TestCheckResourceAttrPair(resourceName, "iam_role_arn", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "log_destination", ""),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", "cloud-watch-logs"),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_name", cloudwatchLogGroupResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "max_aggregation_interval", "60"),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_attachment_id", transitGatewayAttachmentResourceName, "id"),
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

func TestAccVPCFlowLog_LogDestinationType_cloudWatchLogs(t *testing.T) {
	var flowLog ec2.FlowLog
	cloudwatchLogGroupResourceName := "aws_cloudwatch_log_group.test"
	resourceName := "aws_flow_log.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_destinationTypeCloudWatchLogs(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					// We automatically trim :* from ARNs if present
					acctest.CheckResourceAttrRegionalARN(resourceName, "log_destination", "logs", fmt.Sprintf("log-group:%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", "cloud-watch-logs"),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_name", cloudwatchLogGroupResourceName, "name"),
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

func TestAccVPCFlowLog_LogDestinationType_s3(t *testing.T) {
	var flowLog ec2.FlowLog
	s3ResourceName := "aws_s3_bucket.test"
	resourceName := "aws_flow_log.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_destinationTypeS3(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					resource.TestCheckResourceAttrPair(resourceName, "log_destination", s3ResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", "s3"),
					resource.TestCheckResourceAttr(resourceName, "log_group_name", ""),
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

func TestAccVPCFlowLog_LogDestinationTypeS3_invalid(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix("tf-acc-test-flow-log-s3-invalid")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccVPCFlowLogConfig_destinationTypeS3Invalid(rName),
				ExpectError: regexp.MustCompile(`(Access Denied for LogDestination|does not exist)`),
			},
		},
	})
}

func TestAccVPCFlowLog_LogDestinationTypeS3DO_plainText(t *testing.T) {
	var flowLog ec2.FlowLog
	s3ResourceName := "aws_s3_bucket.test"
	resourceName := "aws_flow_log.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_destinationTypeS3DOPlainText(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					resource.TestCheckResourceAttrPair(resourceName, "log_destination", s3ResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", "s3"),
					resource.TestCheckResourceAttr(resourceName, "log_group_name", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_options.0.file_format", "plain-text"),
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

func TestAccVPCFlowLog_LogDestinationTypeS3DOPlainText_hiveCompatible(t *testing.T) {
	var flowLog ec2.FlowLog
	s3ResourceName := "aws_s3_bucket.test"
	resourceName := "aws_flow_log.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_destinationTypeS3DOPlainTextHiveCompatiblePerHour(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					resource.TestCheckResourceAttrPair(resourceName, "log_destination", s3ResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", "s3"),
					resource.TestCheckResourceAttr(resourceName, "log_group_name", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_options.0.file_format", "plain-text"),
					resource.TestCheckResourceAttr(resourceName, "destination_options.0.hive_compatible_partitions", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination_options.0.per_hour_partition", "true"),
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

func TestAccVPCFlowLog_LogDestinationTypeS3DO_parquet(t *testing.T) {
	var flowLog ec2.FlowLog
	s3ResourceName := "aws_s3_bucket.test"
	resourceName := "aws_flow_log.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_destinationTypeS3DOParquet(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					resource.TestCheckResourceAttrPair(resourceName, "log_destination", s3ResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", "s3"),
					resource.TestCheckResourceAttr(resourceName, "log_group_name", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_options.0.file_format", "parquet"),
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

func TestAccVPCFlowLog_LogDestinationTypeS3DOParquet_hiveCompatible(t *testing.T) {
	var flowLog ec2.FlowLog
	s3ResourceName := "aws_s3_bucket.test"
	resourceName := "aws_flow_log.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_destinationTypeS3DOParquetHiveCompatible(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					resource.TestCheckResourceAttrPair(resourceName, "log_destination", s3ResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", "s3"),
					resource.TestCheckResourceAttr(resourceName, "log_group_name", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_options.0.file_format", "parquet"),
					resource.TestCheckResourceAttr(resourceName, "destination_options.0.hive_compatible_partitions", "true"),
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

func TestAccVPCFlowLog_LogDestinationTypeS3DOParquetHiveCompatible_perHour(t *testing.T) {
	var flowLog ec2.FlowLog
	s3ResourceName := "aws_s3_bucket.test"
	resourceName := "aws_flow_log.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_destinationTypeS3DOParquetHiveCompatiblePerHour(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					resource.TestCheckResourceAttrPair(resourceName, "log_destination", s3ResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", "s3"),
					resource.TestCheckResourceAttr(resourceName, "log_group_name", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_options.0.file_format", "parquet"),
					resource.TestCheckResourceAttr(resourceName, "destination_options.0.hive_compatible_partitions", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination_options.0.per_hour_partition", "true"),
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

func TestAccVPCFlowLog_LogDestinationType_maxAggregationInterval(t *testing.T) {
	var flowLog ec2.FlowLog
	resourceName := "aws_flow_log.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_maxAggregationInterval(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					resource.TestCheckResourceAttr(resourceName, "max_aggregation_interval", "60"),
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

func TestAccVPCFlowLog_tags(t *testing.T) {
	var flowLog ec2.FlowLog
	resourceName := "aws_flow_log.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
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
				Config: testAccVPCFlowLogConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccVPCFlowLogConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccVPCFlowLog_disappears(t *testing.T) {
	var flowLog ec2.FlowLog
	resourceName := "aws_flow_log.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCFlowLogConfig_id(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists(resourceName, &flowLog),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceFlowLog(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckFlowLogExists(n string, v *ec2.FlowLog) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Flow Log ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		output, err := tfec2.FindFlowLogByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckFlowLogDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_flow_log" {
			continue
		}

		_, err := tfec2.FindFlowLogByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Flow Log %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccFlowLogConfigBase(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCFlowLogConfig_destinationTypeCloudWatchLogs(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_flow_log" "test" {
  iam_role_arn         = aws_iam_role.test.arn
  log_destination      = aws_cloudwatch_log_group.test.arn
  log_destination_type = "cloud-watch-logs"
  traffic_type         = "ALL"
  vpc_id               = aws_vpc.test.id
}
`, rName))
}

func testAccVPCFlowLogConfig_destinationTypeS3(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_flow_log" "test" {
  log_destination      = aws_s3_bucket.test.arn
  log_destination_type = "s3"
  traffic_type         = "ALL"
  vpc_id               = aws_vpc.test.id
}
`, rName))
}

func testAccVPCFlowLogConfig_destinationTypeS3Invalid(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), `
data "aws_partition" "current" {}

resource "aws_flow_log" "test" {
  log_destination      = "arn:${data.aws_partition.current.partition}:s3:::does-not-exist"
  log_destination_type = "s3"
  traffic_type         = "ALL"
  vpc_id               = aws_vpc.test.id
}
`)
}

func testAccVPCFlowLogConfig_destinationTypeS3DOPlainText(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_flow_log" "test" {
  log_destination      = aws_s3_bucket.test.arn
  log_destination_type = "s3"
  traffic_type         = "ALL"
  vpc_id               = aws_vpc.test.id
  destination_options {
    file_format = "plain-text"
  }
}
`, rName))
}

func testAccVPCFlowLogConfig_destinationTypeS3DOPlainTextHiveCompatiblePerHour(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_flow_log" "test" {
  log_destination      = aws_s3_bucket.test.arn
  log_destination_type = "s3"
  traffic_type         = "ALL"
  vpc_id               = aws_vpc.test.id
  destination_options {
    file_format                = "plain-text"
    hive_compatible_partitions = true
    per_hour_partition         = true
  }
}
`, rName))
}

func testAccVPCFlowLogConfig_destinationTypeS3DOParquet(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_flow_log" "test" {
  log_destination      = aws_s3_bucket.test.arn
  log_destination_type = "s3"
  traffic_type         = "ALL"
  vpc_id               = aws_vpc.test.id
  destination_options {
    file_format = "parquet"
  }
}
`, rName))
}

func testAccVPCFlowLogConfig_destinationTypeS3DOParquetHiveCompatible(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_flow_log" "test" {
  log_destination      = aws_s3_bucket.test.arn
  log_destination_type = "s3"
  traffic_type         = "ALL"
  vpc_id               = aws_vpc.test.id
  destination_options {
    file_format                = "parquet"
    hive_compatible_partitions = true
  }
}
`, rName))
}

func testAccVPCFlowLogConfig_destinationTypeS3DOParquetHiveCompatiblePerHour(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_flow_log" "test" {
  log_destination      = aws_s3_bucket.test.arn
  log_destination_type = "s3"
  traffic_type         = "ALL"
  vpc_id               = aws_vpc.test.id
  destination_options {
    file_format                = "parquet"
    hive_compatible_partitions = true
    per_hour_partition         = true
  }
}
`, rName))
}

func testAccVPCFlowLogConfig_subnetID(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
resource "aws_subnet" "test" {
  cidr_block = "10.0.1.0/24"
  vpc_id     = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_flow_log" "test" {
  iam_role_arn   = aws_iam_role.test.arn
  log_group_name = aws_cloudwatch_log_group.test.name
  subnet_id      = aws_subnet.test.id
  traffic_type   = "ALL"
}
`, rName))
}

func testAccVPCFlowLogConfig_id(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_flow_log" "test" {
  iam_role_arn   = aws_iam_role.test.arn
  log_group_name = aws_cloudwatch_log_group.test.name
  traffic_type   = "ALL"
  vpc_id         = aws_vpc.test.id
}
`, rName))
}

func testAccVPCFlowLogConfig_format(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_flow_log" "test" {
  log_destination      = aws_s3_bucket.test.arn
  log_destination_type = "s3"
  traffic_type         = "ALL"
  vpc_id               = aws_vpc.test.id
  log_format           = "$${version} $${vpc-id} $${subnet-id}"
}
`, rName))
}

func testAccVPCFlowLogConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_flow_log" "test" {
  iam_role_arn   = aws_iam_role.test.arn
  log_group_name = aws_cloudwatch_log_group.test.name
  traffic_type   = "ALL"
  vpc_id         = aws_vpc.test.id

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccVPCFlowLogConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_flow_log" "test" {
  iam_role_arn   = aws_iam_role.test.arn
  log_group_name = aws_cloudwatch_log_group.test.name
  traffic_type   = "ALL"
  vpc_id         = aws_vpc.test.id

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccVPCFlowLogConfig_maxAggregationInterval(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_flow_log" "test" {
  iam_role_arn   = aws_iam_role.test.arn
  log_group_name = aws_cloudwatch_log_group.test.name
  traffic_type   = "ALL"
  vpc_id         = aws_vpc.test.id

  max_aggregation_interval = 60
}
`, rName))
}

func testAccVPCFlowLogConfig_transitGatewayId(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_flow_log" "test" {
  iam_role_arn             = aws_iam_role.test.arn
  log_group_name           = aws_cloudwatch_log_group.test.name
  max_aggregation_interval = 60
  transit_gateway_id       = aws_ec2_transit_gateway.test.id
}
`, rName))
}

func testAccVPCFlowLogConfig_transitGatewayAttachmentId(rName string) string {
	return acctest.ConfigCompose(testAccFlowLogConfigBase(rName), fmt.Sprintf(`
resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block = "10.0.1.0/24"
  vpc_id     = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id
  subnet_ids         = [aws_subnet.test.id]

  tags = {
    Name = %[1]q
  }
}

data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_flow_log" "test" {
  iam_role_arn                  = aws_iam_role.test.arn
  log_group_name                = aws_cloudwatch_log_group.test.name
  max_aggregation_interval      = 60
  transit_gateway_attachment_id = aws_ec2_transit_gateway_vpc_attachment.test.id
}
`, rName))
}
