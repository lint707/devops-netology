package signer_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/signer"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccSignerSigningProfileDataSource_basic(t *testing.T) {
	dataSourceName := "data.aws_signer_signing_profile.test"
	resourceName := "aws_signer_signing_profile.test"
	rString := sdkacctest.RandString(48)
	profileName := fmt.Sprintf("tf_acc_sp_basic_%s", rString)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckSingerSigningProfile(t, "AWSLambda-SHA384-ECDSA") },
		ErrorCheck:               acctest.ErrorCheck(t, signer.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSigningProfileDataSourceConfig_basic(profileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "platform_id", resourceName, "platform_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "signature_validity_period.value", resourceName, "signature_validity_period.value"),
					resource.TestCheckResourceAttrPair(dataSourceName, "signature_validity_period.type", resourceName, "signature_validity_period.type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "platform_display_name", resourceName, "platform_display_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tags", resourceName, "tags"),
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
				),
			},
		},
	})
}

func testAccSigningProfileDataSourceConfig_basic(profileName string) string {
	return fmt.Sprintf(`
resource "aws_signer_signing_profile" "test" {
  platform_id = "AWSLambda-SHA384-ECDSA"
  name        = "%s"
}

data "aws_signer_signing_profile" "test" {
  name = aws_signer_signing_profile.test.name
}`, profileName)
}
