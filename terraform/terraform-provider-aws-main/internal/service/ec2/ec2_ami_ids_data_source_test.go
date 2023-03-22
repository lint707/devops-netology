package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccEC2AMIIDsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAMIIdsDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAMIIDDataSource("data.aws_ami_ids.ubuntu"),
				),
			},
		},
	})
}

func TestAccEC2AMIIDsDataSource_sorted(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAMIIdsDataSourceConfig_sorted(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_ami_ids.test", "ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"data.aws_ami_ids.test", "ids.0",
						"data.aws_ami.amzn_linux_2018_03", "id"),
					resource.TestCheckResourceAttrPair(
						"data.aws_ami_ids.test", "ids.1",
						"data.aws_ami.amzn_linux_2016_09_0", "id"),
				),
			},
			// Make sure when sort_ascending is set, they're sorted in the inverse order
			// it uses the same config / dataset as above so no need to verify the other
			// bits
			{
				Config: testAccAMIIdsDataSourceConfig_sorted(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.aws_ami_ids.test", "ids.0",
						"data.aws_ami.amzn_linux_2016_09_0", "id"),
					resource.TestCheckResourceAttrPair(
						"data.aws_ami_ids.test", "ids.1",
						"data.aws_ami.amzn_linux_2018_03", "id"),
				),
			},
		},
	})
}

const testAccAMIIdsDataSourceConfig_basic = `
data "aws_ami_ids" "ubuntu" {
  owners = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu/images/ubuntu-*-*-amd64-server-*"]
  }
}
`

func testAccAMIIdsDataSourceConfig_sorted(sort_ascending bool) string {
	return fmt.Sprintf(`
data "aws_ami" "amzn_linux_2016_09_0" {
  owners = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-hvm-2016.09.0.20161028-x86_64-gp2"]
  }
}

data "aws_ami" "amzn_linux_2018_03" {
  owners = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-hvm-2018.03.0.20180811-x86_64-gp2"]
  }
}

data "aws_ami_ids" "test" {
  owners = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-hvm-2018.03.0.20180811-x86_64-gp2", "amzn-ami-hvm-2016.09.0.20161028-x86_64-gp2"]
  }

  sort_ascending = "%t"
}
`, sort_ascending)
}
