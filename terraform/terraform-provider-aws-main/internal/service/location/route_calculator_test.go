package location_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/locationservice"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tflocation "github.com/hashicorp/terraform-provider-aws/internal/service/location"
)

func TestAccLocationRouteCalculator_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_location_route_calculator.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, locationservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRouteCalculatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteCalculatorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteCalculatorExists(resourceName),
					acctest.CheckResourceAttrRegionalARN(resourceName, "calculator_arn", "geo", fmt.Sprintf("route-calculator/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "calculator_name", rName),
					acctest.CheckResourceAttrRFC3339(resourceName, "create_time"),
					resource.TestCheckResourceAttr(resourceName, "data_source", "Here"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					acctest.CheckResourceAttrRFC3339(resourceName, "update_time"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccLocationRouteCalculator_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_location_route_calculator.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, locationservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRouteCalculatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteCalculatorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteCalculatorExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tflocation.ResourceRouteCalculator(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccLocationRouteCalculator_description(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_location_route_calculator.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, locationservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRouteCalculatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteCalculatorConfig_description(rName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteCalculatorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRouteCalculatorConfig_description(rName, "description2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteCalculatorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				),
			},
		},
	})
}

func TestAccLocationRouteCalculator_tags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_location_route_calculator.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, locationservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRouteCalculatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteCalculatorConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteCalculatorExists(resourceName),
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
				Config: testAccRouteCalculatorConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteCalculatorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccRouteCalculatorConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteCalculatorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckRouteCalculatorDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).LocationConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_location_route_calculator" {
			continue
		}

		input := &locationservice.DescribeRouteCalculatorInput{
			CalculatorName: aws.String(rs.Primary.ID),
		}

		_, err := conn.DescribeRouteCalculator(input)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, locationservice.ErrCodeResourceNotFoundException) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Expected Location Service Route Calculator to be destroyed, %s found", rs.Primary.ID)
	}

	return nil
}

func testAccCheckRouteCalculatorExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Location Service Route Calculator is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).LocationConn
		_, err := conn.DescribeRouteCalculator(&locationservice.DescribeRouteCalculatorInput{
			CalculatorName: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("Error describing Location Service Route Calculator: %s", err.Error())
		}

		return nil
	}
}

func testAccRouteCalculatorConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_location_route_calculator" "test" {
  calculator_name = %[1]q
  data_source     = "Here"
}
`, rName)
}

func testAccRouteCalculatorConfig_description(rName, description string) string {
	return fmt.Sprintf(`
resource "aws_location_route_calculator" "test" {
  calculator_name = %[1]q
  data_source     = "Here"
  description     = %[2]q
}
`, rName, description)
}

func testAccRouteCalculatorConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_location_route_calculator" "test" {
  calculator_name = %[1]q
  data_source     = "Here"

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccRouteCalculatorConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_location_route_calculator" "test" {
  calculator_name = %[1]q
  data_source     = "Here"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
