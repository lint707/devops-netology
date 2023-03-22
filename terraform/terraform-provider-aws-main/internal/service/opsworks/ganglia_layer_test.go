package opsworks_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/opsworks"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccOpsWorksGangliaLayer_basic(t *testing.T) {
	var v opsworks.Layer
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_ganglia_layer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGangliaLayerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGangliaLayerConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLayerExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", "Ganglia"),
					resource.TestCheckResourceAttrSet(resourceName, "password"),
					resource.TestCheckResourceAttr(resourceName, "url", "/ganglia"),
					resource.TestCheckResourceAttr(resourceName, "username", "opsworks"),
				),
			},
		},
	})
}

// _disappears and _tags for OpsWorks Layers are tested via aws_opsworks_rails_app_layer.

func testAccCheckGangliaLayerDestroy(s *terraform.State) error {
	return testAccCheckLayerDestroy("aws_opsworks_ganglia_layer", s)
}

func testAccGangliaLayerConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccLayerConfig_base(rName), `
resource "aws_opsworks_ganglia_layer" "test" {
  stack_id = aws_opsworks_stack.test.id
  password = "avoid-plaintext-passwords"

  custom_security_group_ids = aws_security_group.test[*].id
}
`)
}
