package ecs_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccECSServiceDataSource_basic(t *testing.T) {
	dataSourceName := "data.aws_ecs_service.test"
	resourceName := "aws_ecs_service.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ecs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "desired_count", dataSourceName, "desired_count"),
					resource.TestCheckResourceAttrPair(resourceName, "launch_type", dataSourceName, "launch_type"),
					resource.TestCheckResourceAttrPair(resourceName, "scheduling_strategy", dataSourceName, "scheduling_strategy"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "service_name"),
					resource.TestCheckResourceAttrPair(resourceName, "task_definition", dataSourceName, "task_definition"),
					resource.TestCheckResourceAttrPair(resourceName, "tags", dataSourceName, "tags"),
				),
			},
		},
	})
}

func testAccServiceDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_ecs_cluster" "test" {
  name = %[1]q
}

resource "aws_ecs_task_definition" "test" {
  family = %[1]q

  container_definitions = <<DEFINITION
[
  {
    "cpu": 128,
    "essential": true,
    "image": "mongo:latest",
    "memory": 128,
    "memoryReservation": 64,
    "name": "mongodb"
  }
]
DEFINITION
}

resource "aws_ecs_service" "test" {
  name            = "mongodb"
  cluster         = aws_ecs_cluster.test.id
  task_definition = aws_ecs_task_definition.test.arn
  desired_count   = 1
  tags = {
    Name = %[1]q
  }
}

data "aws_ecs_service" "test" {
  service_name = aws_ecs_service.test.name
  cluster_arn  = aws_ecs_cluster.test.arn
}
`, rName)
}
