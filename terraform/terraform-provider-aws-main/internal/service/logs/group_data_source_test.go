package logs_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccLogsGroupDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "data.aws_cloudwatch_log_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "name", "aws_cloudwatch_log_group.test", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "arn", "aws_cloudwatch_log_group.test", "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrPair(resourceName, "tags", "aws_cloudwatch_log_group.test", "tags"),
				),
			},
		},
	})
}

func TestAccLogsGroupDataSource_tags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "data.aws_cloudwatch_log_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig_tags(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "name", "aws_cloudwatch_log_group.test", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "arn", "aws_cloudwatch_log_group.test", "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrPair(resourceName, "tags", "aws_cloudwatch_log_group.test", "tags"),
				),
			},
		},
	})
}

func TestAccLogsGroupDataSource_kms(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "data.aws_cloudwatch_log_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig_kms(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "name", "aws_cloudwatch_log_group.test", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "arn", "aws_cloudwatch_log_group.test", "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_id", "aws_cloudwatch_log_group.test", "kms_key_id"),
					resource.TestCheckResourceAttrPair(resourceName, "tags", "aws_cloudwatch_log_group.test", "tags"),
				),
			},
		},
	})
}

func TestAccLogsGroupDataSource_retention(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "data.aws_cloudwatch_log_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig_retention(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "name", "aws_cloudwatch_log_group.test", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "arn", "aws_cloudwatch_log_group.test", "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrPair(resourceName, "tags", "aws_cloudwatch_log_group.test", "tags"),
					resource.TestCheckResourceAttrPair(resourceName, "retention_in_days", "aws_cloudwatch_log_group.test", "retention_in_days"),
				),
			},
		},
	})
}

func testAccGroupDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource aws_cloudwatch_log_group "test" {
  name = "%s"
}

data aws_cloudwatch_log_group "test" {
  name = aws_cloudwatch_log_group.test.name
}
`, rName)
}

func testAccGroupDataSourceConfig_tags(rName string) string {
	return fmt.Sprintf(`
resource aws_cloudwatch_log_group "test" {
  name = "%s"

  tags = {
    Environment = "Production"
    Foo         = "Bar"
    Empty       = ""
  }
}

data aws_cloudwatch_log_group "test" {
  name = aws_cloudwatch_log_group.test.name
}
`, rName)
}

func testAccGroupDataSourceConfig_kms(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "foo" {
  description             = "Terraform acc test %s"
  deletion_window_in_days = 7

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "kms-tf-1",
  "Statement": [
    {
      "Sid": "Enable IAM User Permissions",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": "kms:*",
      "Resource": "*"
    }
  ]
}
POLICY
}

resource aws_cloudwatch_log_group "test" {
  name       = "%s"
  kms_key_id = aws_kms_key.foo.arn
}

data aws_cloudwatch_log_group "test" {
  name = aws_cloudwatch_log_group.test.name
}
`, rName, rName)
}

func testAccGroupDataSourceConfig_retention(rName string) string {
	return fmt.Sprintf(`
resource aws_cloudwatch_log_group "test" {
  name              = "%s"
  retention_in_days = 365
}

data aws_cloudwatch_log_group "test" {
  name = aws_cloudwatch_log_group.test.name
}
`, rName)
}
