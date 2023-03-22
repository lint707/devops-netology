package accessanalyzer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/accessanalyzer"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfaccessanalyzer "github.com/hashicorp/terraform-provider-aws/internal/service/accessanalyzer"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func testAccAnalyzerArchiveRule_basic(t *testing.T) {
	var archiveRule accessanalyzer.ArchiveRuleSummary
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_accessanalyzer_archive_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(accessanalyzer.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, accessanalyzer.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckArchiveRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccArchiveRuleConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckArchiveRuleExists(resourceName, &archiveRule),
					resource.TestCheckResourceAttr(resourceName, "filter.0.criteria", "isPublic"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAnalyzerArchiveRule_updateFilters(t *testing.T) {
	var archiveRule accessanalyzer.ArchiveRuleSummary
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_accessanalyzer_archive_rule.test"

	filters := `
filter {
  criteria = "error"
  exists   = true
}
`

	filtersUpdated := `
filter {
  criteria = "error"
  exists   = true
}

filter {
  criteria = "isPublic"
  eq       = ["false"]
}
`

	filtersRemoved := `
filter {
  criteria = "isPublic"
  eq       = ["true"]
}
`
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(accessanalyzer.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, accessanalyzer.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckArchiveRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccArchiveRuleConfig_updateFilters(rName, filters),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckArchiveRuleExists(resourceName, &archiveRule),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.criteria", "error"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.exists", "true"),
				),
			},
			{
				Config: testAccArchiveRuleConfig_updateFilters(rName, filtersUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckArchiveRuleExists(resourceName, &archiveRule),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.criteria", "error"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.exists", "true"),
					resource.TestCheckResourceAttr(resourceName, "filter.1.criteria", "isPublic"),
					resource.TestCheckResourceAttr(resourceName, "filter.1.eq.0", "false"),
				),
			},
			{
				Config: testAccArchiveRuleConfig_updateFilters(rName, filtersRemoved),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckArchiveRuleExists(resourceName, &archiveRule),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.criteria", "isPublic"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.eq.0", "true"),
				),
			},
		},
	})
}

func testAccAnalyzerArchiveRule_disappears(t *testing.T) {
	var archiveRule accessanalyzer.ArchiveRuleSummary
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_accessanalyzer_archive_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(accessanalyzer.EndpointsID, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, accessanalyzer.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckArchiveRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccArchiveRuleConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckArchiveRuleExists(resourceName, &archiveRule),
					acctest.CheckResourceDisappears(acctest.Provider, tfaccessanalyzer.ResourceArchiveRule(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckArchiveRuleDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).AccessAnalyzerConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_accessanalyzer_archive_rule" {
			continue
		}

		analyzerName, ruleName, err := tfaccessanalyzer.DecodeRuleID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("unable to decode AccessAnalyzer ArchiveRule ID (%s): %s", rs.Primary.ID, err)
		}

		_, err = tfaccessanalyzer.FindArchiveRule(context.Background(), conn, analyzerName, ruleName)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("expected AccessAnalyzer ArchiveRule to be destroyed, %s found", rs.Primary.ID)
	}

	return nil
}

func testAccCheckArchiveRuleExists(name string, archiveRule *accessanalyzer.ArchiveRuleSummary) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AccessAnalyzer ArchiveRule is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).AccessAnalyzerConn
		analyzerName, ruleName, err := tfaccessanalyzer.DecodeRuleID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("unable to decode AccessAnalyzer ArchiveRule ID (%s): %s", rs.Primary.ID, err)
		}

		resp, err := tfaccessanalyzer.FindArchiveRule(context.Background(), conn, analyzerName, ruleName)

		if err != nil {
			return fmt.Errorf("describing AccessAnalyzer ArchiveRule: %s", err.Error())
		}

		*archiveRule = *resp

		return nil
	}
}

func testAccArchiveRuleBaseConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_accessanalyzer_analyzer" "test" {
  analyzer_name = %[1]q
}

`, rName)
}

func testAccArchiveRuleConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccArchiveRuleBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_accessanalyzer_archive_rule" "test" {
  analyzer_name = aws_accessanalyzer_analyzer.test.analyzer_name
  rule_name     = %[1]q

  filter {
    criteria = "isPublic"
    eq       = ["false"]
  }
}
`, rName))
}

func testAccArchiveRuleConfig_updateFilters(rName, filters string) string {
	return acctest.ConfigCompose(
		testAccArchiveRuleBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_accessanalyzer_archive_rule" "test" {
  analyzer_name = aws_accessanalyzer_analyzer.test.analyzer_name
  rule_name     = %[1]q

  %[2]s
}
`, rName, filters))
}
