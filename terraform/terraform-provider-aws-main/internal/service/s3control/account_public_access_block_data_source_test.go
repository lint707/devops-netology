package s3control_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccS3ControlAccountPublicAccessBlockDataSource_basic(t *testing.T) {
	resourceName := "aws_s3_account_public_access_block.test"
	dataSourceName := "data.aws_s3_account_public_access_block.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3control.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountPublicAccessBlockDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "block_public_acls", dataSourceName, "block_public_acls"),
					resource.TestCheckResourceAttrPair(resourceName, "block_public_policy", dataSourceName, "block_public_policy"),
					resource.TestCheckResourceAttrPair(resourceName, "ignore_public_acls", dataSourceName, "ignore_public_acls"),
					resource.TestCheckResourceAttrPair(resourceName, "restrict_public_buckets", dataSourceName, "restrict_public_buckets"),
				),
			},
		},
	})
}

func testAccAccountPublicAccessBlockDataSourceConfig_base() string {
	return `
resource "aws_s3_account_public_access_block" "test" {
  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}
`
}

func testAccAccountPublicAccessBlockDataSourceConfig_basic() string {
	return acctest.ConfigCompose(testAccAccountPublicAccessBlockDataSourceConfig_base(), `
data "aws_s3_account_public_access_block" "test" {
  depends_on = [aws_s3_account_public_access_block.test]
}
`)
}
