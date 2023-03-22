package sfn_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/sfn"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccSFNActivityDataSource_StepFunctions_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sfn_activity.test"
	dataName := "data.aws_sfn_activity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sfn.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccActivityDataSourceConfig_checkARN(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "creation_date", dataName, "creation_date"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataName, "name"),
				),
			},
			{
				Config: testAccActivityDataSourceConfig_checkName(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "creation_date", dataName, "creation_date"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataName, "name"),
				),
			},
		},
	})
}

func testAccActivityDataSourceConfig_checkARN(rName string) string {
	return fmt.Sprintf(`
resource aws_sfn_activity "test" {
  name = "%s"
}

data aws_sfn_activity "test" {
  arn = aws_sfn_activity.test.id
}
`, rName)
}

func testAccActivityDataSourceConfig_checkName(rName string) string {
	return fmt.Sprintf(`
resource aws_sfn_activity "test" {
  name = "%s"
}

data aws_sfn_activity "test" {
  name = aws_sfn_activity.test.name
}
`, rName)
}
