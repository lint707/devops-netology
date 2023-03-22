package meta_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	tfmeta "github.com/hashicorp/terraform-provider-aws/internal/service/meta"
)

func TestAccMetaARNDataSource_basic(t *testing.T) {
	arn := "arn:aws:rds:eu-west-1:123456789012:db:mysql-db" // lintignore:AWSAT003,AWSAT005
	dataSourceName := "data.aws_arn.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, tfmeta.PseudoServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccARNDataSourceConfig_basic(arn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "account", "123456789012"),
					resource.TestCheckResourceAttr(dataSourceName, "id", arn),
					resource.TestCheckResourceAttr(dataSourceName, "partition", "aws"),
					resource.TestCheckResourceAttr(dataSourceName, "region", "eu-west-1"), // lintignore:AWSAT003
					resource.TestCheckResourceAttr(dataSourceName, "resource", "db:mysql-db"),
					resource.TestCheckResourceAttr(dataSourceName, "service", "rds"),
				),
			},
		},
	})
}

func TestAccMetaARNDataSource_s3Bucket(t *testing.T) {
	arn := "arn:aws:s3:::my_corporate_bucket/Development/*" // lintignore:AWSAT005
	dataSourceName := "data.aws_arn.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, tfmeta.PseudoServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccARNDataSourceConfig_basic(arn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "account", ""),
					resource.TestCheckResourceAttr(dataSourceName, "id", arn),
					resource.TestCheckResourceAttr(dataSourceName, "partition", "aws"),
					resource.TestCheckResourceAttr(dataSourceName, "region", ""),
					resource.TestCheckResourceAttr(dataSourceName, "resource", "my_corporate_bucket/Development/*"),
					resource.TestCheckResourceAttr(dataSourceName, "service", "s3"),
				),
			},
		},
	})
}

func testAccARNDataSourceConfig_basic(arn string) string {
	return fmt.Sprintf(`
data "aws_arn" "test" {
  arn = %[1]q
}
`, arn)
}
