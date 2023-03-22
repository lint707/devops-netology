package apigateway_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/apigateway"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfapigateway "github.com/hashicorp/terraform-provider-aws/internal/service/apigateway"
)

func TestAccAPIGatewayUsagePlanKey_basic(t *testing.T) {
	var conf apigateway.UsagePlanKey
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	apiGatewayApiKeyResourceName := "aws_api_gateway_api_key.test"
	apiGatewayUsagePlanResourceName := "aws_api_gateway_usage_plan.test"
	resourceName := "aws_api_gateway_usage_plan_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUsagePlanKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUsagePlanKeyConfig_typeAPI(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUsagePlanKeyExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "key_id", apiGatewayApiKeyResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "key_type", "API_KEY"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "usage_plan_id", apiGatewayUsagePlanResourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccCheckUsagePlanKeyImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAPIGatewayUsagePlanKey_disappears(t *testing.T) {
	var conf apigateway.UsagePlanKey
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_usage_plan_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUsagePlanKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUsagePlanKeyConfig_typeAPI(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUsagePlanKeyExists(resourceName, &conf),
					acctest.CheckResourceDisappears(acctest.Provider, tfapigateway.ResourceUsagePlanKey(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAPIGatewayUsagePlanKey_KeyID_concurrency(t *testing.T) {
	var conf apigateway.UsagePlanKey
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUsagePlanKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUsagePlanKeyConfig_idConcurrency(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUsagePlanKeyExists("aws_api_gateway_usage_plan_key.test.0", &conf),
					testAccCheckUsagePlanKeyExists("aws_api_gateway_usage_plan_key.test.1", &conf),
					testAccCheckUsagePlanKeyExists("aws_api_gateway_usage_plan_key.test.2", &conf),
					testAccCheckUsagePlanKeyExists("aws_api_gateway_usage_plan_key.test.3", &conf),
					testAccCheckUsagePlanKeyExists("aws_api_gateway_usage_plan_key.test.4", &conf),
					testAccCheckUsagePlanKeyExists("aws_api_gateway_usage_plan_key.test.5", &conf),
					testAccCheckUsagePlanKeyExists("aws_api_gateway_usage_plan_key.test.6", &conf),
					testAccCheckUsagePlanKeyExists("aws_api_gateway_usage_plan_key.test.7", &conf),
					testAccCheckUsagePlanKeyExists("aws_api_gateway_usage_plan_key.test.8", &conf),
					testAccCheckUsagePlanKeyExists("aws_api_gateway_usage_plan_key.test.9", &conf),
				),
			},
		},
	})
}

func testAccCheckUsagePlanKeyExists(n string, res *apigateway.UsagePlanKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No API Gateway Usage Plan Key ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayConn

		req := &apigateway.GetUsagePlanKeyInput{
			UsagePlanId: aws.String(rs.Primary.Attributes["usage_plan_id"]),
			KeyId:       aws.String(rs.Primary.Attributes["key_id"]),
		}
		up, err := conn.GetUsagePlanKey(req)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Reading API Gateway Usage Plan Key: %#v", up)

		if *up.Id != rs.Primary.ID {
			return fmt.Errorf("API Gateway Usage Plan Key not found")
		}

		*res = *up

		return nil
	}
}

func testAccCheckUsagePlanKeyDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_api_gateway_usage_plan_key" {
			continue
		}

		req := &apigateway.GetUsagePlanKeyInput{
			UsagePlanId: aws.String(rs.Primary.ID),
			KeyId:       aws.String(rs.Primary.Attributes["key_id"]),
		}
		describe, err := conn.GetUsagePlanKey(req)

		if err == nil {
			if describe.Id != nil && *describe.Id == rs.Primary.ID {
				return fmt.Errorf("API Gateway Usage Plan Key still exists")
			}
		}

		aws2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if aws2err.Code() != apigateway.ErrCodeNotFoundException {
			return err
		}

		return nil
	}

	return nil
}

func testAccCheckUsagePlanKeyImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["usage_plan_id"], rs.Primary.ID), nil
	}
}

func testAccUsagePlanKeyBaseConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = "%[1]s"
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  parent_id   = aws_api_gateway_rest_api.test.root_resource_id
  path_part   = "test"
}

resource "aws_api_gateway_method" "test" {
  rest_api_id   = aws_api_gateway_rest_api.test.id
  resource_id   = aws_api_gateway_resource.test.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "error" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  resource_id = aws_api_gateway_resource.test.id
  http_method = aws_api_gateway_method.test.http_method
  status_code = "400"
}

resource "aws_api_gateway_integration" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  resource_id = aws_api_gateway_resource.test.id
  http_method = aws_api_gateway_method.test.http_method

  type                    = "HTTP"
  uri                     = "https://www.google.de"
  integration_http_method = "GET"
}

resource "aws_api_gateway_integration_response" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  resource_id = aws_api_gateway_resource.test.id
  http_method = aws_api_gateway_integration.test.http_method
  status_code = aws_api_gateway_method_response.error.status_code
}

resource "aws_api_gateway_deployment" "test" {
  depends_on = [aws_api_gateway_integration_response.test]

  description = "This is a test"
  rest_api_id = aws_api_gateway_rest_api.test.id
  stage_name  = "test"
}
`, rName)
}

func testAccUsagePlanKeyConfig_typeAPI(rName string) string {
	return acctest.ConfigCompose(
		testAccUsagePlanKeyBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_api_gateway_api_key" "test" {
  name = %[1]q
}

resource "aws_api_gateway_usage_plan" "test" {
  name = %[1]q

  api_stages {
    api_id = aws_api_gateway_rest_api.test.id
    stage  = aws_api_gateway_deployment.test.stage_name
  }
}

resource "aws_api_gateway_usage_plan_key" "test" {
  key_id        = aws_api_gateway_api_key.test.id
  key_type      = "API_KEY"
  usage_plan_id = aws_api_gateway_usage_plan.test.id
}
`, rName))
}

func testAccUsagePlanKeyConfig_idConcurrency(rName string) string {
	return acctest.ConfigCompose(
		testAccUsagePlanKeyBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_api_gateway_api_key" "test" {
  count = 10

  name = "%[1]s-${count.index}"
}

resource "aws_api_gateway_usage_plan" "test" {
  name = %[1]q

  api_stages {
    api_id = aws_api_gateway_rest_api.test.id
    stage  = aws_api_gateway_deployment.test.stage_name
  }
}

resource "aws_api_gateway_usage_plan_key" "test" {
  count = 10

  key_id        = aws_api_gateway_api_key.test[count.index].id
  key_type      = "API_KEY"
  usage_plan_id = aws_api_gateway_usage_plan.test.id
}
`, rName))
}
