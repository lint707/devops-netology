package wafregional_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfwafregional "github.com/hashicorp/terraform-provider-aws/internal/service/wafregional"
)

func TestAccWAFRegionalRuleGroup_basic(t *testing.T) {
	var rule waf.Rule
	var group waf.RuleGroup
	var idx int

	ruleName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(wafregional.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_basic(ruleName, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleExists("aws_wafregional_rule.test", &rule),
					testAccCheckRuleGroupExists(resourceName, &group),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "waf-regional", regexp.MustCompile(`rulegroup/.+`)),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric_name", groupName),
					computeActivatedRuleWithRuleId(&rule, "COUNT", 50, &idx),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "activated_rule.*", map[string]string{
						"action.0.type": "COUNT",
						"priority":      "50",
						"type":          waf.WafRuleTypeRegular,
					}),
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

func TestAccWAFRegionalRuleGroup_tags(t *testing.T) {
	var rule waf.Rule
	var group waf.RuleGroup

	ruleName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(wafregional.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_tags1(ruleName, groupName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleExists("aws_wafregional_rule.test", &rule),
					testAccCheckRuleGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRuleGroupConfig_tags2(ruleName, groupName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleExists("aws_wafregional_rule.test", &rule),
					testAccCheckRuleGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccRuleGroupConfig_tags1(ruleName, groupName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleExists("aws_wafregional_rule.test", &rule),
					testAccCheckRuleGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccWAFRegionalRuleGroup_changeNameForceNew(t *testing.T) {
	var before, after waf.RuleGroup

	ruleName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	newGroupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(wafregional.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_basic(ruleName, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleGroupExists(resourceName, &before),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric_name", groupName),
				),
			},
			{
				Config: testAccRuleGroupConfig_basic(ruleName, newGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleGroupExists(resourceName, &after),
					resource.TestCheckResourceAttr(resourceName, "name", newGroupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric_name", newGroupName),
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

func TestAccWAFRegionalRuleGroup_disappears(t *testing.T) {
	var group waf.RuleGroup
	ruleName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(wafregional.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_basic(ruleName, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleGroupExists(resourceName, &group),
					testAccCheckRuleGroupDisappears(&group),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccWAFRegionalRuleGroup_changeActivatedRules(t *testing.T) {
	var rule0, rule1, rule2, rule3 waf.Rule
	var groupBefore, groupAfter waf.RuleGroup
	var idx0, idx1, idx2, idx3 int

	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	ruleName1 := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	ruleName2 := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	ruleName3 := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(wafregional.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_basic(ruleName1, groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRuleExists("aws_wafregional_rule.test", &rule0),
					testAccCheckRuleGroupExists(resourceName, &groupBefore),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", "1"),
					computeActivatedRuleWithRuleId(&rule0, "COUNT", 50, &idx0),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "activated_rule.*", map[string]string{
						"action.0.type": "COUNT",
						"priority":      "50",
						"type":          waf.WafRuleTypeRegular,
					}),
				),
			},
			{
				Config: testAccRuleGroupConfig_changeActivateds(ruleName1, ruleName2, ruleName3, groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", "3"),
					testAccCheckRuleGroupExists(resourceName, &groupAfter),

					testAccCheckRuleExists("aws_wafregional_rule.test", &rule1),
					computeActivatedRuleWithRuleId(&rule1, "BLOCK", 10, &idx1),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "activated_rule.*", map[string]string{
						"action.0.type": "BLOCK",
						"priority":      "10",
						"type":          waf.WafRuleTypeRegular,
					}),

					testAccCheckRuleExists("aws_wafregional_rule.test2", &rule2),
					computeActivatedRuleWithRuleId(&rule2, "COUNT", 1, &idx2),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "activated_rule.*", map[string]string{
						"action.0.type": "COUNT",
						"priority":      "1",
						"type":          waf.WafRuleTypeRegular,
					}),

					testAccCheckRuleExists("aws_wafregional_rule.test3", &rule3),
					computeActivatedRuleWithRuleId(&rule3, "BLOCK", 15, &idx3),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "activated_rule.*", map[string]string{
						"action.0.type": "BLOCK",
						"priority":      "15",
						"type":          waf.WafRuleTypeRegular,
					}),
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

func TestAccWAFRegionalRuleGroup_noActivatedRules(t *testing.T) {
	var group waf.RuleGroup
	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(wafregional.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_noActivateds(groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRuleGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", "0"),
				),
			},
		},
	})
}

func testAccCheckRuleGroupDisappears(group *waf.RuleGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFRegionalConn
		region := acctest.Provider.Meta().(*conns.AWSClient).Region

		rResp, err := conn.ListActivatedRulesInRuleGroup(&waf.ListActivatedRulesInRuleGroupInput{
			RuleGroupId: group.RuleGroupId,
		})
		if err != nil {
			return fmt.Errorf("error listing activated rules in WAF Regional Rule Group (%s): %s", aws.StringValue(group.RuleGroupId), err)
		}

		wr := tfwafregional.NewRetryer(conn, region)
		_, err = wr.RetryWithToken(func(token *string) (interface{}, error) {
			req := &waf.UpdateRuleGroupInput{
				ChangeToken: token,
				RuleGroupId: group.RuleGroupId,
			}

			for _, rule := range rResp.ActivatedRules {
				rule := &waf.RuleGroupUpdate{
					Action:        aws.String("DELETE"),
					ActivatedRule: rule,
				}
				req.Updates = append(req.Updates, rule)
			}

			return conn.UpdateRuleGroup(req)
		})
		if err != nil {
			return fmt.Errorf("Error Updating WAF Regional Rule Group: %s", err)
		}

		_, err = wr.RetryWithToken(func(token *string) (interface{}, error) {
			opts := &waf.DeleteRuleGroupInput{
				ChangeToken: token,
				RuleGroupId: group.RuleGroupId,
			}
			return conn.DeleteRuleGroup(opts)
		})
		if err != nil {
			return fmt.Errorf("Error Deleting WAF Regional Rule Group: %s", err)
		}
		return nil
	}
}

func testAccCheckRuleGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_wafregional_rule_group" {
			continue
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFRegionalConn
		resp, err := conn.GetRuleGroup(&waf.GetRuleGroupInput{
			RuleGroupId: aws.String(rs.Primary.ID),
		})

		if err == nil {
			if *resp.RuleGroup.RuleGroupId == rs.Primary.ID {
				return fmt.Errorf("WAF Regional Rule Group %s still exists", rs.Primary.ID)
			}
		}

		if tfawserr.ErrCodeEquals(err, wafregional.ErrCodeWAFNonexistentItemException) {
			return nil
		}

		return err
	}

	return nil
}

func testAccCheckRuleGroupExists(n string, group *waf.RuleGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No WAF Regional Rule Group ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFRegionalConn
		resp, err := conn.GetRuleGroup(&waf.GetRuleGroupInput{
			RuleGroupId: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		if *resp.RuleGroup.RuleGroupId == rs.Primary.ID {
			*group = *resp.RuleGroup
			return nil
		}

		return fmt.Errorf("WAF Regional Rule Group (%s) not found", rs.Primary.ID)
	}
}

func testAccRuleGroupConfig_basic(ruleName, groupName string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_rule" "test" {
  name        = "%[1]s"
  metric_name = "%[1]s"
}

resource "aws_wafregional_rule_group" "test" {
  name        = "%[2]s"
  metric_name = "%[2]s"

  activated_rule {
    action {
      type = "COUNT"
    }

    priority = 50
    rule_id  = aws_wafregional_rule.test.id
  }
}
`, ruleName, groupName)
}

func testAccRuleGroupConfig_tags1(ruleName, groupName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_rule" "test" {
  name        = %[1]q
  metric_name = %[1]q
}

resource "aws_wafregional_rule_group" "test" {
  name        = %[2]q
  metric_name = %[2]q

  activated_rule {
    action {
      type = "COUNT"
    }

    priority = 50
    rule_id  = aws_wafregional_rule.test.id
  }

  tags = {
    %[3]q = %[4]q
  }
}
`, ruleName, groupName, tagKey1, tagValue1)
}

func testAccRuleGroupConfig_tags2(ruleName, groupName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_rule" "test" {
  name        = %[1]q
  metric_name = %[1]q
}

resource "aws_wafregional_rule_group" "test" {
  name        = %[2]q
  metric_name = %[2]q

  activated_rule {
    action {
      type = "COUNT"
    }

    priority = 50
    rule_id  = aws_wafregional_rule.test.id
  }

  tags = {
    %[3]q = %[4]q
    %[5]q = %[6]q
  }
}
`, ruleName, groupName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccRuleGroupConfig_changeActivateds(ruleName1, ruleName2, ruleName3, groupName string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_rule" "test" {
  name        = "%[1]s"
  metric_name = "%[1]s"
}

resource "aws_wafregional_rule" "test2" {
  name        = "%[2]s"
  metric_name = "%[2]s"
}

resource "aws_wafregional_rule" "test3" {
  name        = "%[3]s"
  metric_name = "%[3]s"
}

resource "aws_wafregional_rule_group" "test" {
  name        = "%[4]s"
  metric_name = "%[4]s"

  activated_rule {
    action {
      type = "BLOCK"
    }

    priority = 10
    rule_id  = aws_wafregional_rule.test.id
  }

  activated_rule {
    action {
      type = "COUNT"
    }

    priority = 1
    rule_id  = aws_wafregional_rule.test2.id
  }

  activated_rule {
    action {
      type = "BLOCK"
    }

    priority = 15
    rule_id  = aws_wafregional_rule.test3.id
  }
}
`, ruleName1, ruleName2, ruleName3, groupName)
}

func testAccRuleGroupConfig_noActivateds(groupName string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_rule_group" "test" {
  name        = "%[1]s"
  metric_name = "%[1]s"
}
`, groupName)
}

// computeActivatedRuleWithRuleId calculates index
// which isn't static because ruleId is generated as part of the test
func computeActivatedRuleWithRuleId(rule *waf.Rule, actionType string, priority int, idx *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ruleResource := tfwafregional.ResourceRuleGroup().Schema["activated_rule"].Elem.(*schema.Resource)

		m := map[string]interface{}{
			"action": []interface{}{
				map[string]interface{}{
					"type": actionType,
				},
			},
			"priority": priority,
			"rule_id":  *rule.RuleId,
			"type":     waf.WafRuleTypeRegular,
		}

		f := schema.HashResource(ruleResource)
		*idx = f(m)

		return nil
	}
}
