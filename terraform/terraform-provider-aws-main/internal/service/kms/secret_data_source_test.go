package kms_test

import (
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	tfkms "github.com/hashicorp/terraform-provider-aws/internal/service/kms"
)

func TestAccKMSSecretDataSource_removed(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSecretDataSourceConfig_basic,
				ExpectError: regexp.MustCompile(tfkms.SecretRemovedMessage),
			},
		},
	})
}

const testAccSecretDataSourceConfig_basic = `
data "aws_kms_secret" "testing" {
  secret {
    name    = "secret_name"
    payload = "data-source-removed"

    context = {
      name = "value"
    }
  }
}
`
