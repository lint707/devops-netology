package cur_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/costandusagereportservice"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccCURReportDefinitionDataSource_basic(t *testing.T) {
	resourceName := "aws_cur_report_definition.test"
	datasourceName := "data.aws_cur_report_definition.test"

	reportName := sdkacctest.RandomWithPrefix("tf_acc_test")
	bucketName := fmt.Sprintf("tf-test-bucket-%d", sdkacctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, costandusagereportservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckReportDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReportDefinitionDataSourceConfig_basic(reportName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccReportDefinitionCheckExistsDataSource(datasourceName, resourceName),
					resource.TestCheckResourceAttrPair(datasourceName, "report_name", resourceName, "report_name"),
					resource.TestCheckResourceAttrPair(datasourceName, "time_unit", resourceName, "time_unit"),
					resource.TestCheckResourceAttrPair(datasourceName, "compression", resourceName, "compression"),
					resource.TestCheckResourceAttrPair(datasourceName, "additional_schema_elements.#", resourceName, "additional_schema_elements.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "s3_bucket", resourceName, "s3_bucket"),
					resource.TestCheckResourceAttrPair(datasourceName, "s3_prefix", resourceName, "s3_prefix"),
					resource.TestCheckResourceAttrPair(datasourceName, "s3_region", resourceName, "s3_region"),
					resource.TestCheckResourceAttrPair(datasourceName, "additional_artifacts.#", resourceName, "additional_artifacts.#"),
				),
			},
		},
	})
}

func TestAccCURReportDefinitionDataSource_additional(t *testing.T) {
	resourceName := "aws_cur_report_definition.test"
	datasourceName := "data.aws_cur_report_definition.test"

	reportName := sdkacctest.RandomWithPrefix("tf_acc_test")
	bucketName := fmt.Sprintf("tf-test-bucket-%d", sdkacctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, costandusagereportservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckReportDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReportDefinitionDataSourceConfig_additional(reportName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccReportDefinitionCheckExistsDataSource(datasourceName, resourceName),
					resource.TestCheckResourceAttrPair(datasourceName, "report_name", resourceName, "report_name"),
					resource.TestCheckResourceAttrPair(datasourceName, "time_unit", resourceName, "time_unit"),
					resource.TestCheckResourceAttrPair(datasourceName, "compression", resourceName, "compression"),
					resource.TestCheckResourceAttrPair(datasourceName, "additional_schema_elements.#", resourceName, "additional_schema_elements.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "s3_bucket", resourceName, "s3_bucket"),
					resource.TestCheckResourceAttrPair(datasourceName, "s3_prefix", resourceName, "s3_prefix"),
					resource.TestCheckResourceAttrPair(datasourceName, "s3_region", resourceName, "s3_region"),
					resource.TestCheckResourceAttrPair(datasourceName, "additional_artifacts.#", resourceName, "additional_artifacts.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "refresh_closed_reports", resourceName, "refresh_closed_reports"),
					resource.TestCheckResourceAttrPair(datasourceName, "report_versioning", resourceName, "report_versioning"),
				),
			},
		},
	})
}

func testAccReportDefinitionCheckExistsDataSource(datasourceName, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[datasourceName]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", datasourceName)
		}
		_, ok = s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", resourceName)
		}
		return nil
	}
}

func testAccReportDefinitionDataSourceConfig_basic(reportName string, bucketName string) string {
	return acctest.ConfigCompose(
		testAccRegionProviderConfig(),
		fmt.Sprintf(`
data "aws_billing_service_account" "test" {}

data "aws_partition" "current" {}

resource "aws_s3_bucket" "test" {
  bucket        = "%[2]s"
  force_destroy = true
}

resource "aws_s3_bucket_acl" "test" {
  bucket = aws_s3_bucket.test.id
  acl    = "private"
}

resource "aws_s3_bucket_policy" "test" {
  bucket = aws_s3_bucket.test.id

  policy = <<POLICY
{
  "Version": "2008-10-17",
  "Id": "s3policy",
  "Statement": [
    {
      "Sid": "AllowCURBillingACLPolicy",
      "Effect": "Allow",
      "Principal": {
        "AWS": "${data.aws_billing_service_account.test.arn}"
      },
      "Action": [
        "s3:GetBucketAcl",
        "s3:GetBucketPolicy"
      ],
      "Resource": "${aws_s3_bucket.test.arn}"
    },
    {
      "Sid": "AllowCURPutObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": "${data.aws_billing_service_account.test.arn}"
      },
      "Action": "s3:PutObject",
      "Resource": "arn:${data.aws_partition.current.partition}:s3:::${aws_s3_bucket.test.id}/*"
    }
  ]
}
POLICY
}

resource "aws_cur_report_definition" "test" {
  depends_on = [aws_s3_bucket_policy.test] # needed to avoid "ValidationException: Failed to verify customer bucket permission."

  report_name                = "%[1]s"
  time_unit                  = "DAILY"
  format                     = "textORcsv"
  compression                = "GZIP"
  additional_schema_elements = ["RESOURCES"]
  s3_bucket                  = aws_s3_bucket.test.id
  s3_prefix                  = ""
  s3_region                  = aws_s3_bucket.test.region
  additional_artifacts       = ["REDSHIFT", "QUICKSIGHT"]
}

data "aws_cur_report_definition" "test" {
  report_name = aws_cur_report_definition.test.report_name
}
`, reportName, bucketName))
}

func testAccReportDefinitionDataSourceConfig_additional(reportName string, bucketName string) string {
	return acctest.ConfigCompose(
		testAccRegionProviderConfig(),
		fmt.Sprintf(`
data "aws_billing_service_account" "test" {}

data "aws_partition" "current" {}

resource "aws_s3_bucket" "test" {
  bucket        = "%[2]s"
  force_destroy = true
}

resource "aws_s3_bucket_acl" "test" {
  bucket = aws_s3_bucket.test.id
  acl    = "private"
}

resource "aws_s3_bucket_policy" "test" {
  bucket = aws_s3_bucket.test.id

  policy = <<POLICY
{
  "Version": "2008-10-17",
  "Id": "s3policy",
  "Statement": [
    {
      "Sid": "AllowCURBillingACLPolicy",
      "Effect": "Allow",
      "Principal": {
        "AWS": "${data.aws_billing_service_account.test.arn}"
      },
      "Action": [
        "s3:GetBucketAcl",
        "s3:GetBucketPolicy"
      ],
      "Resource": "${aws_s3_bucket.test.arn}"
    },
    {
      "Sid": "AllowCURPutObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": "${data.aws_billing_service_account.test.arn}"
      },
      "Action": "s3:PutObject",
      "Resource": "arn:${data.aws_partition.current.partition}:s3:::${aws_s3_bucket.test.id}/*"
    }
  ]
}
POLICY
}

resource "aws_cur_report_definition" "test" {
  depends_on = [aws_s3_bucket_policy.test] # needed to avoid "ValidationException: Failed to verify customer bucket permission."

  report_name                = "%[1]s"
  time_unit                  = "DAILY"
  format                     = "textORcsv"
  compression                = "GZIP"
  additional_schema_elements = ["RESOURCES"]
  s3_bucket                  = aws_s3_bucket.test.id
  s3_prefix                  = ""
  s3_region                  = aws_s3_bucket.test.region
  additional_artifacts       = ["REDSHIFT", "QUICKSIGHT"]
  refresh_closed_reports     = true
  report_versioning          = "CREATE_NEW_REPORT"
}

data "aws_cur_report_definition" "test" {
  report_name = aws_cur_report_definition.test.report_name
}
`, reportName, bucketName))
}
