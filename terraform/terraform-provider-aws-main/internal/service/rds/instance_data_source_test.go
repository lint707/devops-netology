package rds_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccRDSInstanceDataSource_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_db_instance.test"
	resourceName := "aws_db_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "address", resourceName, "address"),
					resource.TestCheckResourceAttrPair(dataSourceName, "allocated_storage", resourceName, "allocated_storage"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auto_minor_version_upgrade", resourceName, "auto_minor_version_upgrade"),
					resource.TestCheckResourceAttrPair(dataSourceName, "db_instance_arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "db_instance_class", resourceName, "instance_class"),
					resource.TestCheckResourceAttrPair(dataSourceName, "db_name", resourceName, "db_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "db_subnet_group", resourceName, "db_subnet_group_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enabled_cloudwatch_logs_exports.#", resourceName, "enabled_cloudwatch_logs_exports.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "endpoint", resourceName, "endpoint"),
					resource.TestCheckResourceAttrPair(dataSourceName, "engine", resourceName, "engine"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hosted_zone_id", resourceName, "hosted_zone_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "master_username", resourceName, "username"),
					resource.TestCheckResourceAttrPair(dataSourceName, "multi_az", resourceName, "multi_az"),
					resource.TestCheckResourceAttrPair(dataSourceName, "network_type", resourceName, "network_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "port", resourceName, "port"),
					resource.TestCheckResourceAttrPair(dataSourceName, "resource_id", resourceName, "resource_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tags.%", resourceName, "tags.%"),
				),
			},
		},
	})
}

func testAccInstanceDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccInstanceConfig_orderableClassMariadb(),
		testAccInstanceConfig_baseVPC(rName),
		fmt.Sprintf(`
resource "aws_db_instance" "test" {
  allocated_storage       = 10
  backup_retention_period = 0
  db_subnet_group_name    = aws_db_subnet_group.test.name
  engine                  = data.aws_rds_engine_version.default.engine
  engine_version          = data.aws_rds_engine_version.default.version
  identifier              = %[1]q
  instance_class          = data.aws_rds_orderable_db_instance.test.instance_class
  name                    = "test"
  skip_final_snapshot     = true
  password                = "avoid-plaintext-passwords"
  username                = "tfacctest"

  enabled_cloudwatch_logs_exports = [
    "audit",
    "error",
  ]

  tags = {
    Name = %[1]q
  }
}

data "aws_db_instance" "test" {
  db_instance_identifier = aws_db_instance.test.identifier
}
`, rName))
}
