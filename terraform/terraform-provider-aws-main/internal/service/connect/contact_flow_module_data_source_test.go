package connect_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/connect"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func testAccContactFlowModuleDataSource_contactFlowModuleID(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	resourceName := "aws_connect_contact_flow_module.test"
	datasourceName := "data.aws_connect_contact_flow_module.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, connect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactFlowModuleDataSourceConfig_id(rName, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(datasourceName, "contact_flow_module_id", resourceName, "contact_flow_module_id"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_id", resourceName, "instance_id"),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(datasourceName, "content", resourceName, "content"),
					resource.TestCheckResourceAttrSet(datasourceName, "state"),
					resource.TestCheckResourceAttrSet(datasourceName, "status"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
				),
			},
		},
	})
}

func testAccContactFlowModuleDataSource_name(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")
	resourceName := "aws_connect_contact_flow_module.test"
	datasourceName := "data.aws_connect_contact_flow_module.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, connect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactFlowModuleDataSourceConfig_name(rName, rName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(datasourceName, "contact_flow_module_id", resourceName, "contact_flow_module_id"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_id", resourceName, "instance_id"),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(datasourceName, "content", resourceName, "content"),
					resource.TestCheckResourceAttrSet(datasourceName, "state"),
					resource.TestCheckResourceAttrSet(datasourceName, "status"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
				),
			},
		},
	})
}

func testAccContactFlowModuleBaseDataSourceConfig(rName, rName2 string) string {
	return fmt.Sprintf(`
resource "aws_connect_instance" "test" {
  identity_management_type = "CONNECT_MANAGED"
  inbound_calls_enabled    = true
  instance_alias           = %[1]q
  outbound_calls_enabled   = true
}

resource "aws_connect_contact_flow_module" "test" {
  instance_id = aws_connect_instance.test.id
  name        = %[2]q
  description = "Test Contact Flow Module Description"
  content     = file("./test-fixtures/connect_contact_flow_module.json")

  tags = {
    "Name"        = "Test Contact Flow Module",
    "Application" = "Terraform",
    "Method"      = "Create"
  }
}
    `, rName, rName2)
}

func testAccContactFlowModuleDataSourceConfig_id(rName, rName2 string) string {
	return acctest.ConfigCompose(
		testAccContactFlowModuleBaseDataSourceConfig(rName, rName2),
		`
data "aws_connect_contact_flow_module" "test" {
  instance_id            = aws_connect_instance.test.id
  contact_flow_module_id = aws_connect_contact_flow_module.test.contact_flow_module_id
}
`)
}

func testAccContactFlowModuleDataSourceConfig_name(rName, rName2 string) string {
	return acctest.ConfigCompose(
		testAccContactFlowModuleBaseDataSourceConfig(rName, rName2),
		`
data "aws_connect_contact_flow_module" "test" {
  instance_id = aws_connect_instance.test.id
  name        = aws_connect_contact_flow_module.test.name
}
`)
}
