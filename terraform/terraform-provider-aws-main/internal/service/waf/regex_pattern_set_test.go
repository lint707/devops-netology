package waf_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfwaf "github.com/hashicorp/terraform-provider-aws/internal/service/waf"
)

// Serialized acceptance tests due to WAF account limits
// https://docs.aws.amazon.com/waf/latest/developerguide/limits.html
func TestAccWAFRegexPatternSet_serial(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic":          testAccRegexPatternSet_basic,
		"changePatterns": testAccRegexPatternSet_changePatterns,
		"noPatterns":     testAccRegexPatternSet_noPatterns,
		"disappears":     testAccRegexPatternSet_disappears,
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccRegexPatternSet_basic(t *testing.T) {
	var v waf.RegexPatternSet
	patternSetName := fmt.Sprintf("tfacc-%s", sdkacctest.RandString(5))
	resourceName := "aws_waf_regex_pattern_set.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, waf.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRegexPatternSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegexPatternSetConfig_basic(patternSetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegexPatternSetExists(resourceName, &v),
					acctest.MatchResourceAttrGlobalARN(resourceName, "arn", "waf", regexp.MustCompile(`regexpatternset/.+`)),
					resource.TestCheckResourceAttr(resourceName, "name", patternSetName),
					resource.TestCheckResourceAttr(resourceName, "regex_pattern_strings.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "regex_pattern_strings.*", "one"),
					resource.TestCheckTypeSetElemAttr(resourceName, "regex_pattern_strings.*", "two"),
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

func testAccRegexPatternSet_changePatterns(t *testing.T) {
	var before, after waf.RegexPatternSet
	patternSetName := fmt.Sprintf("tfacc-%s", sdkacctest.RandString(5))
	resourceName := "aws_waf_regex_pattern_set.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, waf.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRegexPatternSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegexPatternSetConfig_basic(patternSetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRegexPatternSetExists(resourceName, &before),
					resource.TestCheckResourceAttr(resourceName, "name", patternSetName),
					resource.TestCheckResourceAttr(resourceName, "regex_pattern_strings.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "regex_pattern_strings.*", "one"),
					resource.TestCheckTypeSetElemAttr(resourceName, "regex_pattern_strings.*", "two"),
				),
			},
			{
				Config: testAccRegexPatternSetConfig_changes(patternSetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRegexPatternSetExists(resourceName, &after),
					resource.TestCheckResourceAttr(resourceName, "name", patternSetName),
					resource.TestCheckResourceAttr(resourceName, "regex_pattern_strings.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "regex_pattern_strings.*", "two"),
					resource.TestCheckTypeSetElemAttr(resourceName, "regex_pattern_strings.*", "three"),
					resource.TestCheckTypeSetElemAttr(resourceName, "regex_pattern_strings.*", "four"),
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

func testAccRegexPatternSet_noPatterns(t *testing.T) {
	var patternSet waf.RegexPatternSet
	patternSetName := fmt.Sprintf("tfacc-%s", sdkacctest.RandString(5))
	resourceName := "aws_waf_regex_pattern_set.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, waf.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRegexPatternSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegexPatternSetConfig_nos(patternSetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRegexPatternSetExists(resourceName, &patternSet),
					resource.TestCheckResourceAttr(resourceName, "name", patternSetName),
					resource.TestCheckResourceAttr(resourceName, "regex_pattern_strings.#", "0"),
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

func testAccRegexPatternSet_disappears(t *testing.T) {
	var v waf.RegexPatternSet
	patternSetName := fmt.Sprintf("tfacc-%s", sdkacctest.RandString(5))
	resourceName := "aws_waf_regex_pattern_set.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, waf.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRegexPatternSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegexPatternSetConfig_basic(patternSetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegexPatternSetExists(resourceName, &v),
					testAccCheckRegexPatternSetDisappears(&v),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRegexPatternSetDisappears(set *waf.RegexPatternSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFConn

		wr := tfwaf.NewRetryer(conn)
		_, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
			req := &waf.UpdateRegexPatternSetInput{
				ChangeToken:       token,
				RegexPatternSetId: set.RegexPatternSetId,
			}

			for _, pattern := range set.RegexPatternStrings {
				update := &waf.RegexPatternSetUpdate{
					Action:             aws.String("DELETE"),
					RegexPatternString: pattern,
				}
				req.Updates = append(req.Updates, update)
			}

			return conn.UpdateRegexPatternSet(req)
		})
		if err != nil {
			return fmt.Errorf("Failed updating WAF Regex Pattern Set: %s", err)
		}

		_, err = wr.RetryWithToken(func(token *string) (interface{}, error) {
			opts := &waf.DeleteRegexPatternSetInput{
				ChangeToken:       token,
				RegexPatternSetId: set.RegexPatternSetId,
			}
			return conn.DeleteRegexPatternSet(opts)
		})
		if err != nil {
			return fmt.Errorf("Failed deleting WAF Regex Pattern Set: %s", err)
		}

		return nil
	}
}

func testAccCheckRegexPatternSetExists(n string, v *waf.RegexPatternSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No WAF Regex Pattern Set ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFConn
		resp, err := conn.GetRegexPatternSet(&waf.GetRegexPatternSetInput{
			RegexPatternSetId: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		if *resp.RegexPatternSet.RegexPatternSetId == rs.Primary.ID {
			*v = *resp.RegexPatternSet
			return nil
		}

		return fmt.Errorf("WAF Regex Pattern Set (%s) not found", rs.Primary.ID)
	}
}

func testAccCheckRegexPatternSetDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_waf_regex_pattern_set" {
			continue
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFConn
		resp, err := conn.GetRegexPatternSet(&waf.GetRegexPatternSetInput{
			RegexPatternSetId: aws.String(rs.Primary.ID),
		})

		if err == nil {
			if *resp.RegexPatternSet.RegexPatternSetId == rs.Primary.ID {
				return fmt.Errorf("WAF Regex Pattern Set %s still exists", rs.Primary.ID)
			}
		}

		// Return nil if the Regex Pattern Set is already destroyed
		if tfawserr.ErrCodeEquals(err, waf.ErrCodeNonexistentItemException) {
			return nil
		}

		return err
	}

	return nil
}

func testAccRegexPatternSetConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_regex_pattern_set" "test" {
  name                  = "%s"
  regex_pattern_strings = ["one", "two"]
}
`, name)
}

func testAccRegexPatternSetConfig_changes(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_regex_pattern_set" "test" {
  name                  = "%s"
  regex_pattern_strings = ["two", "three", "four"]
}
`, name)
}

func testAccRegexPatternSetConfig_nos(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_regex_pattern_set" "test" {
  name = "%s"
}
`, name)
}
