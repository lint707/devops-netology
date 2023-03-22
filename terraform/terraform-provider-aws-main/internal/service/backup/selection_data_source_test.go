package backup_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/backup"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccBackupSelectionDataSource_basic(t *testing.T) {
	datasourceName := "data.aws_backup_selection.test"
	resourceName := "aws_backup_selection.test"
	rInt := sdkacctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, backup.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSelectionDataSourceConfig_nonExistent,
				ExpectError: regexp.MustCompile(`Error getting Backup Selection`),
			},
			{
				Config: testAccSelectionDataSourceConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "iam_role_arn", resourceName, "iam_role_arn"),
					resource.TestCheckResourceAttrPair(datasourceName, "resources.#", resourceName, "resources.#"),
				),
			},
		},
	})
}

const testAccSelectionDataSourceConfig_nonExistent = `
data "aws_backup_selection" "test" {
  plan_id      = "tf-acc-test-does-not-exist"
  selection_id = "tf-acc-test-dne"
}
`

func testAccSelectionDataSourceConfig_basic(rInt int) string {
	return fmt.Sprintf(`
data "aws_caller_identity" "current" {}

data "aws_partition" "current" {}

data "aws_region" "current" {}

resource "aws_backup_vault" "test" {
  name = "tf_acc_test_backup_vault_%[1]d"
}

resource "aws_backup_plan" "test" {
  name = "tf_acc_test_backup_plan_%[1]d"

  rule {
    rule_name         = "tf_acc_test_backup_rule_%[1]d"
    target_vault_name = aws_backup_vault.test.name
    schedule          = "cron(0 12 * * ? *)"
  }
}

resource "aws_backup_selection" "test" {
  plan_id      = aws_backup_plan.test.id
  name         = "tf_acc_test_backup_selection_%[1]d"
  iam_role_arn = "arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:role/service-role/AWSBackupDefaultServiceRole"

  selection_tag {
    type  = "STRINGEQUALS"
    key   = "foo"
    value = "bar"
  }

  condition {}
  not_resources = []

  resources = [
    "arn:${data.aws_partition.current.partition}:ec2:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:volume/*"
  ]
}

data "aws_backup_selection" "test" {
  plan_id      = aws_backup_plan.test.id
  selection_id = aws_backup_selection.test.id
}
`, rInt)
}
