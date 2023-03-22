package elbv2_test

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfelbv2 "github.com/hashicorp/terraform-provider-aws/internal/service/elbv2"
)

func TestLBListenerARNFromRuleARN(t *testing.T) {
	cases := []struct {
		name     string
		arn      string
		expected string
	}{
		{
			name:     "valid listener rule arn",
			arn:      "arn:aws:elasticloadbalancing:us-east-1:012345678912:listener-rule/app/name/0123456789abcdef/abcdef0123456789/456789abcedf1234", //lintignore:AWSAT003,AWSAT005
			expected: "arn:aws:elasticloadbalancing:us-east-1:012345678912:listener/app/name/0123456789abcdef/abcdef0123456789",                       //lintignore:AWSAT003,AWSAT005
		},
		{
			name:     "listener arn",
			arn:      "arn:aws:elasticloadbalancing:us-east-1:012345678912:listener/app/name/0123456789abcdef/abcdef0123456789", //lintignore:AWSAT003,AWSAT005
			expected: "",
		},
		{
			name:     "some other arn",
			arn:      "arn:aws:elasticloadbalancing:us-east-1:123456:targetgroup/my-targets/73e2d6bc24d8a067", //lintignore:AWSAT003,AWSAT005
			expected: "",
		},
		{
			name:     "not an arn",
			arn:      "blah blah blah",
			expected: "",
		},
		{
			name:     "empty arn",
			arn:      "",
			expected: "",
		},
	}

	for _, tc := range cases {
		actual := tfelbv2.ListenerARNFromRuleARN(tc.arn)
		if actual != tc.expected {
			t.Fatalf("incorrect arn returned: %q\nExpected: %s\n     Got: %s", tc.name, tc.expected, actual)
		}
	}
}

func TestAccELBV2ListenerRule_basic(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-basic-%s", sdkacctest.RandString(13))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", sdkacctest.RandString(10))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"
	targetGroupResourceName := "aws_lb_target_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_basic(lbName, targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.type", "forward"),
					resource.TestCheckResourceAttrPair(resourceName, "action.0.target_group_arn", targetGroupResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_cognito.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":           "0",
						"http_header.#":           "0",
						"http_request_method.#":   "0",
						"path_pattern.#":          "1",
						"path_pattern.0.values.#": "1",
						"query_string.#":          "0",
						"source_ip.#":             "0",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.path_pattern.0.values.*", "/static/*"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_tags(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-basic-%s", sdkacctest.RandString(13))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", sdkacctest.RandString(10))

	resourceName := "aws_lb_listener_rule.static"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_tags1(lbName, targetGroupName, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccListenerRuleConfig_tags2(lbName, targetGroupName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccListenerRuleConfig_tags1(lbName, targetGroupName, "key2", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_forwardWeighted(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-weighted-%s", sdkacctest.RandString(13))
	targetGroupName1 := fmt.Sprintf("testtargetgroup-%s", sdkacctest.RandString(10))
	targetGroupName2 := fmt.Sprintf("testtargetgroup-%s", sdkacctest.RandString(10))

	resourceName := "aws_lb_listener_rule.weighted"
	frontEndListenerResourceName := "aws_lb_listener.front_end"
	targetGroup1ResourceName := "aws_lb_target_group.test1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_forwardWeighted(lbName, targetGroupName1, targetGroupName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.type", "forward"),
					resource.TestCheckResourceAttr(resourceName, "action.0.forward.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.forward.0.target_group.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.0.forward.0.stickiness.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "action.0.forward.0.stickiness.0.duration", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_cognito.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
				),
			},
			{
				Config: testAccListenerRuleConfig_changeForwardWeightedStickiness(lbName, targetGroupName1, targetGroupName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.type", "forward"),
					resource.TestCheckResourceAttr(resourceName, "action.0.forward.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.forward.0.target_group.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.0.forward.0.stickiness.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "action.0.forward.0.stickiness.0.duration", "3600"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_cognito.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
				),
			},
			{
				Config: testAccListenerRuleConfig_changeForwardWeightedToBasic(lbName, targetGroupName1, targetGroupName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.type", "forward"),
					resource.TestCheckResourceAttrPair(resourceName, "action.0.target_group_arn", targetGroup1ResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_cognito.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_backwardsCompatibility(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-basic-%s", sdkacctest.RandString(13))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", sdkacctest.RandString(10))

	resourceName := "aws_alb_listener_rule.static"
	frontEndListenerResourceName := "aws_alb_listener.front_end"
	targetGroupResourceName := "aws_alb_target_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_backwardsCompatibility(lbName, targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.type", "forward"),
					resource.TestCheckResourceAttrPair(resourceName, "action.0.target_group_arn", targetGroupResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_cognito.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":           "0",
						"http_header.#":           "0",
						"http_request_method.#":   "0",
						"path_pattern.#":          "1",
						"path_pattern.0.values.#": "1",
						"query_string.#":          "0",
						"source_ip.#":             "0",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.path_pattern.0.values.*", "/static/*"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_redirect(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-redirect-%s", sdkacctest.RandString(14))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_redirect(lbName, "null"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.type", "redirect"),
					resource.TestCheckResourceAttr(resourceName, "action.0.target_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.host", "#{host}"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.path", "/#{path}"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.query", "#{query}"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.status_code", "HTTP_301"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_cognito.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
				),
			},
			{
				Config: testAccListenerRuleConfig_redirect(lbName, "param1=value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.type", "redirect"),
					resource.TestCheckResourceAttr(resourceName, "action.0.target_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.host", "#{host}"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.path", "/#{path}"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.query", "param1=value1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.status_code", "HTTP_301"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_cognito.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
				),
			},
			{
				Config: testAccListenerRuleConfig_redirect(lbName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.type", "redirect"),
					resource.TestCheckResourceAttr(resourceName, "action.0.target_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.host", "#{host}"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.path", "/#{path}"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.query", ""),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.0.status_code", "HTTP_301"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_cognito.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_fixedResponse(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-fixedresponse-%s", sdkacctest.RandString(9))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_fixedResponse(lbName, "Fixed response content"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.type", "fixed-response"),
					resource.TestCheckResourceAttr(resourceName, "action.0.target_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.0.content_type", "text/plain"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.0.message_body", "Fixed response content"),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.0.status_code", "200"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_cognito.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
				),
			},
		},
	})
}

// Updating Action breaks Condition change logic GH-11323 and GH-11362
func TestAccELBV2ListenerRule_updateFixedResponse(t *testing.T) {
	var rule elbv2.Rule
	lbName := fmt.Sprintf("testrule-basic-%s", sdkacctest.RandString(13))

	resourceName := "aws_lb_listener_rule.static"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_fixedResponse(lbName, "Fixed Response 1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.0.message_body", "Fixed Response 1"),
				),
			},
			{
				Config: testAccListenerRuleConfig_fixedResponse(lbName, "Fixed Response 2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "action.0.fixed_response.0.message_body", "Fixed Response 2"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_updateRulePriority(t *testing.T) {
	var before, after elbv2.Rule
	lbName := fmt.Sprintf("testrule-basic-%s", sdkacctest.RandString(13))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", sdkacctest.RandString(10))

	resourceName := "aws_lb_listener_rule.static"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_basic(lbName, targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &before),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
				),
			},
			{
				Config: testAccListenerRuleConfig_updatePriority(lbName, targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &after),
					testAccCheckListenerRuleNotRecreated(t, &before, &after),
					resource.TestCheckResourceAttr(resourceName, "priority", "101"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_changeListenerRuleARNForcesNew(t *testing.T) {
	var before, after elbv2.Rule
	lbName := fmt.Sprintf("testrule-basic-%s", sdkacctest.RandString(13))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", sdkacctest.RandString(10))

	resourceName := "aws_lb_listener_rule.static"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_basic(lbName, targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &before),
				),
			},
			{
				Config: testAccListenerRuleConfig_changeARN(lbName, targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &after),
					testAccCheckListenerRuleRecreated(t, &before, &after),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_priority(t *testing.T) {
	var rule elbv2.Rule
	lbName := fmt.Sprintf("testrule-basic-%s", sdkacctest.RandString(13))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", sdkacctest.RandString(10))

	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_priorityFirst(lbName, targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists("aws_lb_listener_rule.first", &rule),
					resource.TestCheckResourceAttr("aws_lb_listener_rule.first", "priority", "1"),
				),
			},
			{
				Config: testAccListenerRuleConfig_priorityLast(lbName, targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists("aws_lb_listener_rule.last", &rule),
					resource.TestCheckResourceAttr("aws_lb_listener_rule.last", "priority", "4"),
				),
			},
			{
				Config: testAccListenerRuleConfig_priorityStatic(lbName, targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists("aws_lb_listener_rule.last", &rule),
					resource.TestCheckResourceAttr("aws_lb_listener_rule.last", "priority", "7"),
				),
			},
			{
				Config: testAccListenerRuleConfig_priorityLast(lbName, targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists("aws_lb_listener_rule.last", &rule),
					resource.TestCheckResourceAttr("aws_lb_listener_rule.last", "priority", "7"),
				),
			},
			{
				Config: testAccListenerRuleConfig_priorityParallelism(lbName, targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists("aws_lb_listener_rule.last", &rule),
					resource.TestCheckResourceAttrSet("aws_lb_listener_rule.parallelism.0", "priority"),
					resource.TestCheckResourceAttrSet("aws_lb_listener_rule.parallelism.1", "priority"),
					resource.TestCheckResourceAttrSet("aws_lb_listener_rule.parallelism.2", "priority"),
					resource.TestCheckResourceAttrSet("aws_lb_listener_rule.parallelism.3", "priority"),
					resource.TestCheckResourceAttrSet("aws_lb_listener_rule.parallelism.4", "priority"),
					resource.TestCheckResourceAttrSet("aws_lb_listener_rule.parallelism.5", "priority"),
					resource.TestCheckResourceAttrSet("aws_lb_listener_rule.parallelism.6", "priority"),
					resource.TestCheckResourceAttrSet("aws_lb_listener_rule.parallelism.7", "priority"),
					resource.TestCheckResourceAttrSet("aws_lb_listener_rule.parallelism.8", "priority"),
					resource.TestCheckResourceAttrSet("aws_lb_listener_rule.parallelism.9", "priority"),
				),
			},
			{
				Config: testAccListenerRuleConfig_priority50000(lbName, targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRuleExists("aws_lb_listener_rule.priority50000", &rule),
					resource.TestCheckResourceAttr("aws_lb_listener_rule.priority50000", "priority", "50000"),
				),
			},
			{
				Config:      testAccListenerRuleConfig_priority50001(lbName, targetGroupName),
				ExpectError: regexp.MustCompile(`Error creating LB Listener Rule: ValidationError`),
			},
			{
				Config:      testAccListenerRuleConfig_priorityInUse(lbName, targetGroupName),
				ExpectError: regexp.MustCompile(`Error creating LB Listener Rule: PriorityInUse`),
			},
		},
	})
}

func TestAccELBV2ListenerRule_cognito(t *testing.T) {
	var conf elbv2.Rule
	key := acctest.TLSRSAPrivateKeyPEM(2048)
	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(key, "example.com")
	lbName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resourceName := "aws_lb_listener_rule.cognito"
	frontEndListenerResourceName := "aws_lb_listener.front_end"
	targetGroupResourceName := "aws_lb_target_group.test"
	cognitoPoolResourceName := "aws_cognito_user_pool.test"
	cognitoPoolClientResourceName := "aws_cognito_user_pool_client.test"
	cognitoPoolDomainResourceName := "aws_cognito_user_pool_domain.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_cognito(lbName, key, certificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.type", "authenticate-cognito"),
					resource.TestCheckResourceAttrPair(resourceName, "action.0.authenticate_cognito.0.user_pool_arn", cognitoPoolResourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "action.0.authenticate_cognito.0.user_pool_client_id", cognitoPoolClientResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "action.0.authenticate_cognito.0.user_pool_domain", cognitoPoolDomainResourceName, "domain"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_cognito.0.authentication_request_extra_params.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_cognito.0.authentication_request_extra_params.param", "test"),
					resource.TestCheckResourceAttr(resourceName, "action.1.order", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.1.type", "forward"),
					resource.TestCheckResourceAttrPair(resourceName, "action.1.target_group_arn", targetGroupResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_oidc(t *testing.T) {
	var conf elbv2.Rule
	key := acctest.TLSRSAPrivateKeyPEM(2048)
	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(key, "example.com")
	lbName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resourceName := "aws_lb_listener_rule.oidc"
	frontEndListenerResourceName := "aws_lb_listener.front_end"
	targetGroupResourceName := "aws_lb_target_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_oidc(lbName, key, certificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.type", "authenticate-oidc"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.0.authorization_endpoint", "https://example.com/authorization_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.0.client_id", "s6BhdRkqt3"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.0.client_secret", "7Fjfp0ZBr1KtDRbnfVdmIw"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.0.issuer", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.0.token_endpoint", "https://example.com/token_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.0.user_info_endpoint", "https://example.com/user_info_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.0.authentication_request_extra_params.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0.authenticate_oidc.0.authentication_request_extra_params.param", "test"),
					resource.TestCheckResourceAttr(resourceName, "action.1.order", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.1.type", "forward"),
					resource.TestCheckResourceAttrPair(resourceName, "action.1.target_group_arn", targetGroupResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_Action_order(t *testing.T) {
	var rule elbv2.Rule
	key := acctest.TLSRSAPrivateKeyPEM(2048)
	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(key, "example.com")
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_lb_listener_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_actionOrder(rName, key, certificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.1.order", "2"),
				),
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/6171
func TestAccELBV2ListenerRule_ActionOrder_recreates(t *testing.T) {
	var rule elbv2.Rule
	key := acctest.TLSRSAPrivateKeyPEM(2048)
	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(key, "example.com")
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_lb_listener_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_actionOrder(rName, key, certificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.1.order", "2"),
					testAccCheckListenerRuleActionOrderDisappears(&rule, 1),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccELBV2ListenerRule_conditionAttributesCount(t *testing.T) {
	err_many := regexp.MustCompile("Only one of host_header, http_header, http_request_method, path_pattern, query_string or source_ip can be set in a condition block")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccListenerRuleConfig_conditionAttributesCountHTTPHeader(),
				ExpectError: err_many,
			},
			{
				Config:      testAccListenerRuleConfig_conditionAttributesCountHTTPRequestMethod(),
				ExpectError: err_many,
			},
			{
				Config:      testAccListenerRuleConfig_conditionAttributesCountPathPattern(),
				ExpectError: err_many,
			},
			{
				Config:      testAccListenerRuleConfig_conditionAttributesCountQueryString(),
				ExpectError: err_many,
			},
			{
				Config:      testAccListenerRuleConfig_conditionAttributesCountSourceIP(),
				ExpectError: err_many,
			},
		},
	})
}

func TestAccELBV2ListenerRule_conditionHostHeader(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-hostHeader-%s", sdkacctest.RandString(12))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_conditionHostHeader(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":          "1",
						"host_header.0.values.#": "2",
						"http_header.#":          "0",
						"http_request_method.#":  "0",
						"path_pattern.#":         "0",
						"query_string.#":         "0",
						"source_ip.#":            "0",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.host_header.0.values.*", "example.com"),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.host_header.0.values.*", "www.example.com"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_conditionHTTPHeader(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-httpHeader-%s", sdkacctest.RandString(12))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_conditionHTTPHeader(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":                  "0",
						"http_header.#":                  "1",
						"http_header.0.http_header_name": "X-Forwarded-For",
						"http_header.0.values.#":         "2",
						"http_request_method.#":          "0",
						"path_pattern.#":                 "0",
						"query_string.#":                 "0",
						"source_ip.#":                    "0",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.http_header.0.values.*", "10.0.0.*"),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.http_header.0.values.*", "192.168.1.*"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":                  "0",
						"http_header.#":                  "1",
						"http_header.0.http_header_name": "Zz9~|_^.-+*'&%$#!0aA",
						"http_header.0.values.#":         "1",
						"http_request_method.#":          "0",
						"path_pattern.#":                 "0",
						"query_string.#":                 "0",
						"source_ip.#":                    "0",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.http_header.0.values.*", "RFC7230 Validity"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_ConditionHTTPHeader_invalid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccListenerRuleConfig_conditionHTTPHeaderInvalid(),
				ExpectError: regexp.MustCompile(`expected value of condition.0.http_header.0.http_header_name to match regular expression`),
			},
		},
	})
}

func TestAccELBV2ListenerRule_conditionHTTPRequestMethod(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-httpRequest-%s", sdkacctest.RandString(11))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_conditionHTTPRequestMethod(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":                  "0",
						"http_header.#":                  "0",
						"http_request_method.#":          "1",
						"http_request_method.0.values.#": "2",
						"path_pattern.#":                 "0",
						"query_string.#":                 "0",
						"source_ip.#":                    "0",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.http_request_method.0.values.*", "GET"),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.http_request_method.0.values.*", "POST"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_conditionPathPattern(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-pathPattern-%s", sdkacctest.RandString(11))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_conditionPathPattern(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":           "0",
						"http_header.#":           "0",
						"http_request_method.#":   "0",
						"path_pattern.#":          "1",
						"path_pattern.0.values.#": "2",
						"query_string.#":          "0",
						"source_ip.#":             "0",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.path_pattern.0.values.*", "/cgi-bin/*"),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.path_pattern.0.values.*", "/public/*"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_conditionQueryString(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-queryString-%s", sdkacctest.RandString(11))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_conditionQueryString(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":         "0",
						"http_header.#":         "0",
						"http_request_method.#": "0",
						"path_pattern.#":        "0",
						"query_string.#":        "2",
						"source_ip.#":           "0",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*.query_string.*", map[string]string{
						"key":   "",
						"value": "surprise",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*.query_string.*", map[string]string{
						"key":   "",
						"value": "blank",
					}),
					// TODO: TypeSet check helpers cannot make distinction between the 2 set items
					// because we had to write a new check for the "downstream" nested set
					// a distinguishing attribute on the outer set would be solve this.
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":         "0",
						"http_header.#":         "0",
						"http_request_method.#": "0",
						"path_pattern.#":        "0",
						"query_string.#":        "2",
						"source_ip.#":           "0",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*.query_string.*", map[string]string{
						"key":   "foo",
						"value": "baz",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*.query_string.*", map[string]string{
						"key":   "foo",
						"value": "bar",
					}),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_conditionSourceIP(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-sourceIp-%s", sdkacctest.RandString(14))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_conditionSourceIP(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":         "0",
						"http_header.#":         "0",
						"http_request_method.#": "0",
						"path_pattern.#":        "0",
						"query_string.#":        "0",
						"source_ip.#":           "1",
						"source_ip.0.values.#":  "2",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.source_ip.0.values.*", "dead:cafe::/64"),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.source_ip.0.values.*", "192.168.0.0/16"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_conditionUpdateMixed(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-mixed-%s", sdkacctest.RandString(17))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"

	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_conditionMixed(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"path_pattern.#":          "1",
						"path_pattern.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.path_pattern.0.values.*", "/public/*"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"source_ip.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.source_ip.0.values.*", "192.168.0.0/16"),
				),
			},
			{
				Config: testAccListenerRuleConfig_conditionMixedUpdated(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"path_pattern.#":          "1",
						"path_pattern.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.path_pattern.0.values.*", "/public/*"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"source_ip.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.source_ip.0.values.*", "dead:cafe::/64"),
				),
			},
			{
				Config: testAccListenerRuleConfig_conditionMixedUpdated2(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"path_pattern.#":          "1",
						"path_pattern.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.path_pattern.0.values.*", "/cgi-bin/*"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"source_ip.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.source_ip.0.values.*", "dead:cafe::/64"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_conditionMultiple(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-condMulti-%s", sdkacctest.RandString(13))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_conditionMultiple(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "5"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":                  "0",
						"http_header.#":                  "1",
						"http_header.0.http_header_name": "X-Forwarded-For",
						"http_header.0.values.#":         "1",
						"http_request_method.#":          "0",
						"path_pattern.#":                 "0",
						"query_string.#":                 "0",
						"source_ip.#":                    "0",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.http_header.0.values.*", "192.168.1.*"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":         "0",
						"http_header.#":         "0",
						"http_request_method.#": "0",
						"path_pattern.#":        "0",
						"query_string.#":        "0",
						"source_ip.#":           "1",
						"source_ip.0.values.#":  "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.source_ip.0.values.*", "192.168.0.0/16"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":                  "0",
						"http_header.#":                  "0",
						"http_request_method.#":          "1",
						"http_request_method.0.values.#": "1",
						"path_pattern.#":                 "0",
						"query_string.#":                 "0",
						"source_ip.#":                    "0",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.http_request_method.0.values.*", "GET"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":           "0",
						"http_header.#":           "0",
						"http_request_method.#":   "0",
						"path_pattern.#":          "1",
						"path_pattern.0.values.#": "1",
						"query_string.#":          "0",
						"source_ip.#":             "0",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.path_pattern.0.values.*", "/public/*"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":          "1",
						"host_header.0.values.#": "1",
						"http_header.#":          "0",
						"http_request_method.#":  "0",
						"path_pattern.#":         "0",
						"query_string.#":         "0",
						"source_ip.#":            "0",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.host_header.0.values.*", "example.com"),
				),
			},
		},
	})
}

func TestAccELBV2ListenerRule_conditionUpdateMultiple(t *testing.T) {
	var conf elbv2.Rule
	lbName := fmt.Sprintf("testrule-condMulti-%s", sdkacctest.RandString(13))

	resourceName := "aws_lb_listener_rule.static"
	frontEndListenerResourceName := "aws_lb_listener.front_end"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerRuleConfig_conditionMultiple(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "5"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"http_header.#":                  "1",
						"http_header.0.http_header_name": "X-Forwarded-For",
						"http_header.0.values.#":         "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.http_header.0.values.*", "192.168.1.*"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"source_ip.#":          "1",
						"source_ip.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.source_ip.0.values.*", "192.168.0.0/16"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"http_request_method.#":          "1",
						"http_request_method.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.http_request_method.0.values.*", "GET"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"path_pattern.#":          "1",
						"path_pattern.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.path_pattern.0.values.*", "/public/*"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":          "1",
						"host_header.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.host_header.0.values.*", "example.com"),
				),
			},
			{
				Config: testAccListenerRuleConfig_conditionMultipleUpdated(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckListenerRuleExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile(fmt.Sprintf(`listener-rule/app/%s/.+$`, lbName))),
					resource.TestCheckResourceAttrPair(resourceName, "listener_arn", frontEndListenerResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "5"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"http_header.#":                  "1",
						"http_header.0.http_header_name": "X-Forwarded-For",
						"http_header.0.values.#":         "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.http_header.0.values.*", "192.168.2.*"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"source_ip.#":          "1",
						"source_ip.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.source_ip.0.values.*", "192.168.0.0/24"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"http_request_method.#":          "1",
						"http_request_method.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.http_request_method.0.values.*", "POST"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"path_pattern.#":          "1",
						"path_pattern.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.path_pattern.0.values.*", "/public/2/*"),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "condition.*", map[string]string{
						"host_header.#":          "1",
						"host_header.0.values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "condition.*.host_header.0.values.*", "example.com"),
				),
			},
		},
	})
}

func testAccCheckListenerRuleActionOrderDisappears(rule *elbv2.Rule, actionOrderToDelete int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var newActions []*elbv2.Action

		for i, action := range rule.Actions {
			if int(aws.Int64Value(action.Order)) == actionOrderToDelete {
				newActions = append(rule.Actions[:i], rule.Actions[i+1:]...)
				break
			}
		}

		if len(newActions) == 0 {
			return fmt.Errorf("Unable to find action order %d from actions: %#v", actionOrderToDelete, rule.Actions)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ELBV2Conn

		input := &elbv2.ModifyRuleInput{
			Actions: newActions,
			RuleArn: rule.RuleArn,
		}

		_, err := conn.ModifyRule(input)

		return err
	}
}

func testAccCheckListenerRuleNotRecreated(t *testing.T, before, after *elbv2.Rule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.StringValue(before.RuleArn), aws.StringValue(after.RuleArn); before != after {
			t.Fatalf("ELBv2 Listener Rule (%s) was recreated: %s", before, after)
		}

		return nil
	}
}

func testAccCheckListenerRuleRecreated(t *testing.T, before, after *elbv2.Rule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.StringValue(before.RuleArn), aws.StringValue(after.RuleArn); before == after {
			t.Fatalf("ELBv2 Listener Rule (%s) was not recreated", before)
		}

		return nil
	}
}

func testAccCheckListenerRuleExists(n string, res *elbv2.Rule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Listener Rule ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ELBV2Conn

		describe, err := conn.DescribeRules(&elbv2.DescribeRulesInput{
			RuleArns: []*string{aws.String(rs.Primary.ID)},
		})

		if err != nil {
			return err
		}

		if len(describe.Rules) != 1 ||
			*describe.Rules[0].RuleArn != rs.Primary.ID {
			return errors.New("Listener Rule not found")
		}

		*res = *describe.Rules[0]
		return nil
	}
}

func testAccCheckListenerRuleDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ELBV2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_lb_listener_rule" && rs.Type != "aws_alb_listener_rule" {
			continue
		}

		describe, err := conn.DescribeRules(&elbv2.DescribeRulesInput{
			RuleArns: []*string{aws.String(rs.Primary.ID)},
		})

		if err == nil {
			if len(describe.Rules) != 0 &&
				*describe.Rules[0].RuleArn == rs.Primary.ID {
				return fmt.Errorf("Listener Rule %q still exists", rs.Primary.ID)
			}
		}

		// Verify the error
		if tfawserr.ErrCodeEquals(err, elbv2.ErrCodeRuleNotFoundException) {
			return nil
		} else {
			return fmt.Errorf("Unexpected error checking LB Listener Rule destroyed: %s", err)
		}
	}

	return nil
}

func testAccListenerRuleConfig_basic(lbName, targetGroupName string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "static" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = aws_lb_target_group.test.id
    type             = "forward"
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-basic"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-basic-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, lbName, targetGroupName)
}

func testAccListenerRuleConfig_forwardWeighted(lbName, targetGroupName1 string, targetGroupName2 string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "weighted" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 100

  action {
    type = "forward"

    forward {
      target_group {
        arn    = aws_lb_target_group.test1.arn
        weight = 1
      }

      target_group {
        arn    = aws_lb_target_group.test2.arn
        weight = 1
      }
    }
  }

  condition {
    path_pattern {
      values = ["/weighted/*"]
    }
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = aws_lb_target_group.test1.arn
    type             = "forward"
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test1" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

resource "aws_lb_target_group" "test2" {
  name     = %[3]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-basic"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-basic-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, lbName, targetGroupName1, targetGroupName2)
}

func testAccListenerRuleConfig_changeForwardWeightedStickiness(lbName, targetGroupName1 string, targetGroupName2 string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "weighted" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 100

  action {
    type = "forward"

    forward {
      target_group {
        arn    = aws_lb_target_group.test1.arn
        weight = 1
      }

      target_group {
        arn    = aws_lb_target_group.test2.arn
        weight = 1
      }

      stickiness {
        enabled  = true
        duration = 3600
      }
    }
  }

  condition {
    path_pattern {
      values = ["/weighted/*"]
    }
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = aws_lb_target_group.test1.arn
    type             = "forward"
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test1" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

resource "aws_lb_target_group" "test2" {
  name     = %[3]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-basic"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-basic-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, lbName, targetGroupName1, targetGroupName2)
}

func testAccListenerRuleConfig_changeForwardWeightedToBasic(lbName, targetGroupName1 string, targetGroupName2 string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "weighted" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test1.arn
  }

  condition {
    path_pattern {
      values = ["/weighted/*"]
    }
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = aws_lb_target_group.test1.arn
    type             = "forward"
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test1" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

resource "aws_lb_target_group" "test2" {
  name     = %[3]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-basic"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-basic-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, lbName, targetGroupName1, targetGroupName2)
}

func testAccListenerRuleConfig_backwardsCompatibility(lbName, targetGroupName string) string {
	return fmt.Sprintf(`
resource "aws_alb_listener_rule" "static" {
  listener_arn = aws_alb_listener.front_end.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }
}

resource "aws_alb_listener" "front_end" {
  load_balancer_arn = aws_alb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = aws_alb_target_group.test.id
    type             = "forward"
  }
}

resource "aws_alb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_alb_target_group" "test" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-bc"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-bc-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, lbName, targetGroupName)
}

func testAccListenerRuleConfig_redirect(lbName, query string) string {
	if query != "null" {
		query = strconv.Quote(query)
	}

	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "static" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 100

  action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      query       = %[2]s
      status_code = "HTTP_301"
    }
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-redirect"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-redirect-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, lbName, query)
}

func testAccListenerRuleConfig_fixedResponse(lbName, response string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "static" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 100

  action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = %[1]q
      status_code  = "200"
    }
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "Fixed response content"
      status_code  = "200"
    }
  }
}

resource "aws_lb" "alb_test" {
  name            = %[2]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[2]q
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-fixedresponse"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-fixedresponse-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[2]q
  }
}
`, response, lbName)
}

func testAccListenerRuleConfig_updatePriority(lbName, targetGroupName string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "static" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 101

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = aws_lb_target_group.test.id
    type             = "forward"
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-update-rule-priority"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-update-rule-priority-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, lbName, targetGroupName)
}

func testAccListenerRuleConfig_changeARN(lbName, targetGroupName string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "static" {
  listener_arn = aws_lb_listener.front_end_ruleupdate.arn
  priority     = 101

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = aws_lb_target_group.test.id
    type             = "forward"
  }
}

resource "aws_lb_listener" "front_end_ruleupdate" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "8080"

  default_action {
    target_group_arn = aws_lb_target_group.test.id
    type             = "forward"
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-change-rule-arn"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-change-rule-arn-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, lbName, targetGroupName)
}

func testAccListenerRuleConfig_priorityBase(lbName, targetGroupName string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = aws_lb_target_group.test.id
    type             = "forward"
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-priority"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-priority-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, lbName, targetGroupName)
}

func testAccListenerRuleConfig_priorityFirst(lbName, targetGroupName string) string {
	return acctest.ConfigCompose(testAccListenerRuleConfig_priorityBase(lbName, targetGroupName), `
resource "aws_lb_listener_rule" "first" {
  listener_arn = aws_lb_listener.front_end.arn

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/first/*"]
    }
  }
}

resource "aws_lb_listener_rule" "third" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 3

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/third/*"]
    }
  }

  depends_on = [aws_lb_listener_rule.first]
}
`)
}

func testAccListenerRuleConfig_priorityLast(lbName, targetGroupName string) string {
	return acctest.ConfigCompose(testAccListenerRuleConfig_priorityFirst(lbName, targetGroupName), `
resource "aws_lb_listener_rule" "last" {
  listener_arn = aws_lb_listener.front_end.arn

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/last/*"]
    }
  }
}
`)
}

func testAccListenerRuleConfig_priorityStatic(lbName, targetGroupName string) string {
	return acctest.ConfigCompose(testAccListenerRuleConfig_priorityFirst(lbName, targetGroupName), `
resource "aws_lb_listener_rule" "last" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 7

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/last/*"]
    }
  }
}
`)
}

func testAccListenerRuleConfig_priorityParallelism(lbName, targetGroupName string) string {
	return acctest.ConfigCompose(testAccListenerRuleConfig_priorityStatic(lbName, targetGroupName), `
resource "aws_lb_listener_rule" "parallelism" {
  count = 10

  listener_arn = aws_lb_listener.front_end.arn

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/${count.index}/*"]
    }
  }
}
`)
}

func testAccListenerRuleConfig_priority50000(lbName, targetGroupName string) string {
	return acctest.ConfigCompose(testAccListenerRuleConfig_priorityBase(lbName, targetGroupName), `
resource "aws_lb_listener_rule" "priority50000" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 50000

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/50000/*"]
    }
  }
}
`)
}

// priority out of range (1, 50000)
func testAccListenerRuleConfig_priority50001(lbName, targetGroupName string) string {
	return acctest.ConfigCompose(testAccListenerRuleConfig_priority50000(lbName, targetGroupName), `
resource "aws_lb_listener_rule" "priority50001" {
  listener_arn = aws_lb_listener.front_end.arn

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/50001/*"]
    }
  }
}
`)
}

func testAccListenerRuleConfig_priorityInUse(lbName, targetGroupName string) string {
	return acctest.ConfigCompose(testAccListenerRuleConfig_priority50000(lbName, targetGroupName), `
resource "aws_lb_listener_rule" "priority50000_in_use" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 50000

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/50000_in_use/*"]
    }
  }
}
`)
}

func testAccListenerRuleConfig_cognito(rName, key, certificate string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "cognito" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 100

  action {
    type = "authenticate-cognito"

    authenticate_cognito {
      user_pool_arn       = aws_cognito_user_pool.test.arn
      user_pool_client_id = aws_cognito_user_pool_client.test.id
      user_pool_domain    = aws_cognito_user_pool_domain.test.domain

      authentication_request_extra_params = {
        param = "test"
      }
    }
  }

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }
}

resource "aws_iam_server_certificate" "test" {
  name             = %[1]q
  certificate_body = %[2]q
  private_key      = %[3]q
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTPS"
  port              = "443"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_iam_server_certificate.test.arn

  default_action {
    target_group_arn = aws_lb_target_group.test.id
    type             = "forward"
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-cognito"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-cognito-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_cognito_user_pool" "test" {
  name = "%[1]s-pool"
}

resource "aws_cognito_user_pool_client" "test" {
  name                                 = "%[1]s-pool-client"
  user_pool_id                         = aws_cognito_user_pool.test.id
  generate_secret                      = true
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code", "implicit"]
  allowed_oauth_scopes                 = ["phone", "email", "openid", "profile", "aws.cognito.signin.user.admin"]
  callback_urls                        = ["https://www.example.com/callback", "https://www.example.com/redirect"]
  default_redirect_uri                 = "https://www.example.com/redirect"
  logout_urls                          = ["https://www.example.com/login"]
}

resource "aws_cognito_user_pool_domain" "test" {
  domain       = "%[1]s-pool-domain"
  user_pool_id = aws_cognito_user_pool.test.id
}
`, rName, acctest.TLSPEMEscapeNewlines(certificate), acctest.TLSPEMEscapeNewlines(key))
}

func testAccListenerRuleConfig_oidc(rName, key, certificate string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "oidc" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 100

  action {
    type = "authenticate-oidc"

    authenticate_oidc {
      authorization_endpoint = "https://example.com/authorization_endpoint"
      client_id              = "s6BhdRkqt3"
      client_secret          = "7Fjfp0ZBr1KtDRbnfVdmIw"
      issuer                 = "https://example.com"
      token_endpoint         = "https://example.com/token_endpoint"
      user_info_endpoint     = "https://example.com/user_info_endpoint"

      authentication_request_extra_params = {
        param = "test"
      }
    }
  }

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }
}

resource "aws_iam_server_certificate" "test" {
  name             = %[1]q
  certificate_body = %[2]q
  private_key      = %[3]q
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTPS"
  port              = "443"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_iam_server_certificate.test.arn

  default_action {
    target_group_arn = aws_lb_target_group.test.id
    type             = "forward"
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-cognito"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-cognito-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName, acctest.TLSPEMEscapeNewlines(certificate), acctest.TLSPEMEscapeNewlines(key))
}

func testAccListenerRuleConfig_actionOrder(rName, key, certificate string) string {
	return fmt.Sprintf(`
variable "rName" {
  default = %[1]q
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_lb_listener_rule" "test" {
  listener_arn = aws_lb_listener.test.arn

  action {
    order = 1
    type  = "authenticate-oidc"

    authenticate_oidc {
      authorization_endpoint = "https://example.com/authorization_endpoint"
      client_id              = "s6BhdRkqt3"
      client_secret          = "7Fjfp0ZBr1KtDRbnfVdmIw"
      issuer                 = "https://example.com"
      token_endpoint         = "https://example.com/token_endpoint"
      user_info_endpoint     = "https://example.com/user_info_endpoint"

      authentication_request_extra_params = {
        param = "test"
      }
    }
  }

  action {
    order            = 2
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }
}

resource "aws_iam_server_certificate" "test" {
  certificate_body = %[2]q
  name             = var.rName
  private_key      = %[3]q
}

resource "aws_lb_listener" "test" {
  load_balancer_arn = aws_lb.test.id
  protocol          = "HTTPS"
  port              = "443"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_iam_server_certificate.test.arn

  default_action {
    target_group_arn = aws_lb_target_group.test.id
    type             = "forward"
  }
}

resource "aws_lb" "test" {
  internal        = true
  name            = var.rName
  security_groups = [aws_security_group.test.id]
  subnets         = aws_subnet.test[*].id
}

resource "aws_lb_target_group" "test" {
  name     = var.rName
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = var.rName
  }
}

resource "aws_subnet" "test" {
  count = 2

  availability_zone       = data.aws_availability_zones.available.names[count.index]
  cidr_block              = "10.0.${count.index}.0/24"
  map_public_ip_on_launch = true
  vpc_id                  = aws_vpc.test.id

  tags = {
    Name = var.rName
  }
}

resource "aws_security_group" "test" {
  name   = var.rName
  vpc_id = aws_vpc.test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = var.rName
  }
}
`, rName, acctest.TLSPEMEscapeNewlines(certificate), acctest.TLSPEMEscapeNewlines(key))
}

func testAccListenerRuleConfig_condition_error(condition string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

data "aws_region" "current" {}

resource "aws_lb_listener_rule" "error" {
  listener_arn = "arn:${data.aws_partition.current.partition}:elasticloadbalancing:${data.aws_region.current.name}:111111111111:listener/app/example/1234567890abcdef/1234567890abcdef"
  priority     = 100

  action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "Static"
      status_code  = 200
    }
  }

  %s
}
`, condition)
}

func testAccListenerRuleConfig_conditionAttributesCountHTTPHeader() string {
	return testAccListenerRuleConfig_condition_error(`
condition {
  host_header {
    values = ["example.com"]
  }

  http_header {
    http_header_name = "X-Clacks-Overhead"
    values           = ["GNU Terry Pratchett"]
  }
}
`)
}

func testAccListenerRuleConfig_conditionAttributesCountHTTPRequestMethod() string {
	return testAccListenerRuleConfig_condition_error(`
condition {
  host_header {
    values = ["example.com"]
  }

  http_request_method {
    values = ["POST"]
  }
}
`)
}

func testAccListenerRuleConfig_conditionAttributesCountPathPattern() string {
	return testAccListenerRuleConfig_condition_error(`
condition {
  host_header {
    values = ["example.com"]
  }

  path_pattern {
    values = ["/"]
  }
}
`)
}

func testAccListenerRuleConfig_conditionAttributesCountQueryString() string {
	return testAccListenerRuleConfig_condition_error(`
condition {
  host_header {
    values = ["example.com"]
  }

  query_string {
    key   = "foo"
    value = "bar"
  }
}
`)
}

func testAccListenerRuleConfig_conditionAttributesCountSourceIP() string {
	return testAccListenerRuleConfig_condition_error(`
condition {
  host_header {
    values = ["example.com"]
  }

  source_ip {
    values = ["192.168.0.0/16"]
  }
}
`)
}

func testAccListenerRuleConfig_conditionBase(condition, lbName string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "static" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 100

  action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "Static"
      status_code  = 200
    }
  }

  %[1]s
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "Not Found"
      status_code  = 404
    }
  }
}

resource "aws_lb" "alb_test" {
  name            = %[2]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[2]q
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[2]q
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "%[2]s-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[2]q
  }
}
`, condition, lbName)
}

func testAccListenerRuleConfig_conditionHostHeader(lbName string) string {
	return testAccListenerRuleConfig_conditionBase(`
condition {
  host_header {
    values = ["example.com", "www.example.com"]
  }
}
`, lbName)
}

func testAccListenerRuleConfig_conditionHTTPHeader(lbName string) string {
	return testAccListenerRuleConfig_conditionBase(`
condition {
  http_header {
    http_header_name = "X-Forwarded-For"
    values           = ["192.168.1.*", "10.0.0.*"]
  }
}

condition {
  http_header {
    http_header_name = "Zz9~|_^.-+*'&%$#!0aA"
    values           = ["RFC7230 Validity"]
  }
}
`, lbName)
}

func testAccListenerRuleConfig_conditionHTTPHeaderInvalid() string {
	return `
data "aws_partition" "current" {}

data "aws_region" "current" {}

resource "aws_lb_listener_rule" "static" {
  listener_arn = "arn:${data.aws_partition.current.partition}:elasticloadbalancing:${data.aws_region.current.name}:111111111111:listener/app/test/xxxxxxxxxxxxxxxx/xxxxxxxxxxxxxxxx"
  priority     = 100

  action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "Static"
      status_code  = 200
    }
  }

  condition {
    http_header {
      http_header_name = "Invalid@"
      values           = ["RFC7230 Validity"]
    }
  }
}
`
}

func testAccListenerRuleConfig_conditionHTTPRequestMethod(lbName string) string {
	return testAccListenerRuleConfig_conditionBase(`
condition {
  http_request_method {
    values = ["GET", "POST"]
  }
}
`, lbName)
}

func testAccListenerRuleConfig_conditionPathPattern(lbName string) string {
	return testAccListenerRuleConfig_conditionBase(`
condition {
  path_pattern {
    values = ["/public/*", "/cgi-bin/*"]
  }
}
`, lbName)
}

func testAccListenerRuleConfig_conditionQueryString(lbName string) string {
	return testAccListenerRuleConfig_conditionBase(`
condition {
  query_string {
    value = "surprise"
  }

  query_string {
    key   = ""
    value = "blank"
  }
}

condition {
  query_string {
    key   = "foo"
    value = "bar"
  }

  query_string {
    key   = "foo"
    value = "baz"
  }
}
`, lbName)
}

func testAccListenerRuleConfig_conditionSourceIP(lbName string) string {
	return testAccListenerRuleConfig_conditionBase(`
condition {
  source_ip {
    values = [
      "192.168.0.0/16",
      "dead:cafe::/64",
    ]
  }
}
`, lbName)
}

func testAccListenerRuleConfig_conditionMixed(lbName string) string {
	return testAccListenerRuleConfig_conditionBase(`
condition {
  path_pattern {
    values = ["/public/*"]
  }
}

condition {
  source_ip {
    values = [
      "192.168.0.0/16",
    ]
  }
}
`, lbName)
}

// Update new style condition without modifying deprecated. Issue GH-11323
func testAccListenerRuleConfig_conditionMixedUpdated(lbName string) string {
	return testAccListenerRuleConfig_conditionBase(`
condition {
  path_pattern {
    values = ["/public/*"]
  }
}

condition {
  source_ip {
    values = [
      "dead:cafe::/64",
    ]
  }
}
`, lbName)
}

// Then update deprecated syntax without touching new. Issue GH-11362
func testAccListenerRuleConfig_conditionMixedUpdated2(lbName string) string {
	return testAccListenerRuleConfig_conditionBase(`
condition {
  path_pattern {
    values = ["/cgi-bin/*"]
  }
}

condition {
  source_ip {
    values = [
      "dead:cafe::/64",
    ]
  }
}
`, lbName)
}

// Currently a maximum of 5 condition values per rule
func testAccListenerRuleConfig_conditionMultiple(lbName string) string {
	return testAccListenerRuleConfig_conditionBase(`
condition {
  host_header {
    values = ["example.com"]
  }
}

condition {
  http_header {
    http_header_name = "X-Forwarded-For"
    values           = ["192.168.1.*"]
  }
}

condition {
  http_request_method {
    values = ["GET"]
  }
}

condition {
  path_pattern {
    values = ["/public/*"]
  }
}

condition {
  source_ip {
    values = ["192.168.0.0/16"]
  }
}
`, lbName)
}

func testAccListenerRuleConfig_conditionMultipleUpdated(lbName string) string {
	return testAccListenerRuleConfig_conditionBase(`
condition {
  host_header {
    values = ["example.com"]
  }
}

condition {
  http_header {
    http_header_name = "X-Forwarded-For"
    values           = ["192.168.2.*"]
  }
}

condition {
  http_request_method {
    values = ["POST"]
  }
}

condition {
  path_pattern {
    values = ["/public/2/*"]
  }
}

condition {
  source_ip {
    values = ["192.168.0.0/24"]
  }
}
`, lbName)
}

func testAccListenerRuleConfig_tags1(lbName, targetGroupName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "static" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }

  tags = {
    %[3]q = %[4]q
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = aws_lb_target_group.test.id
    type             = "forward"
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-basic"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-basic-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, lbName, targetGroupName, tagKey1, tagValue1)
}

func testAccListenerRuleConfig_tags2(lbName, targetGroupName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_lb_listener_rule" "static" {
  listener_arn = aws_lb_listener.front_end.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test.arn
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }

  tags = {
    %[3]q = %[4]q
    %[5]q = %[6]q
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.alb_test.id
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = aws_lb_target_group.test.id
    type             = "forward"
  }
}

resource "aws_lb" "alb_test" {
  name            = %[1]q
  internal        = true
  security_groups = [aws_security_group.alb_test.id]
  subnets         = aws_subnet.alb_test[*].id

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.alb_test.id

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-lb-listener-rule-basic"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = aws_vpc.alb_test.id
  cidr_block              = element(var.subnets, count.index)
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "tf-acc-lb-listener-rule-basic-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = aws_vpc.alb_test.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, lbName, targetGroupName, tagKey1, tagValue1, tagKey2, tagValue2)
}
