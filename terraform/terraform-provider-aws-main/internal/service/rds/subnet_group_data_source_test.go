package rds_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccRDSSubnetGroupDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_db_subnet_group.test"
	dataSourceName := "data.aws_db_subnet_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetGroupDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "arn", dataSourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_ids.#", dataSourceName, "subnet_ids.#"),
					resource.TestCheckResourceAttrPair(resourceName, "supported_network_types.#", dataSourceName, "supported_network_types.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_id"),
				),
			},
		},
	})
}

func testAccSubnetGroupDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccSubnetGroupConfig_basic(rName), `
data "aws_db_subnet_group" "test" {
  name = aws_db_subnet_group.test.name
}
`)
}
