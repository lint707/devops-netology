package ec2_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestProtocolStateFunc(t *testing.T) {
	cases := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    "tcp",
			expected: "tcp",
		},
		{
			input:    6,
			expected: "",
		},
		{
			input:    "17",
			expected: "udp",
		},
		{
			input:    "all",
			expected: "-1",
		},
		{
			input:    "-1",
			expected: "-1",
		},
		{
			input:    -1,
			expected: "",
		},
		{
			input:    "1",
			expected: "icmp",
		},
		{
			input:    "icmp",
			expected: "icmp",
		},
		{
			input:    1,
			expected: "",
		},
		{
			input:    "icmpv6",
			expected: "icmpv6",
		},
		{
			input:    "58",
			expected: "icmpv6",
		},
		{
			input:    58,
			expected: "",
		},
	}
	for _, c := range cases {
		result := tfec2.ProtocolStateFunc(c.input)
		if result != c.expected {
			t.Errorf("Error matching protocol, expected (%s), got (%s)", c.expected, result)
		}
	}
}

func TestProtocolForValue(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "tcp",
			expected: "tcp",
		},
		{
			input:    "6",
			expected: "tcp",
		},
		{
			input:    "udp",
			expected: "udp",
		},
		{
			input:    "17",
			expected: "udp",
		},
		{
			input:    "all",
			expected: "-1",
		},
		{
			input:    "-1",
			expected: "-1",
		},
		{
			input:    "tCp",
			expected: "tcp",
		},
		{
			input:    "6",
			expected: "tcp",
		},
		{
			input:    "UDp",
			expected: "udp",
		},
		{
			input:    "17",
			expected: "udp",
		},
		{
			input:    "ALL",
			expected: "-1",
		},
		{
			input:    "icMp",
			expected: "icmp",
		},
		{
			input:    "1",
			expected: "icmp",
		},
		{
			input:    "icMpv6",
			expected: "icmpv6",
		},
		{
			input:    "58",
			expected: "icmpv6",
		},
	}

	for _, c := range cases {
		result := tfec2.ProtocolForValue(c.input)
		if result != c.expected {
			t.Errorf("Error matching protocol, expected (%s), got (%s)", c.expected, result)
		}
	}
}

func calcSecurityGroupChecksum(rules []interface{}) int {
	sum := 0
	for _, rule := range rules {
		sum += tfec2.SecurityGroupRuleHash(rule)
	}
	return sum
}

func TestSecurityGroupExpandCollapseRules(t *testing.T) {
	expected_compact_list := []interface{}{
		map[string]interface{}{
			"protocol":    "tcp",
			"from_port":   int(443),
			"to_port":     int(443),
			"description": "block with description",
			"self":        true,
			"cidr_blocks": []interface{}{
				"10.0.0.1/32",
				"10.0.0.2/32",
				"10.0.0.3/32",
			},
		},
		map[string]interface{}{
			"protocol":    "tcp",
			"from_port":   int(443),
			"to_port":     int(443),
			"description": "block with another description",
			"self":        false,
			"cidr_blocks": []interface{}{
				"192.168.0.1/32",
				"192.168.0.2/32",
			},
		},
		map[string]interface{}{
			"protocol":    "-1",
			"from_port":   int(8000),
			"to_port":     int(8080),
			"description": "",
			"self":        false,
			"ipv6_cidr_blocks": []interface{}{
				"fd00::1/128",
				"fd00::2/128",
			},
			"security_groups": schema.NewSet(schema.HashString, []interface{}{
				"sg-11111",
				"sg-22222",
				"sg-33333",
			}),
		},
		map[string]interface{}{
			"protocol":    "udp",
			"from_port":   int(10000),
			"to_port":     int(10000),
			"description": "",
			"self":        false,
			"prefix_list_ids": []interface{}{
				"pl-111111",
				"pl-222222",
			},
		},
	}

	expected_expanded_list := []interface{}{
		map[string]interface{}{
			"protocol":    "tcp",
			"from_port":   int(443),
			"to_port":     int(443),
			"description": "block with description",
			"self":        true,
		},
		map[string]interface{}{
			"protocol":    "tcp",
			"from_port":   int(443),
			"to_port":     int(443),
			"description": "block with description",
			"self":        false,
			"cidr_blocks": []interface{}{
				"10.0.0.1/32",
			},
		},
		map[string]interface{}{
			"protocol":    "tcp",
			"from_port":   int(443),
			"to_port":     int(443),
			"description": "block with description",
			"self":        false,
			"cidr_blocks": []interface{}{
				"10.0.0.2/32",
			},
		},
		map[string]interface{}{
			"protocol":    "tcp",
			"from_port":   int(443),
			"to_port":     int(443),
			"description": "block with description",
			"self":        false,
			"cidr_blocks": []interface{}{
				"10.0.0.3/32",
			},
		},
		map[string]interface{}{
			"protocol":    "tcp",
			"from_port":   int(443),
			"to_port":     int(443),
			"description": "block with another description",
			"self":        false,
			"cidr_blocks": []interface{}{
				"192.168.0.1/32",
			},
		},
		map[string]interface{}{
			"protocol":    "tcp",
			"from_port":   int(443),
			"to_port":     int(443),
			"description": "block with another description",
			"self":        false,
			"cidr_blocks": []interface{}{
				"192.168.0.2/32",
			},
		},
		map[string]interface{}{
			"protocol":    "-1",
			"from_port":   int(8000),
			"to_port":     int(8080),
			"description": "",
			"self":        false,
			"ipv6_cidr_blocks": []interface{}{
				"fd00::1/128",
			},
		},
		map[string]interface{}{
			"protocol":    "-1",
			"from_port":   int(8000),
			"to_port":     int(8080),
			"description": "",
			"self":        false,
			"ipv6_cidr_blocks": []interface{}{
				"fd00::2/128",
			},
		},
		map[string]interface{}{
			"protocol":    "-1",
			"from_port":   int(8000),
			"to_port":     int(8080),
			"description": "",
			"self":        false,
			"security_groups": schema.NewSet(schema.HashString, []interface{}{
				"sg-11111",
			}),
		},
		map[string]interface{}{
			"protocol":    "-1",
			"from_port":   int(8000),
			"to_port":     int(8080),
			"description": "",
			"self":        false,
			"security_groups": schema.NewSet(schema.HashString, []interface{}{
				"sg-22222",
			}),
		},
		map[string]interface{}{
			"protocol":    "-1",
			"from_port":   int(8000),
			"to_port":     int(8080),
			"description": "",
			"self":        false,
			"security_groups": schema.NewSet(schema.HashString, []interface{}{
				"sg-33333",
			}),
		},
		map[string]interface{}{
			"protocol":    "udp",
			"from_port":   int(10000),
			"to_port":     int(10000),
			"description": "",
			"self":        false,
			"prefix_list_ids": []interface{}{
				"pl-111111",
			},
		},
		map[string]interface{}{
			"protocol":    "udp",
			"from_port":   int(10000),
			"to_port":     int(10000),
			"description": "",
			"self":        false,
			"prefix_list_ids": []interface{}{
				"pl-222222",
			},
		},
	}

	expected_compact_set := schema.NewSet(tfec2.SecurityGroupRuleHash, expected_compact_list)
	actual_expanded_list := tfec2.SecurityGroupExpandRules(expected_compact_set).List()

	if calcSecurityGroupChecksum(expected_expanded_list) != calcSecurityGroupChecksum(actual_expanded_list) {
		t.Fatalf("error matching expanded set for tfec2.SecurityGroupExpandRules()")
	}

	actual_collapsed_list := tfec2.SecurityGroupCollapseRules("ingress", expected_expanded_list)

	if calcSecurityGroupChecksum(expected_compact_list) != calcSecurityGroupChecksum(actual_collapsed_list) {
		t.Fatalf("error matching collapsed set for tfec2.SecurityGroupCollapseRules()")
	}
}

func TestSecurityGroupIPPermGather(t *testing.T) {
	raw := []*ec2.IpPermission{
		{
			IpProtocol: aws.String("tcp"),
			FromPort:   aws.Int64(1),
			ToPort:     aws.Int64(int64(-1)),
			IpRanges:   []*ec2.IpRange{{CidrIp: aws.String("0.0.0.0/0")}},
			UserIdGroupPairs: []*ec2.UserIdGroupPair{
				{
					GroupId:     aws.String("sg-11111"),
					Description: aws.String("desc"),
				},
			},
		},
		{
			IpProtocol: aws.String("tcp"),
			FromPort:   aws.Int64(80),
			ToPort:     aws.Int64(80),
			UserIdGroupPairs: []*ec2.UserIdGroupPair{
				// VPC
				{
					GroupId: aws.String("sg-22222"),
				},
			},
		},
		{
			IpProtocol: aws.String("tcp"),
			FromPort:   aws.Int64(443),
			ToPort:     aws.Int64(443),
			UserIdGroupPairs: []*ec2.UserIdGroupPair{
				// Classic
				{
					UserId:    aws.String("12345"),
					GroupId:   aws.String("sg-33333"),
					GroupName: aws.String("ec2_classic"),
				},
				{
					UserId:    aws.String("amazon-elb"),
					GroupId:   aws.String("sg-d2c979d3"),
					GroupName: aws.String("amazon-elb-sg"),
				},
			},
		},
		{
			IpProtocol: aws.String("-1"),
			FromPort:   aws.Int64(0),
			ToPort:     aws.Int64(0),
			PrefixListIds: []*ec2.PrefixListId{
				{
					PrefixListId: aws.String("pl-12345678"),
					Description:  aws.String("desc"),
				},
			},
			UserIdGroupPairs: []*ec2.UserIdGroupPair{
				// VPC
				{
					GroupId: aws.String("sg-22222"),
				},
			},
		},
	}

	local := []map[string]interface{}{
		{
			"protocol":    "tcp",
			"from_port":   int64(1),
			"to_port":     int64(-1),
			"cidr_blocks": []string{"0.0.0.0/0"},
			"self":        true,
			"description": "desc",
		},
		{
			"protocol":  "tcp",
			"from_port": int64(80),
			"to_port":   int64(80),
			"security_groups": schema.NewSet(schema.HashString, []interface{}{
				"sg-22222",
			}),
		},
		{
			"protocol":  "tcp",
			"from_port": int64(443),
			"to_port":   int64(443),
			"security_groups": schema.NewSet(schema.HashString, []interface{}{
				"ec2_classic",
				"amazon-elb/amazon-elb-sg",
			}),
		},
		{
			"protocol":        "-1",
			"from_port":       int64(0),
			"to_port":         int64(0),
			"prefix_list_ids": []string{"pl-12345678"},
			"security_groups": schema.NewSet(schema.HashString, []interface{}{
				"sg-22222",
			}),
			"description": "desc",
		},
	}

	out := tfec2.SecurityGroupIPPermGather("sg-11111", raw, aws.String("12345"))
	for _, i := range out {
		// loop and match rules, because the ordering is not guarneteed
		for _, l := range local {
			if i["from_port"] == l["from_port"] {
				if i["to_port"] != l["to_port"] {
					t.Fatalf("to_port does not match")
				}

				if _, ok := i["cidr_blocks"]; ok {
					if !reflect.DeepEqual(i["cidr_blocks"], l["cidr_blocks"]) {
						t.Fatalf("error matching cidr_blocks")
					}
				}

				if _, ok := i["security_groups"]; ok {
					outSet := i["security_groups"].(*schema.Set)
					localSet := l["security_groups"].(*schema.Set)

					if !outSet.Equal(localSet) {
						t.Fatalf("Security Group sets are not equal")
					}
				}
			}
		}
	}
}

func TestAccVPCSecurityGroup_basic(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_name(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile(`security-group/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", ""),
					acctest.CheckResourceAttrAccountID(resourceName, "owner_id"),
					resource.TestCheckResourceAttr(resourceName, "revoke_rules_on_delete", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_basicEC2Classic(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckEC2Classic(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupEC2ClassicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_ec2Classic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSecurityGroupEC2ClassicExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", ""),
					acctest.CheckResourceAttrAccountID(resourceName, "owner_id"),
					resource.TestCheckResourceAttr(resourceName, "revoke_rules_on_delete", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", ""),
				),
			},
			{
				Config:                  testAccVPCSecurityGroupConfig_ec2Classic(rName),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_disappears(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_name(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceSecurityGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVPCSecurityGroup_nameGenerated(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_nameGenerated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					acctest.CheckResourceAttrNameGenerated(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", resource.UniqueIdPrefix),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/17017
func TestAccVPCSecurityGroup_nameTerraformPrefix(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix("terraform-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_name(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_namePrefix(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_namePrefix(rName, "tf-acc-test-prefix-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "name", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-test-prefix-"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/17017
func TestAccVPCSecurityGroup_namePrefixTerraform(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_namePrefix(rName, "terraform-test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "name", "terraform-test"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "terraform-test"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_tags(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
			{
				Config: testAccVPCSecurityGroupConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccVPCSecurityGroupConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccVPCSecurityGroup_allowAll(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_allowAll(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_sourceSecurityGroup(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_sourceSecurityGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_ipRangeAndSecurityGroupWithSameRules(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_ipRangeAndSecurityGroupWithSameRules(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_ipRangesWithSameRules(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_ipRangesWithSameRules(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_egressMode(t *testing.T) {
	var securityGroup1, securityGroup2, securityGroup3 ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_egressModeBlocks(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &securityGroup1),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
			{
				Config: testAccVPCSecurityGroupConfig_egressModeNoBlocks(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &securityGroup2),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "2"),
				),
			},
			{
				Config: testAccVPCSecurityGroupConfig_egressModeZeroed(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &securityGroup3),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "0"),
				),
			},
		},
	})
}

func TestAccVPCSecurityGroup_ingressMode(t *testing.T) {
	var securityGroup1, securityGroup2, securityGroup3 ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_ingressModeBlocks(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &securityGroup1),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
			{
				Config: testAccVPCSecurityGroupConfig_ingressModeNoBlocks(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &securityGroup2),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "2"),
				),
			},
			{
				Config: testAccVPCSecurityGroupConfig_ingressModeZeroed(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &securityGroup3),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "0"),
				),
			},
		},
	})
}

func TestAccVPCSecurityGroup_ruleGathering(t *testing.T) {
	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_ruleGathering(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":      "0",
						"description":        "egress for all ipv6",
						"from_port":          "0",
						"ipv6_cidr_blocks.#": "1",
						"ipv6_cidr_blocks.0": "::/0",
						"prefix_list_ids.#":  "0",
						"protocol":           "-1",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "0",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "0.0.0.0/0",
						"description":        "egress for all ipv4",
						"from_port":          "0",
						"ipv6_cidr_blocks.#": "0",
						"prefix_list_ids.#":  "0",
						"protocol":           "-1",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "0",
					}),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "5"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "192.168.0.0/16",
						"description":        "ingress from 192.168.0.0/16",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "80",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "0",
						"description":        "ingress from all ipv6",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "1",
						"ipv6_cidr_blocks.0": "::/0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "80",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "2",
						"cidr_blocks.0":      "10.0.2.0/24",
						"cidr_blocks.1":      "10.0.3.0/24",
						"description":        "ingress from 10.0.0.0/16",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "80",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "2",
						"cidr_blocks.0":      "10.0.0.0/24",
						"cidr_blocks.1":      "10.0.1.0/24",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "true",
						"to_port":            "80",
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

// This test should fail to destroy the Security Groups and VPC, due to a
// dependency cycle added outside of terraform's management. There is a sweeper
// 'aws_vpc' and 'aws_security_group' that cleans these up, however, the test is
// written to allow Terraform to clean it up because we do go and revoke the
// cyclic rules that were added.
func TestAccVPCSecurityGroup_forceRevokeRulesTrue(t *testing.T) {
	var primary ec2.SecurityGroup
	var secondary ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.primary"
	resourceName2 := "aws_security_group.secondary"

	// Add rules to create a cycle between primary and secondary. This prevents
	// Terraform/AWS from being able to destroy the groups
	testAddCycle := testAddRuleCycle(&primary, &secondary)
	// Remove the rules that created the cycle; Terraform/AWS can now destroy them
	testRemoveCycle := testRemoveRuleCycle(&primary, &secondary)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			// create the configuration with 2 security groups, then create a
			// dependency cycle such that they cannot be deleted
			{
				Config: testAccVPCSecurityGroupConfig_revokeBase(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &primary),
					testAccCheckSecurityGroupExists(resourceName2, &secondary),
					testAddCycle,
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
			// Verify the DependencyViolation error by using a configuration with the
			// groups removed. Terraform tries to destroy them but cannot. Expect a
			// DependencyViolation error
			{
				Config:      testAccVPCSecurityGroupConfig_revokeBaseRemoved(rName),
				ExpectError: regexp.MustCompile("DependencyViolation"),
			},
			// Restore the config (a no-op plan) but also remove the dependencies
			// between the groups with testRemoveCycle
			{
				Config: testAccVPCSecurityGroupConfig_revokeBase(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &primary),
					testAccCheckSecurityGroupExists(resourceName2, &secondary),
					testRemoveCycle,
				),
			},
			// Again try to apply the config with the sgs removed; it should work
			{
				Config: testAccVPCSecurityGroupConfig_revokeBaseRemoved(rName),
			},
			////
			// now test with revoke_rules_on_delete
			////
			// create the configuration with 2 security groups, then create a
			// dependency cycle such that they cannot be deleted. In this
			// configuration, each Security Group has `revoke_rules_on_delete`
			// specified, and should delete with no issue
			{
				Config: testAccVPCSecurityGroupConfig_revokeTrue(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &primary),
					testAccCheckSecurityGroupExists(resourceName2, &secondary),
					testAddCycle,
				),
			},
			// Again try to apply the config with the sgs removed; it should work,
			// because we've told the SGs to forcefully revoke their rules first
			{
				Config: testAccVPCSecurityGroupConfig_revokeBaseRemoved(rName),
			},
		},
	})
}

func TestAccVPCSecurityGroup_forceRevokeRulesFalse(t *testing.T) {
	var primary ec2.SecurityGroup
	var secondary ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.primary"
	resourceName2 := "aws_security_group.secondary"

	// Add rules to create a cycle between primary and secondary. This prevents
	// Terraform/AWS from being able to destroy the groups
	testAddCycle := testAddRuleCycle(&primary, &secondary)
	// Remove the rules that created the cycle; Terraform/AWS can now destroy them
	testRemoveCycle := testRemoveRuleCycle(&primary, &secondary)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			// create the configuration with 2 security groups, then create a
			// dependency cycle such that they cannot be deleted. These Security
			// Groups are configured to explicitly not revoke rules on delete,
			// `revoke_rules_on_delete = false`
			{
				Config: testAccVPCSecurityGroupConfig_revokeFalse(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &primary),
					testAccCheckSecurityGroupExists(resourceName2, &secondary),
					testAddCycle,
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
			// Verify the DependencyViolation error by using a configuration with the
			// groups removed, and the Groups not configured to revoke their ruls.
			// Terraform tries to destroy them but cannot. Expect a
			// DependencyViolation error
			{
				Config:      testAccVPCSecurityGroupConfig_revokeBaseRemoved(rName),
				ExpectError: regexp.MustCompile("DependencyViolation"),
			},
			// Restore the config (a no-op plan) but also remove the dependencies
			// between the groups with testRemoveCycle
			{
				Config: testAccVPCSecurityGroupConfig_revokeFalse(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &primary),
					testAccCheckSecurityGroupExists(resourceName2, &secondary),
					testRemoveCycle,
				),
			},
			// Again try to apply the config with the sgs removed; it should work
			{
				Config: testAccVPCSecurityGroupConfig_revokeBaseRemoved(rName),
			},
		},
	})
}

func TestAccVPCSecurityGroup_change(t *testing.T) {
	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
			{
				Config: testAccVPCSecurityGroupConfig_changed(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "10.0.0.0/8",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "10.0.0.0/8",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "9000",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "2",
						"cidr_blocks.0":      "0.0.0.0/0",
						"cidr_blocks.1":      "10.0.0.0/8",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
				),
			},
		},
	})
}

func TestAccVPCSecurityGroup_ipv6(t *testing.T) {
	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_ipv6(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":      "0",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "1",
						"ipv6_cidr_blocks.0": "::/0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "0",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "1",
						"ipv6_cidr_blocks.0": "::/0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_self(t *testing.T) {
	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	checkSelf := func(s *terraform.State) (err error) {
		if len(group.IpPermissions) > 0 &&
			len(group.IpPermissions[0].UserIdGroupPairs) > 0 &&
			aws.StringValue(group.IpPermissions[0].UserIdGroupPairs[0].GroupId) == aws.StringValue(group.GroupId) {
			return nil
		}

		return fmt.Errorf("Security Group does not contain \"self\" rule: %#v", group)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_self(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"protocol":  "tcp",
						"from_port": "80",
						"to_port":   "8000",
						"self":      "true",
					}),
					checkSelf,
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_vpc(t *testing.T) {
	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_vpc(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"protocol":      "tcp",
						"from_port":     "80",
						"to_port":       "8000",
						"cidr_blocks.#": "1",
						"cidr_blocks.0": "10.0.0.0/8",
					}),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"protocol":      "tcp",
						"from_port":     "80",
						"to_port":       "8000",
						"cidr_blocks.#": "1",
						"cidr_blocks.0": "10.0.0.0/8",
					}),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_vpcNegOneIngress(t *testing.T) {
	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_vpcNegativeOneIngress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"protocol":      "-1",
						"from_port":     "0",
						"to_port":       "0",
						"cidr_blocks.#": "1",
						"cidr_blocks.0": "10.0.0.0/8",
					}),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_vpcProtoNumIngress(t *testing.T) {
	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_vpcProtocolNumberIngress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"protocol":      "50",
						"from_port":     "0",
						"to_port":       "0",
						"cidr_blocks.#": "1",
						"cidr_blocks.0": "10.0.0.0/8",
					}),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_multiIngress(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test1"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_multiIngress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_ruleDescription(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_ruleDescription(rName, "Egress description", "Ingress description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "10.0.0.0/8",
						"description":        "Egress description",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"prefix_list_ids.#":  "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "10.0.0.0/8",
						"description":        "Ingress description",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
			// Change just the rule descriptions.
			{
				Config: testAccVPCSecurityGroupConfig_ruleDescription(rName, "New egress description", "New ingress description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "10.0.0.0/8",
						"description":        "New egress description",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"prefix_list_ids.#":  "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "10.0.0.0/8",
						"description":        "New ingress description",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
				),
			},
			// Remove just the rule descriptions.
			{
				Config: testAccVPCSecurityGroupConfig_emptyRuleDescription(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":     "1",
						"cidr_blocks.0":     "10.0.0.0/8",
						"description":       "",
						"from_port":         "80",
						"protocol":          "tcp",
						"security_groups.#": "0",
						"self":              "false",
						"to_port":           "8000",
					}),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":     "1",
						"cidr_blocks.0":     "10.0.0.0/8",
						"description":       "",
						"from_port":         "80",
						"protocol":          "tcp",
						"security_groups.#": "0",
						"self":              "false",
						"to_port":           "8000",
					}),
				),
			},
		},
	})
}

func TestAccVPCSecurityGroup_defaultEgressVPC(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_defaultEgress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

// Testing drift detection with groups containing the same port and types
func TestAccVPCSecurityGroup_drift(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_drift(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "10.0.0.0/8",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "206.0.0.0/8",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				// In rules with cidr_block drift, import only creates a single ingress
				// rule with the cidr_blocks de-normalized. During subsequent apply, its
				// normalized to create the 2 ingress rules seen in checks above.
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete", "ingress", "egress"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_driftComplex(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test1"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_driftComplex(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "10.0.0.0/8",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"prefix_list_ids.#":  "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "206.0.0.0/8",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"prefix_list_ids.#":  "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "10.0.0.0/8",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "206.0.0.0/8",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				// In rules with cidr_block drift, import only creates a single ingress
				// rule with the cidr_blocks de-normalized. During subsequent apply, its
				// normalized to create the 2 ingress rules seen in checks above.
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete", "ingress", "egress"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_invalidCIDRBlock(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccVPCSecurityGroupConfig_invalidIngressCIDR,
				ExpectError: regexp.MustCompile("invalid CIDR address: 1.2.3.4/33"),
			},
			{
				Config:      testAccVPCSecurityGroupConfig_invalidEgressCIDR,
				ExpectError: regexp.MustCompile("invalid CIDR address: 1.2.3.4/33"),
			},
			{
				Config:      testAccVPCSecurityGroupConfig_invalidIPv6IngressCIDR,
				ExpectError: regexp.MustCompile("invalid CIDR address: ::/244"),
			},
			{
				Config:      testAccVPCSecurityGroupConfig_invalidIPv6EgressCIDR,
				ExpectError: regexp.MustCompile("invalid CIDR address: ::/244"),
			},
		},
	})
}

func TestAccVPCSecurityGroup_cidrAndGroups(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test1"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_combinedCIDRAndGroups(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_ingressWithCIDRAndSGsVPC(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test1"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_ingressWithCIDRAndSGs(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "10.0.0.0/8",
						"description":        "",
						"from_port":          "80",
						"ipv6_cidr_blocks.#": "0",
						"prefix_list_ids.#":  "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "8000",
					}),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "192.168.0.1/32",
						"description":        "",
						"from_port":          "22",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "22",
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_ingressWithCIDRAndSGsClassic(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test1"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckEC2Classic(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupEC2ClassicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_ingressWithCIDRAndSGsEC2Classic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupEC2ClassicExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ingress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "192.168.0.1/32",
						"description":        "",
						"from_port":          "22",
						"ipv6_cidr_blocks.#": "0",
						"protocol":           "tcp",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "22",
					}),
				),
			},
			{
				Config:                  testAccVPCSecurityGroupConfig_ingressWithCIDRAndSGsEC2Classic(rName),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_egressWithPrefixList(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_prefixListEgress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_ingressWithPrefixList(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_prefixListIngress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_ipv4AndIPv6Egress(t *testing.T) {
	var group ec2.SecurityGroup
	resourceName := "aws_security_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_ipv4andIPv6Egress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":      "1",
						"cidr_blocks.0":      "0.0.0.0/0",
						"description":        "",
						"from_port":          "0",
						"ipv6_cidr_blocks.#": "0",
						"prefix_list_ids.#":  "0",
						"protocol":           "-1",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "0",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "egress.*", map[string]string{
						"cidr_blocks.#":      "0",
						"description":        "",
						"from_port":          "0",
						"ipv6_cidr_blocks.#": "1",
						"ipv6_cidr_blocks.0": "::/0",
						"prefix_list_ids.#":  "0",
						"protocol":           "-1",
						"security_groups.#":  "0",
						"self":               "false",
						"to_port":            "0",
					}),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revoke_rules_on_delete", "egress"},
			},
		},
	})
}

func TestAccVPCSecurityGroup_failWithDiffMismatch(t *testing.T) {
	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCSecurityGroupConfig_failWithDiffMismatch(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "egress.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "2"),
				),
			},
		},
	})
}

func TestAccVPCSecurityGroup_ruleLimitExceededAppend(t *testing.T) {
	ruleLimit := testAccSecurityGroupRulesPerGroupLimitFromEnv()

	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			// create a valid SG just under the limit
			{
				Config: testAccVPCSecurityGroupConfig_ruleLimit(rName, 0, ruleLimit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					testAccCheckSecurityGroupRuleCount(&group, 0, ruleLimit),
					resource.TestCheckResourceAttr(resourceName, "egress.#", strconv.Itoa(ruleLimit)),
				),
			},
			// append a rule to step over the limit
			{
				Config:      testAccVPCSecurityGroupConfig_ruleLimit(rName, 0, ruleLimit+1),
				ExpectError: regexp.MustCompile("RulesPerSecurityGroupLimitExceeded"),
			},
			{
				PreConfig: func() {
					// should have the original rules still
					err := testSecurityGroupRuleCount(aws.StringValue(group.GroupId), 0, ruleLimit)
					if err != nil {
						t.Fatalf("PreConfig check failed: %s", err)
					}
				},
				// running the original config again now should restore the rules
				Config: testAccVPCSecurityGroupConfig_ruleLimit(rName, 0, ruleLimit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					testAccCheckSecurityGroupRuleCount(&group, 0, ruleLimit),
					resource.TestCheckResourceAttr(resourceName, "egress.#", strconv.Itoa(ruleLimit)),
				),
			},
		},
	})
}

func TestAccVPCSecurityGroup_ruleLimitCIDRBlockExceededAppend(t *testing.T) {
	ruleLimit := testAccSecurityGroupRulesPerGroupLimitFromEnv()

	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			// create a valid SG just under the limit
			{
				Config: testAccVPCSecurityGroupConfig_cidrBlockRuleLimit(rName, 0, ruleLimit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					testAccCheckSecurityGroupRuleCount(&group, 0, 1),
				),
			},
			// append a rule to step over the limit
			{
				Config:      testAccVPCSecurityGroupConfig_cidrBlockRuleLimit(rName, 0, ruleLimit+1),
				ExpectError: regexp.MustCompile("RulesPerSecurityGroupLimitExceeded"),
			},
			{
				PreConfig: func() {
					// should have the original cidr blocks still in 1 rule
					err := testSecurityGroupRuleCount(aws.StringValue(group.GroupId), 0, 1)
					if err != nil {
						t.Fatalf("PreConfig check failed: %s", err)
					}

					id := aws.StringValue(group.GroupId)

					conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

					match, err := tfec2.FindSecurityGroupByID(conn, id)
					if tfresource.NotFound(err) {
						t.Fatalf("PreConfig check failed: Security Group (%s) not found: %s", id, err)
					}
					if err != nil {
						t.Fatalf("PreConfig check failed: %s", err)
					}

					if cidrCount := len(match.IpPermissionsEgress[0].IpRanges); cidrCount != ruleLimit {
						t.Fatalf("PreConfig check failed: rule does not have previous IP ranges, has %d", cidrCount)
					}
				},
				// running the original config again now should restore the rules
				Config: testAccVPCSecurityGroupConfig_cidrBlockRuleLimit(rName, 0, ruleLimit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					testAccCheckSecurityGroupRuleCount(&group, 0, 1),
				),
			},
		},
	})
}

func TestAccVPCSecurityGroup_ruleLimitExceededPrepend(t *testing.T) {
	ruleLimit := testAccSecurityGroupRulesPerGroupLimitFromEnv()

	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			// create a valid SG just under the limit
			{
				Config: testAccVPCSecurityGroupConfig_ruleLimit(rName, 0, ruleLimit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					testAccCheckSecurityGroupRuleCount(&group, 0, ruleLimit),
				),
			},
			// prepend a rule to step over the limit
			{
				Config:      testAccVPCSecurityGroupConfig_ruleLimit(rName, 1, ruleLimit+1),
				ExpectError: regexp.MustCompile("RulesPerSecurityGroupLimitExceeded"),
			},
			{
				PreConfig: func() {
					// should have the original rules still (limit - 1 because of the shift)
					err := testSecurityGroupRuleCount(aws.StringValue(group.GroupId), 0, ruleLimit-1)
					if err != nil {
						t.Fatalf("PreConfig check failed: %s", err)
					}
				},
				// running the original config again now should restore the rules
				Config: testAccVPCSecurityGroupConfig_ruleLimit(rName, 0, ruleLimit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					testAccCheckSecurityGroupRuleCount(&group, 0, ruleLimit),
				),
			},
		},
	})
}

func TestAccVPCSecurityGroup_ruleLimitExceededAllNew(t *testing.T) {
	ruleLimit := testAccSecurityGroupRulesPerGroupLimitFromEnv()

	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			// create a valid SG just under the limit
			{
				Config: testAccVPCSecurityGroupConfig_ruleLimit(rName, 0, ruleLimit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					testAccCheckSecurityGroupRuleCount(&group, 0, ruleLimit),
				),
			},
			// add a rule to step over the limit with entirely new rules
			{
				Config:      testAccVPCSecurityGroupConfig_ruleLimit(rName, 100, ruleLimit+1),
				ExpectError: regexp.MustCompile("RulesPerSecurityGroupLimitExceeded"),
			},
			{
				// all the rules should have been revoked and the add failed
				PreConfig: func() {
					err := testSecurityGroupRuleCount(aws.StringValue(group.GroupId), 0, 0)
					if err != nil {
						t.Fatalf("PreConfig check failed: %s", err)
					}
				},
				// running the original config again now should restore the rules
				Config: testAccVPCSecurityGroupConfig_ruleLimit(rName, 0, ruleLimit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					testAccCheckSecurityGroupRuleCount(&group, 0, ruleLimit),
				),
			},
		},
	})
}

func TestAccVPCSecurityGroup_rulesDropOnError(t *testing.T) {
	var group ec2.SecurityGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			// Create a valid security group with some rules and make sure it exists
			{
				Config: testAccVPCSecurityGroupConfig_rulesDropOnErrorInit(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
				),
			},
			// Add a bad rule to trigger API error
			{
				Config:      testAccVPCSecurityGroupConfig_rulesDropOnErrorAddBadRule(rName),
				ExpectError: regexp.MustCompile("InvalidGroupId.Malformed"),
			},
			// All originally added rules must survive. This will return non-empty plan if anything changed.
			{
				Config:   testAccVPCSecurityGroupConfig_rulesDropOnErrorInit(rName),
				PlanOnly: true,
			},
		},
	})
}

// cycleIPPermForGroup returns an IpPermission struct with a configured
// UserIdGroupPair for the groupid given. Used in
// TestAccAWSSecurityGroup_forceRevokeRules_should_fail to create a cyclic rule
// between 2 security groups
func cycleIPPermForGroup(groupId string) *ec2.IpPermission {
	var perm ec2.IpPermission
	perm.FromPort = aws.Int64(0)
	perm.ToPort = aws.Int64(0)
	perm.IpProtocol = aws.String("icmp")
	perm.UserIdGroupPairs = make([]*ec2.UserIdGroupPair, 1)
	perm.UserIdGroupPairs[0] = &ec2.UserIdGroupPair{
		GroupId: aws.String(groupId),
	}
	return &perm
}

// testAddRuleCycle returns a TestCheckFunc to use at the end of a test, such
// that a Security Group Rule cyclic dependency will be created between the two
// Security Groups. A companion function, testRemoveRuleCycle, will undo this.
func testAddRuleCycle(primary, secondary *ec2.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if primary.GroupId == nil {
			return fmt.Errorf("Primary SG not set for TestAccAWSSecurityGroup_forceRevokeRules_should_fail")
		}
		if secondary.GroupId == nil {
			return fmt.Errorf("Secondary SG not set for TestAccAWSSecurityGroup_forceRevokeRules_should_fail")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		// cycle from primary to secondary
		perm1 := cycleIPPermForGroup(aws.StringValue(secondary.GroupId))
		// cycle from secondary to primary
		perm2 := cycleIPPermForGroup(aws.StringValue(primary.GroupId))

		req1 := &ec2.AuthorizeSecurityGroupEgressInput{
			GroupId:       primary.GroupId,
			IpPermissions: []*ec2.IpPermission{perm1},
		}
		req2 := &ec2.AuthorizeSecurityGroupEgressInput{
			GroupId:       secondary.GroupId,
			IpPermissions: []*ec2.IpPermission{perm2},
		}

		var err error
		_, err = conn.AuthorizeSecurityGroupEgress(req1)
		if err != nil {
			return fmt.Errorf("Error authorizing primary security group %s rules: %w", aws.StringValue(primary.GroupId), err)
		}
		_, err = conn.AuthorizeSecurityGroupEgress(req2)
		if err != nil {
			return fmt.Errorf("Error authorizing secondary security group %s rules: %w", aws.StringValue(secondary.GroupId), err)
		}
		return nil
	}
}

// testRemoveRuleCycle removes the cyclic dependency between two security groups
// that was added in testAddRuleCycle
func testRemoveRuleCycle(primary, secondary *ec2.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if primary.GroupId == nil {
			return fmt.Errorf("Primary SG not set for TestAccAWSSecurityGroup_forceRevokeRules_should_fail")
		}
		if secondary.GroupId == nil {
			return fmt.Errorf("Secondary SG not set for TestAccAWSSecurityGroup_forceRevokeRules_should_fail")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn
		for _, sg := range []*ec2.SecurityGroup{primary, secondary} {
			var err error
			if sg.IpPermissions != nil {
				req := &ec2.RevokeSecurityGroupIngressInput{
					GroupId:       sg.GroupId,
					IpPermissions: sg.IpPermissions,
				}

				if _, err = conn.RevokeSecurityGroupIngress(req); err != nil {
					return fmt.Errorf("Error revoking default ingress rule for Security Group in testRemoveCycle (%s): %w", aws.StringValue(primary.GroupId), err)
				}
			}

			if sg.IpPermissionsEgress != nil {
				req := &ec2.RevokeSecurityGroupEgressInput{
					GroupId:       sg.GroupId,
					IpPermissions: sg.IpPermissionsEgress,
				}

				if _, err = conn.RevokeSecurityGroupEgress(req); err != nil {
					return fmt.Errorf("Error revoking default egress rule for Security Group in testRemoveCycle (%s): %w", aws.StringValue(sg.GroupId), err)
				}
			}
		}
		return nil
	}
}

func testAccCheckSecurityGroupDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_security_group" {
			continue
		}

		_, err := tfec2.FindSecurityGroupByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("VPC Security Group (%s) still exists.", rs.Primary.ID)
	}

	return nil
}

func testAccCheckSecurityGroupEC2ClassicDestroy(s *terraform.State) error { // nosemgrep:ci.ec2-in-func-name
	conn := acctest.ProviderEC2Classic.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_security_group" {
			continue
		}

		_, err := tfec2.FindSecurityGroupByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("EC2 Classic Security Group (%s) still exists.", rs.Primary.ID)
	}

	return nil
}

func testAccCheckSecurityGroupExists(n string, v *ec2.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC Security Group ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		output, err := tfec2.FindSecurityGroupByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckSecurityGroupEC2ClassicExists(n string, v *ec2.SecurityGroup) resource.TestCheckFunc { // nosemgrep:ci.ec2-in-func-name
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 Classic Security Group ID is set")
		}

		conn := acctest.ProviderEC2Classic.Meta().(*conns.AWSClient).EC2Conn

		output, err := tfec2.FindSecurityGroupByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

// testAccSecurityGroupRulesPerGroupLimitFromEnv returns security group rules per group limit
// Currently this information is not available from any EC2 or Trusted Advisor API
// Prefers the EC2_SECURITY_GROUP_RULES_PER_GROUP_LIMIT environment variable or defaults to 50
func testAccSecurityGroupRulesPerGroupLimitFromEnv() int {
	const defaultLimit = 50
	const envVar = "EC2_SECURITY_GROUP_RULES_PER_GROUP_LIMIT"

	envLimitStr := os.Getenv(envVar)
	if envLimitStr == "" {
		return defaultLimit
	}
	envLimitInt, err := strconv.Atoi(envLimitStr)
	if err != nil {
		log.Printf("[WARN] Error converting %q environment variable value %q to integer: %s", envVar, envLimitStr, err)
		return defaultLimit
	}
	if envLimitInt <= 50 {
		return defaultLimit
	}
	return envLimitInt
}

func testAccCheckSecurityGroupRuleCount(group *ec2.SecurityGroup, expectedIngressCount, expectedEgressCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id := aws.StringValue(group.GroupId)
		return testSecurityGroupRuleCount(id, expectedIngressCount, expectedEgressCount)
	}
}

func testSecurityGroupRuleCount(id string, expectedIngressCount, expectedEgressCount int) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	group, err := tfec2.FindSecurityGroupByID(conn, id)
	if tfresource.NotFound(err) {
		return fmt.Errorf("Security Group (%s) not found: %w", id, err)
	}
	if err != nil {
		return err
	}

	if actual := len(group.IpPermissions); actual != expectedIngressCount {
		return fmt.Errorf("Security group ingress rule count %d does not match %d", actual, expectedIngressCount)
	}

	if actual := len(group.IpPermissionsEgress); actual != expectedEgressCount {
		return fmt.Errorf("Security group egress rule count %d does not match %d", actual, expectedEgressCount)
	}

	return nil
}

func testAccVPCSecurityGroupConfig_name(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ec2Classic(rName string) string { // nosemgrep:ci.ec2-in-func-name
	return acctest.ConfigCompose(acctest.ConfigEC2ClassicRegionProvider(), fmt.Sprintf(`
resource "aws_security_group" "test" {
  name = %[1]q
}
`, rName))
}

func testAccVPCSecurityGroupConfig_nameGenerated(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  vpc_id = aws_vpc.test.id
}
`, rName)
}

func testAccVPCSecurityGroupConfig_namePrefix(rName, namePrefix string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name_prefix = %[2]q
  vpc_id      = aws_vpc.test.id
}
`, rName, namePrefix)
}

func testAccVPCSecurityGroupConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccVPCSecurityGroupConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccVPCSecurityGroupConfig_ruleLimit(rName string, egressStartIndex, egressRulesCount int) string {
	var egressRules strings.Builder
	for i := egressStartIndex; i < egressRulesCount+egressStartIndex; i++ {
		fmt.Fprintf(&egressRules, `
  egress {
    protocol    = "tcp"
    from_port   = "${80 + %[1]d}"
    to_port     = "${80 + %[1]d}"
    cidr_blocks = ["${cidrhost("10.1.0.0/16", %[1]d)}/32"]
  }
`, i)
	}

	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  # egress rules to exhaust the limit
  %[2]s
}
`, rName, egressRules.String())
}

func testAccVPCSecurityGroupConfig_cidrBlockRuleLimit(rName string, egressStartIndex, egressRulesCount int) string {
	var cidrBlocks strings.Builder
	for i := egressStartIndex; i < egressRulesCount+egressStartIndex; i++ {
		fmt.Fprintf(&cidrBlocks, `
		"${cidrhost("10.1.0.0/16", %[1]d)}/32",
`, i)
	}

	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  egress {
    protocol  = "tcp"
    from_port = "80"
    to_port   = "80"
    # cidr_blocks to exhaust the limit
    cidr_blocks = [
		%[2]s
    ]
  }
}
`, rName, cidrBlocks.String())
}

func testAccVPCSecurityGroupConfig_emptyRuleDescription(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    protocol    = "6"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
    description = ""
  }

  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
    description = ""
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ipv6(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    protocol         = "6"
    from_port        = 80
    to_port          = 8000
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    protocol         = "tcp"
    from_port        = 80
    to_port          = 8000
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    protocol    = "6"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_revokeBaseRemoved(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_revokeBase(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "primary" {
  name   = "%[1]s-primary"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "secondary" {
  name   = "%[1]s-secondary"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_revokeFalse(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "primary" {
  name   = "%[1]s-primary"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  revoke_rules_on_delete = false
}

resource "aws_security_group" "secondary" {
  name   = "%[1]s-secondary"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  revoke_rules_on_delete = false
}
`, rName)
}

func testAccVPCSecurityGroupConfig_revokeTrue(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "primary" {
  name   = "%[1]s-primary"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  revoke_rules_on_delete = true
}

resource "aws_security_group" "secondary" {
  name   = "%[1]s-secondary"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  revoke_rules_on_delete = true
}
`, rName)
}

func testAccVPCSecurityGroupConfig_changed(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 9000
    cidr_blocks = ["10.0.0.0/8"]
  }

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["0.0.0.0/0", "10.0.0.0/8"]
  }

  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ruleDescription(rName, egressDescription, ingressDescription string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    protocol    = "6"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
    description = %[2]q
  }

  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
    description = %[3]q
  }

  tags = {
    Name = %[1]q
  }
}
`, rName, ingressDescription, egressDescription)
}

func testAccVPCSecurityGroupConfig_self(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    protocol  = "tcp"
    from_port = 80
    to_port   = 8000
    self      = true
  }

  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_vpc(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_vpcNegativeOneIngress(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["10.0.0.0/8"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_vpcProtocolNumberIngress(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    protocol    = "50"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["10.0.0.0/8"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_multiIngress(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test1" {
  name   = "%[1]s-1"
  vpc_id = aws_vpc.test.id

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test2" {
  name   = "%[1]s-2"
  vpc_id = aws_vpc.test.id

  ingress {
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["10.0.0.0/8"]
  }

  ingress {
    protocol    = "tcp"
    from_port   = 800
    to_port     = 800
    cidr_blocks = ["10.0.0.0/8"]
  }

  ingress {
    protocol        = "tcp"
    from_port       = 80
    to_port         = 8000
    security_groups = [aws_security_group.test1.id]
  }

  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_defaultEgress(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_drift(rName string) string {
	return fmt.Sprintf(`
resource "aws_security_group" "test" {
  name = %[1]q

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["206.0.0.0/8"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_driftComplex(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test2" {
  name   = "%[1]s-2"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test1" {
  name   = "%[1]s-1"
  vpc_id = aws_vpc.test.id

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["206.0.0.0/8"]
  }

  ingress {
    protocol        = "tcp"
    from_port       = 22
    to_port         = 22
    security_groups = [aws_security_group.test2.id]
  }

  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["206.0.0.0/8"]
  }

  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }

  egress {
    protocol        = "tcp"
    from_port       = 22
    to_port         = 22
    security_groups = [aws_security_group.test2.id]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

const testAccVPCSecurityGroupConfig_invalidIngressCIDR = `
resource "aws_security_group" "test" {
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["1.2.3.4/33"]
  }
}
`

const testAccVPCSecurityGroupConfig_invalidEgressCIDR = `
resource "aws_security_group" "test" {
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["1.2.3.4/33"]
  }
}
`

const testAccVPCSecurityGroupConfig_invalidIPv6IngressCIDR = `
resource "aws_security_group" "test" {
  ingress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    ipv6_cidr_blocks = ["::/244"]
  }
}
`

const testAccVPCSecurityGroupConfig_invalidIPv6EgressCIDR = `
resource "aws_security_group" "test" {
  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    ipv6_cidr_blocks = ["::/244"]
  }
}
`

func testAccVPCSecurityGroupConfig_combinedCIDRAndGroups(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test2" {
  name   = "%[1]s-2"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test3" {
  name   = "%[1]s-3"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test4" {
  name   = "%[1]s-4"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test1" {
  name   = "%[1]s-1"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16", "10.1.0.0/16", "10.7.0.0/16"]

    security_groups = [
      aws_security_group.test2.id,
      aws_security_group.test3.id,
      aws_security_group.test4.id,
    ]
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ingressWithCIDRAndSGs(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test2" {
  name   = "%[1]s-2"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test1" {
  name   = "%[1]s-1"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  ingress {
    protocol  = "tcp"
    from_port = "22"
    to_port   = "22"

    cidr_blocks = [
      "192.168.0.1/32",
    ]
  }

  ingress {
    protocol        = "tcp"
    from_port       = 80
    to_port         = 8000
    cidr_blocks     = ["10.0.0.0/8"]
    security_groups = [aws_security_group.test2.id]
  }

  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = ["10.0.0.0/8"]
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ingressWithCIDRAndSGsEC2Classic(rName string) string { // nosemgrep:ci.ec2-in-func-name
	return acctest.ConfigCompose(acctest.ConfigEC2ClassicRegionProvider(), fmt.Sprintf(`
resource "aws_security_group" "test2" {
  name = "%[1]s-2"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test1" {
  name = "%[1]s-1"

  ingress {
    protocol  = "tcp"
    from_port = "22"
    to_port   = "22"

    cidr_blocks = [
      "192.168.0.1/32",
    ]
  }

  ingress {
    protocol        = "tcp"
    from_port       = 80
    to_port         = 8000
    cidr_blocks     = ["10.0.0.0/8"]
    security_groups = [aws_security_group.test2.name]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

// fails to apply in one pass with the error "diffs didn't match during apply"
// GH-2027
func testAccVPCSecurityGroupConfig_failWithDiffMismatch(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test3" {
  vpc_id = aws_vpc.main.id
  name   = "%[1]s-3"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test2" {
  vpc_id = aws_vpc.main.id
  name   = "%[1]s-2"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test1" {
  vpc_id = aws_vpc.main.id
  name   = "%[1]s-1"

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [aws_security_group.test2.id]
  }

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [aws_security_group.test3.id]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_allowAll(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group_rule" "allow_all-1" {
  type        = "ingress"
  from_port   = 0
  to_port     = 65535
  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.test.id
}

resource "aws_security_group_rule" "allow_all-2" {
  type      = "ingress"
  from_port = 65534
  to_port   = 65535
  protocol  = "tcp"

  self              = true
  security_group_id = aws_security_group.test.id
}
`, rName)
}

func testAccVPCSecurityGroupConfig_sourceSecurityGroup(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = "%[1]s-1"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test2" {
  name   = "%[1]s-2"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test3" {
  name   = "%[1]s-3"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group_rule" "allow_test2" {
  type      = "ingress"
  from_port = 0
  to_port   = 0
  protocol  = "tcp"

  source_security_group_id = aws_security_group.test.id
  security_group_id        = aws_security_group.test2.id
}

resource "aws_security_group_rule" "allow_test3" {
  type      = "ingress"
  from_port = 0
  to_port   = 0
  protocol  = "tcp"

  source_security_group_id = aws_security_group.test.id
  security_group_id        = aws_security_group.test3.id
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ipRangeAndSecurityGroupWithSameRules(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = "%[1]s-1"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test2" {
  name   = "%[1]s-2"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group_rule" "allow_security_group" {
  type      = "ingress"
  from_port = 0
  to_port   = 0
  protocol  = "tcp"

  source_security_group_id = aws_security_group.test2.id
  security_group_id        = aws_security_group.test.id
}

resource "aws_security_group_rule" "allow_cidr_block" {
  type      = "ingress"
  from_port = 0
  to_port   = 0
  protocol  = "tcp"

  cidr_blocks       = ["10.0.0.0/32"]
  security_group_id = aws_security_group.test.id
}

resource "aws_security_group_rule" "allow_ipv6_cidr_block" {
  type      = "ingress"
  from_port = 0
  to_port   = 0
  protocol  = "tcp"

  ipv6_cidr_blocks  = ["::/0"]
  security_group_id = aws_security_group.test.id
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ipRangesWithSameRules(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group_rule" "allow_cidr_block" {
  type      = "ingress"
  from_port = 0
  to_port   = 0
  protocol  = "tcp"

  cidr_blocks       = ["10.0.0.0/32"]
  security_group_id = aws_security_group.test.id
}

resource "aws_security_group_rule" "allow_ipv6_cidr_block" {
  type      = "ingress"
  from_port = 0
  to_port   = 0
  protocol  = "tcp"

  ipv6_cidr_blocks  = ["::/0"]
  security_group_id = aws_security_group.test.id
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ipv4andIPv6Egress(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_prefixListEgress(rName string) string {
	return fmt.Sprintf(`
data "aws_region" "current" {}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc_endpoint" "test" {
  vpc_id          = aws_vpc.test.id
  service_name    = "com.amazonaws.${data.aws_region.current.name}.s3"
  route_table_ids = [aws_route_table.test.id]

  tags = {
    Name = %[1]q
  }

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAll",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "*",
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  egress {
    protocol        = "-1"
    from_port       = 0
    to_port         = 0
    prefix_list_ids = [aws_vpc_endpoint.test.prefix_list_id]
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_prefixListIngress(rName string) string {
	return fmt.Sprintf(`
data "aws_region" "current" {}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc_endpoint" "test" {
  vpc_id          = aws_vpc.test.id
  service_name    = "com.amazonaws.${data.aws_region.current.name}.s3"
  route_table_ids = [aws_route_table.test.id]

  tags = {
    Name = %[1]q
  }

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAll",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "*",
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  ingress {
    protocol        = "-1"
    from_port       = 0
    to_port         = 0
    prefix_list_ids = [aws_vpc_endpoint.test.prefix_list_id]
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ruleGathering(rName string) string {
	return fmt.Sprintf(`
data "aws_region" "current" {}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc_endpoint" "test" {
  vpc_id          = aws_vpc.test.id
  service_name    = "com.amazonaws.${data.aws_region.current.name}.s3"
  route_table_ids = [aws_route_table.test.id]

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAll",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "*",
      "Resource": "*"
    }
  ]
}
POLICY

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "source1" {
  name   = "%[1]s-source1"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "source2" {
  name   = "%[1]s-source2"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["10.0.0.0/24", "10.0.1.0/24"]
    self        = true
  }

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["10.0.2.0/24", "10.0.3.0/24"]
    description = "ingress from 10.0.0.0/16"
  }

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["192.168.0.0/16"]
    description = "ingress from 192.168.0.0/16"
  }

  ingress {
    protocol         = "tcp"
    from_port        = 80
    to_port          = 80
    ipv6_cidr_blocks = ["::/0"]
    description      = "ingress from all ipv6"
  }

  ingress {
    protocol        = "tcp"
    from_port       = 80
    to_port         = 80
    security_groups = [aws_security_group.source1.id, aws_security_group.source2.id]
    description     = "ingress from other security groups"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "egress for all ipv4"
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    ipv6_cidr_blocks = ["::/0"]
    description      = "egress for all ipv6"
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    prefix_list_ids = [aws_vpc_endpoint.test.prefix_list_id]
    description     = "egress for vpc endpoints"
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_rulesDropOnErrorInit(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test_ref0" {
  name   = "%[1]s-ref0"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test_ref1" {
  name   = "%[1]s-ref1"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  ingress {
    protocol  = "tcp"
    from_port = "80"
    to_port   = "80"
    security_groups = [
      aws_security_group.test_ref0.id,
      aws_security_group.test_ref1.id,
    ]
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_rulesDropOnErrorAddBadRule(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test_ref0" {
  name   = "%[1]s-ref0"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test_ref1" {
  name   = "%[1]s-ref1"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  ingress {
    protocol  = "tcp"
    from_port = "80"
    to_port   = "80"
    security_groups = [
      aws_security_group.test_ref0.id,
      aws_security_group.test_ref1.id,
      "sg-malformed", # non-existent rule to trigger API error
    ]
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_egressModeBlocks(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name = %[1]q

  tags = {
    Name = %[1]q
  }

  vpc_id = aws_vpc.test.id

  egress {
    cidr_blocks = [aws_vpc.test.cidr_block]
    from_port   = 0
    protocol    = "tcp"
    to_port     = 0
  }

  egress {
    cidr_blocks = [aws_vpc.test.cidr_block]
    from_port   = 0
    protocol    = "udp"
    to_port     = 0
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_egressModeNoBlocks(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name = %[1]q

  tags = {
    Name = %[1]q
  }

  vpc_id = aws_vpc.test.id
}
`, rName)
}

func testAccVPCSecurityGroupConfig_egressModeZeroed(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name = %[1]q

  tags = {
    Name = %[1]q
  }

  egress = []

  vpc_id = aws_vpc.test.id
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ingressModeBlocks(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name = %[1]q

  tags = {
    Name = %[1]q
  }

  vpc_id = aws_vpc.test.id

  ingress {
    cidr_blocks = [aws_vpc.test.cidr_block]
    from_port   = 0
    protocol    = "tcp"
    to_port     = 0
  }

  ingress {
    cidr_blocks = [aws_vpc.test.cidr_block]
    from_port   = 0
    protocol    = "udp"
    to_port     = 0
  }
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ingressModeNoBlocks(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name = %[1]q

  tags = {
    Name = %[1]q
  }

  vpc_id = aws_vpc.test.id
}
`, rName)
}

func testAccVPCSecurityGroupConfig_ingressModeZeroed(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name = %[1]q

  tags = {
    Name = %[1]q
  }

  ingress = []

  vpc_id = aws_vpc.test.id
}
`, rName)
}
