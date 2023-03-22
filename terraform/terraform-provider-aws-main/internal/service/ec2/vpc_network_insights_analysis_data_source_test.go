package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccVPCNetworkInsightsAnalysisDataSource_basic(t *testing.T) {
	resourceName := "aws_ec2_network_insights_analysis.test"
	datasourceName := "data.aws_ec2_network_insights_analysis.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkInsightsAnalysisDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "alternate_path_hints.#", resourceName, "alternate_path_hints.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(datasourceName, "explanations.#", resourceName, "explanations.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "filter_in_arns.#", resourceName, "filter_in_arns.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "forward_path_components.#", resourceName, "forward_path_components.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "network_insights_analysis_id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "network_insights_path_id", resourceName, "network_insights_path_id"),
					resource.TestCheckResourceAttrPair(datasourceName, "path_found", resourceName, "path_found"),
					resource.TestCheckResourceAttrPair(datasourceName, "return_path_components.#", resourceName, "return_path_components.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "start_date", resourceName, "start_date"),
					resource.TestCheckResourceAttrPair(datasourceName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(datasourceName, "status_message", resourceName, "status_message"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "warning_message", resourceName, "warning_message"),
				),
			},
		},
	})
}

func testAccVPCNetworkInsightsAnalysisDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccVPCNetworkInsightsAnalysisConfig_base(rName), fmt.Sprintf(`
resource "aws_ec2_network_insights_analysis" "test" {
  network_insights_path_id = aws_ec2_network_insights_path.test.id
  wait_for_completion      = true

  tags = {
    Name = %[1]q
  }
}

data "aws_ec2_network_insights_analysis" "test" {
  network_insights_analysis_id = aws_ec2_network_insights_analysis.test.id
}
`, rName))
}
