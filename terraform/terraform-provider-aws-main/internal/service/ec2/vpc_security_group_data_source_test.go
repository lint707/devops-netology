package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccVPCSecurityGroupDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccSecurityGroupCheckDataSource("data.aws_security_group.by_id"),
					testAccSecurityGroupCheckDataSource("data.aws_security_group.by_tag"),
					testAccSecurityGroupCheckDataSource("data.aws_security_group.by_filter"),
					testAccSecurityGroupCheckDataSource("data.aws_security_group.by_name"),
				),
			},
		},
	})
}

func testAccSecurityGroupCheckDataSource(dataSourceName string) resource.TestCheckFunc {
	resourceName := "aws_security_group.test"

	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
		resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
		resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
		resource.TestCheckResourceAttrPair(dataSourceName, "tags.%", resourceName, "tags.%"),
		resource.TestCheckResourceAttrPair(dataSourceName, "vpc_id", resourceName, "vpc_id"),
	)
}

func testAccVPCSecurityGroupDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "172.16.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  vpc_id = aws_vpc.test.id
  name   = %[1]q

  tags = {
    Name = %[1]q
  }
}

data "aws_security_group" "by_id" {
  id = aws_security_group.test.id
}

data "aws_security_group" "by_name" {
  name = aws_security_group.test.name
}

data "aws_security_group" "by_tag" {
  tags = {
    Name = aws_security_group.test.tags["Name"]
  }
}

data "aws_security_group" "by_filter" {
  filter {
    name   = "group-name"
    values = [aws_security_group.test.name]
  }
}
`, rName)
}
