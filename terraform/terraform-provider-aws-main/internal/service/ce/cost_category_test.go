package ce_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/costexplorer"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfce "github.com/hashicorp/terraform-provider-aws/internal/service/ce"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccCECostCategory_basic(t *testing.T) {
	var output costexplorer.CostCategory
	resourceName := "aws_ce_cost_category.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCostCategoryDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, costexplorer.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccCostCategoryConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCostCategoryExists(resourceName, &output),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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

func TestAccCECostCategory_disappears(t *testing.T) {
	var output costexplorer.CostCategory
	resourceName := "aws_ce_cost_category.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCostCategoryDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, costexplorer.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccCostCategoryConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCostCategoryExists(resourceName, &output),
					acctest.CheckResourceDisappears(acctest.Provider, tfce.ResourceCostCategory(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCECostCategory_complete(t *testing.T) {
	var output costexplorer.CostCategory
	resourceName := "aws_ce_cost_category.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCostCategoryDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, costexplorer.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccCostCategoryConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCostCategoryExists(resourceName, &output),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccCostCategoryConfig_operandAnd(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCostCategoryExists(resourceName, &output),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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

func TestAccCECostCategory_splitCharge(t *testing.T) {
	var output costexplorer.CostCategory
	resourceName := "aws_ce_cost_category.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCostCategoryDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, costexplorer.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccCostCategoryConfig_splitCharges(rName, "PROPORTIONAL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCostCategoryExists(resourceName, &output),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccCostCategoryConfig_splitCharges(rName, "EVEN"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCostCategoryExists(resourceName, &output),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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

func TestAccCECostCategory_tags(t *testing.T) {
	var output costexplorer.CostCategory
	resourceName := "aws_ce_cost_category.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCostCategoryDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, costexplorer.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccCostCategoryConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCostCategoryExists(resourceName, &output),
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
				Config: testAccCostCategoryConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCostCategoryExists(resourceName, &output),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccCostCategoryConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCostCategoryExists(resourceName, &output),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckCostCategoryExists(n string, v *costexplorer.CostCategory) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CE Cost Category ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).CEConn

		output, err := tfce.FindCostCategoryByARN(context.Background(), conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckCostCategoryDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).CEConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ce_cost_category" {
			continue
		}

		_, err := tfce.FindCostCategoryByARN(context.Background(), conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("CE Cost Category %s still exists", rs.Primary.ID)
	}

	return nil

}

func testAccCostCategoryConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_ce_cost_category" "test" {
  name         = %[1]q
  rule_version = "CostCategoryExpression.v1"
  rule {
    value = "production"
    rule {
      dimension {
        key           = "LINKED_ACCOUNT_NAME"
        values        = ["-prod"]
        match_options = ["ENDS_WITH"]
      }
    }
    type = "REGULAR"
  }
  rule {
    value = "staging"
    rule {
      dimension {
        key           = "LINKED_ACCOUNT_NAME"
        values        = ["-stg"]
        match_options = ["ENDS_WITH"]
      }
    }
    type = "REGULAR"
  }
  rule {
    value = "testing"
    rule {
      dimension {
        key           = "LINKED_ACCOUNT_NAME"
        values        = ["-dev"]
        match_options = ["ENDS_WITH"]
      }
    }
    type = "REGULAR"
  }
}
`, rName)
}

func testAccCostCategoryConfig_operandAnd(rName string) string {
	return fmt.Sprintf(`
resource "aws_ce_cost_category" "test" {
  name         = %[1]q
  rule_version = "CostCategoryExpression.v1"
  rule {
    value = "production"
    rule {
      and {
        dimension {
          key           = "LINKED_ACCOUNT_NAME"
          values        = ["-prod"]
          match_options = ["ENDS_WITH"]
        }
      }
      and {
        dimension {
          key           = "LINKED_ACCOUNT_NAME"
          values        = ["-stg"]
          match_options = ["ENDS_WITH"]
        }
      }
      and {
        dimension {
          key           = "LINKED_ACCOUNT_NAME"
          values        = ["-dev"]
          match_options = ["ENDS_WITH"]
        }
      }
    }
    type = "REGULAR"
  }
}
`, rName)
}

func testAccCostCategoryConfig_splitCharges(rName, method string) string {
	return fmt.Sprintf(`
resource "aws_ce_cost_category" "test1" {
  name         = "%[1]s-1"
  rule_version = "CostCategoryExpression.v1"

  rule {
    value = "production"

    rule {
      dimension {
        key           = "LINKED_ACCOUNT_NAME"
        values        = ["-prod"]
        match_options = ["ENDS_WITH"]
      }
    }

    type = "REGULAR"
  }

  rule {
    value = "staging"

    rule {
      dimension {
        key           = "LINKED_ACCOUNT_NAME"
        values        = ["-stg"]
        match_options = ["ENDS_WITH"]
      }
    }

    type = "REGULAR"
  }

  rule {
    value = "testing"

    rule {
      dimension {
        key           = "LINKED_ACCOUNT_NAME"
        values        = ["-dev"]
        match_options = ["ENDS_WITH"]
      }
    }

    type = "REGULAR"
  }
}

resource "aws_ce_cost_category" "test2" {
  name         = "%[1]s-2"
  rule_version = "CostCategoryExpression.v1"

  rule {
    value = "production"

    rule {
      and {
        dimension {
          key           = "LINKED_ACCOUNT_NAME"
          values        = ["-prod"]
          match_options = ["ENDS_WITH"]
        }
      }

      and {
        dimension {
          key           = "LINKED_ACCOUNT_NAME"
          values        = ["-stg"]
          match_options = ["ENDS_WITH"]
        }
      }

      and {
        dimension {
          key           = "LINKED_ACCOUNT_NAME"
          values        = ["-dev"]
          match_options = ["ENDS_WITH"]
        }
      }
    }

    type = "REGULAR"
  }
}

resource "aws_ce_cost_category" "test" {
  name         = %[1]q
  rule_version = "CostCategoryExpression.v1"

  rule {
    value = "production"
    rule {
      dimension {
        key           = "LINKED_ACCOUNT_NAME"
        values        = ["-prod"]
        match_options = ["ENDS_WITH"]
      }
    }

    type = "REGULAR"
  }

  split_charge_rule {
    method  = %[2]q
    source  = aws_ce_cost_category.test1.id
    targets = [aws_ce_cost_category.test2.id]
  }
}
`, rName, method)
}

func testAccCostCategoryConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_ce_cost_category" "test" {
  name         = %[1]q
  rule_version = "CostCategoryExpression.v1"

  rule {
    value = "production"

    rule {
      dimension {
        key           = "LINKED_ACCOUNT_NAME"
        values        = ["-prod"]
        match_options = ["ENDS_WITH"]
      }
    }

    type = "REGULAR"
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccCostCategoryConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_ce_cost_category" "test" {
  name         = %[1]q
  rule_version = "CostCategoryExpression.v1"

  rule {
    value = "production"

    rule {
      dimension {
        key           = "LINKED_ACCOUNT_NAME"
        values        = ["-prod"]
        match_options = ["ENDS_WITH"]
      }
    }

    type = "REGULAR"
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
