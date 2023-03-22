package opsworks_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/opsworks"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccOpsWorksHAProxyLayer_basic(t *testing.T) {
	var v opsworks.Layer
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_haproxy_layer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckHAProxyLayerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHAProxyLayerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLayerExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_method", "OPTIONS"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_url", "/"),
					resource.TestCheckResourceAttr(resourceName, "name", "HAProxy"),
					resource.TestCheckResourceAttr(resourceName, "stats_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "stats_password"),
					resource.TestCheckResourceAttr(resourceName, "stats_url", "/haproxy?stats"),
					resource.TestCheckResourceAttr(resourceName, "stats_user", "opsworks"),
				),
			},
		},
	})
}

// _disappears and _tags for OpsWorks Layers are tested via aws_opsworks_rails_app_layer.

func testAccCheckHAProxyLayerDestroy(s *terraform.State) error {
	return testAccCheckLayerDestroy("aws_opsworks_haproxy_layer", s)
}

func testAccHAProxyLayerConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccLayerConfig_base(rName), `
resource "aws_opsworks_haproxy_layer" "test" {
  stack_id       = aws_opsworks_stack.test.id
  stats_password = "avoid-plaintext-passwords"

  custom_security_group_ids = aws_security_group.test[*].id
}
`)
}
