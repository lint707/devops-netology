package waf_test

import (
	"fmt"
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

func TestAccWAFByteMatchSet_basic(t *testing.T) {
	var v waf.ByteMatchSet
	byteMatchSet := fmt.Sprintf("byteMatchSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_waf_byte_match_set.byte_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, waf.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckByteMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccByteMatchSetConfig_basic(byteMatchSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckByteMatchSetExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", byteMatchSet),
					resource.TestCheckResourceAttr(resourceName, "byte_match_tuples.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer1",
						"text_transformation":   "NONE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer2",
						"text_transformation":   "NONE",
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

func TestAccWAFByteMatchSet_changeNameForceNew(t *testing.T) {
	var before, after waf.ByteMatchSet
	byteMatchSet := fmt.Sprintf("byteMatchSet-%s", sdkacctest.RandString(5))
	byteMatchSetNewName := fmt.Sprintf("byteMatchSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_waf_byte_match_set.byte_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, waf.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckByteMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccByteMatchSetConfig_basic(byteMatchSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckByteMatchSetExists(resourceName, &before),
					resource.TestCheckResourceAttr(resourceName, "name", byteMatchSet),
					resource.TestCheckResourceAttr(resourceName, "byte_match_tuples.#", "2"),
				),
			},
			{
				Config: testAccByteMatchSetConfig_changeName(byteMatchSetNewName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckByteMatchSetExists(resourceName, &after),
					resource.TestCheckResourceAttr(resourceName, "name", byteMatchSetNewName),
					resource.TestCheckResourceAttr(resourceName, "byte_match_tuples.#", "2"),
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

func TestAccWAFByteMatchSet_changeTuples(t *testing.T) {
	var before, after waf.ByteMatchSet
	byteMatchSetName := fmt.Sprintf("byteMatchSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_waf_byte_match_set.byte_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, waf.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckByteMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccByteMatchSetConfig_basic(byteMatchSetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckByteMatchSetExists(resourceName, &before),
					resource.TestCheckResourceAttr(resourceName, "name", byteMatchSetName),
					resource.TestCheckResourceAttr(resourceName, "byte_match_tuples.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer1",
						"text_transformation":   "NONE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer2",
						"text_transformation":   "NONE",
					}),
				),
			},
			{
				Config: testAccByteMatchSetConfig_changeTuples(byteMatchSetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckByteMatchSetExists(resourceName, &after),
					resource.TestCheckResourceAttr(resourceName, "name", byteMatchSetName),
					resource.TestCheckResourceAttr(resourceName, "byte_match_tuples.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer1",
						"text_transformation":   "NONE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "",
						"field_to_match.0.type": "METHOD",
						"positional_constraint": "CONTAINS_WORD",
						"target_string":         "blah",
						"text_transformation":   "URL_DECODE",
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

func TestAccWAFByteMatchSet_noTuples(t *testing.T) {
	var byteSet waf.ByteMatchSet
	byteMatchSetName := fmt.Sprintf("byteMatchSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_waf_byte_match_set.byte_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, waf.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckByteMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccByteMatchSetConfig_noTuples(byteMatchSetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckByteMatchSetExists(resourceName, &byteSet),
					resource.TestCheckResourceAttr(resourceName, "name", byteMatchSetName),
					resource.TestCheckResourceAttr(resourceName, "byte_match_tuples.#", "0"),
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

func TestAccWAFByteMatchSet_disappears(t *testing.T) {
	var v waf.ByteMatchSet
	byteMatchSet := fmt.Sprintf("byteMatchSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_waf_byte_match_set.byte_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, waf.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckByteMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccByteMatchSetConfig_basic(byteMatchSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckByteMatchSetExists(resourceName, &v),
					testAccCheckByteMatchSetDisappears(&v),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckByteMatchSetDisappears(v *waf.ByteMatchSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFConn

		wr := tfwaf.NewRetryer(conn)
		_, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
			req := &waf.UpdateByteMatchSetInput{
				ChangeToken:    token,
				ByteMatchSetId: v.ByteMatchSetId,
			}

			for _, ByteMatchTuple := range v.ByteMatchTuples {
				ByteMatchUpdate := &waf.ByteMatchSetUpdate{
					Action: aws.String("DELETE"),
					ByteMatchTuple: &waf.ByteMatchTuple{
						FieldToMatch:         ByteMatchTuple.FieldToMatch,
						PositionalConstraint: ByteMatchTuple.PositionalConstraint,
						TargetString:         ByteMatchTuple.TargetString,
						TextTransformation:   ByteMatchTuple.TextTransformation,
					},
				}
				req.Updates = append(req.Updates, ByteMatchUpdate)
			}

			return conn.UpdateByteMatchSet(req)
		})
		if err != nil {
			return fmt.Errorf("Error updating ByteMatchSet: %s", err)
		}

		_, err = wr.RetryWithToken(func(token *string) (interface{}, error) {
			opts := &waf.DeleteByteMatchSetInput{
				ChangeToken:    token,
				ByteMatchSetId: v.ByteMatchSetId,
			}
			return conn.DeleteByteMatchSet(opts)
		})
		if err != nil {
			return fmt.Errorf("Error deleting ByteMatchSet: %s", err)
		}

		return nil
	}
}

func testAccCheckByteMatchSetExists(n string, v *waf.ByteMatchSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No WAF ByteMatchSet ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFConn
		resp, err := conn.GetByteMatchSet(&waf.GetByteMatchSetInput{
			ByteMatchSetId: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		if *resp.ByteMatchSet.ByteMatchSetId == rs.Primary.ID {
			*v = *resp.ByteMatchSet
			return nil
		}

		return fmt.Errorf("WAF ByteMatchSet (%s) not found", rs.Primary.ID)
	}
}

func testAccCheckByteMatchSetDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_waf_byte_match_set" {
			continue
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFConn
		resp, err := conn.GetByteMatchSet(
			&waf.GetByteMatchSetInput{
				ByteMatchSetId: aws.String(rs.Primary.ID),
			})

		if err == nil {
			if *resp.ByteMatchSet.ByteMatchSetId == rs.Primary.ID {
				return fmt.Errorf("WAF ByteMatchSet %s still exists", rs.Primary.ID)
			}
		}

		if tfawserr.ErrCodeEquals(err, waf.ErrCodeNonexistentItemException) {
			continue
		}

		return err
	}

	return nil
}

func testAccByteMatchSetConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_byte_match_set" "byte_set" {
  name = "%s"

  byte_match_tuples {
    text_transformation   = "NONE"
    target_string         = "badrefer1"
    positional_constraint = "CONTAINS"

    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }

  byte_match_tuples {
    text_transformation   = "NONE"
    target_string         = "badrefer2"
    positional_constraint = "CONTAINS"

    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }
}
`, name)
}

func testAccByteMatchSetConfig_changeName(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_byte_match_set" "byte_set" {
  name = "%s"

  byte_match_tuples {
    text_transformation   = "NONE"
    target_string         = "badrefer1"
    positional_constraint = "CONTAINS"

    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }

  byte_match_tuples {
    text_transformation   = "NONE"
    target_string         = "badrefer2"
    positional_constraint = "CONTAINS"

    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }
}
`, name)
}

func testAccByteMatchSetConfig_changeTuples(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_byte_match_set" "byte_set" {
  name = "%s"

  byte_match_tuples {
    text_transformation   = "NONE"
    target_string         = "badrefer1"
    positional_constraint = "CONTAINS"

    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }

  byte_match_tuples {
    text_transformation   = "URL_DECODE"
    target_string         = "blah"
    positional_constraint = "CONTAINS_WORD"

    field_to_match {
      type = "METHOD"
      # data field omitted as the type is neither "HEADER" nor "SINGLE_QUERY_ARG"
    }
  }
}
`, name)
}

func testAccByteMatchSetConfig_noTuples(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_byte_match_set" "byte_set" {
  name = "%s"
}
`, name)
}
