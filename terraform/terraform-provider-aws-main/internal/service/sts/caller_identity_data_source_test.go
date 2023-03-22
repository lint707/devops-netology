package sts_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccSTSCallerIdentityDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sts.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCallerIdentityConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckCallerIdentityAccountID("data.aws_caller_identity.current"),
				),
			},
		},
	})
}

const testAccCallerIdentityConfig_basic = `
data "aws_caller_identity" "current" {}
`
