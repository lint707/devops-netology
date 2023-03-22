package apigateway_test

import (
	"fmt"
	"regexp"
	"strconv"
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

func TestAccAPIGatewayAuthorizer_basic(t *testing.T) {
	var conf apigateway.Authorizer
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_authorizer.test"
	lambdaResourceName := "aws_lambda_function.test"
	roleResourceName := "aws_iam_role.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_lambda(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &conf),
					acctest.MatchResourceAttrRegionalARNNoAccount(resourceName, "arn", "apigateway", regexp.MustCompile(`/restapis/.+/authorizers/.+`)),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttr(resourceName, "identity_source", "method.request.header.Authorization"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "TOKEN"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_credentials", roleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", strconv.Itoa(tfapigateway.DefaultAuthorizerTTL)),
					resource.TestCheckResourceAttr(resourceName, "identity_validation_expression", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAuthorizerImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				Config: testAccAuthorizerConfig_lambdaUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttr(resourceName, "identity_source", "method.request.header.Authorization"),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_modified"),
					resource.TestCheckResourceAttr(resourceName, "type", "TOKEN"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_credentials", roleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", strconv.Itoa(360)),
					resource.TestCheckResourceAttr(resourceName, "identity_validation_expression", ".*"),
				),
			},
		},
	})
}

func TestAccAPIGatewayAuthorizer_cognito(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_authorizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_cognito(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "COGNITO_USER_POOLS"),
					resource.TestCheckResourceAttr(resourceName, "provider_arns.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAuthorizerImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				Config: testAccAuthorizerConfig_cognitoUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "COGNITO_USER_POOLS"),
					resource.TestCheckResourceAttr(resourceName, "provider_arns.#", "3"),
				),
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/16613
func TestAccAPIGatewayAuthorizer_Cognito_authorizerCredentials(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_authorizer.test"
	iamRoleResourceName := "aws_iam_role.lambda"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_cognitoCredentials(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_credentials", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "COGNITO_USER_POOLS"),
					resource.TestCheckResourceAttr(resourceName, "provider_arns.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAuthorizerImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAPIGatewayAuthorizer_switchAuthType(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_authorizer.test"
	lambdaResourceName := "aws_lambda_function.test"
	roleResourceName := "aws_iam_role.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_lambda(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "TOKEN"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_credentials", roleResourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAuthorizerImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				Config: testAccAuthorizerConfig_cognito(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "COGNITO_USER_POOLS"),
					resource.TestCheckResourceAttr(resourceName, "provider_arns.#", "2"),
				),
			},
			{
				Config: testAccAuthorizerConfig_lambdaUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_modified"),
					resource.TestCheckResourceAttr(resourceName, "type", "TOKEN"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_credentials", roleResourceName, "arn"),
				),
			},
		},
	})
}

func TestAccAPIGatewayAuthorizer_switchAuthorizerTTL(t *testing.T) {
	var conf apigateway.Authorizer
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_authorizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_lambda(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", strconv.Itoa(tfapigateway.DefaultAuthorizerTTL)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAuthorizerImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				Config: testAccAuthorizerConfig_lambdaUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", strconv.Itoa(360)),
				),
			},
			{
				Config: testAccAuthorizerConfig_lambdaNoCache(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", strconv.Itoa(0)),
				),
			},
			{
				Config: testAccAuthorizerConfig_lambda(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", strconv.Itoa(tfapigateway.DefaultAuthorizerTTL)),
				),
			},
		},
	})
}

func TestAccAPIGatewayAuthorizer_authTypeValidation(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAuthorizerConfig_authTypeValidationDefaultToken(rName),
				ExpectError: regexp.MustCompile(`authorizer_uri must be set non-empty when authorizer type is TOKEN`),
			},
			{
				Config:      testAccAuthorizerConfig_authTypeValidationRequest(rName),
				ExpectError: regexp.MustCompile(`authorizer_uri must be set non-empty when authorizer type is REQUEST`),
			},
			{
				Config:      testAccAuthorizerConfig_authTypeValidationCognito(rName),
				ExpectError: regexp.MustCompile(`provider_arns must be set non-empty when authorizer type is COGNITO_USER_POOLS`),
			},
		},
	})
}

func TestAccAPIGatewayAuthorizer_Zero_ttl(t *testing.T) {
	var conf apigateway.Authorizer
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_authorizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_lambdaNoCache(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAuthorizerImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAPIGatewayAuthorizer_disappears(t *testing.T) {
	var conf apigateway.Authorizer
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_authorizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_lambda(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &conf),
					acctest.CheckResourceDisappears(acctest.Provider, tfapigateway.ResourceAuthorizer(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAuthorizerExists(n string, res *apigateway.Authorizer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No API Gateway Authorizer ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayConn

		req := &apigateway.GetAuthorizerInput{
			AuthorizerId: aws.String(rs.Primary.ID),
			RestApiId:    aws.String(rs.Primary.Attributes["rest_api_id"]),
		}
		describe, err := conn.GetAuthorizer(req)
		if err != nil {
			return err
		}

		*res = *describe

		return nil
	}
}

func testAccCheckAuthorizerDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_api_gateway_authorizer" {
			continue
		}

		req := &apigateway.GetAuthorizerInput{
			AuthorizerId: aws.String(rs.Primary.ID),
			RestApiId:    aws.String(rs.Primary.Attributes["rest_api_id"]),
		}
		_, err := conn.GetAuthorizer(req)

		if err == nil {
			return fmt.Errorf("API Gateway Authorizer still exists")
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

func testAccAuthorizerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["rest_api_id"], rs.Primary.ID), nil
	}
}

func testAccAuthorizerBaseConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_iam_role" "test" {
  name = %[1]q
  path = "/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "test" {
  name = %[1]q
  role = aws_iam_role.test.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "lambda:InvokeFunction",
      "Effect": "Allow",
      "Resource": "${aws_lambda_function.test.arn}"
    }
  ]
}
EOF
}

resource "aws_iam_role" "lambda" {
  name = "%[1]s-lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "test" {
  filename         = "test-fixtures/lambdatest.zip"
  source_code_hash = filebase64sha256("test-fixtures/lambdatest.zip")
  function_name    = %[1]q
  role             = aws_iam_role.lambda.arn
  handler          = "exports.example"
  runtime          = "nodejs12.x"
}
`, rName)
}

func testAccAuthorizerConfig_lambda(rName string) string {
	return testAccAuthorizerBaseConfig(rName) + fmt.Sprintf(`
resource "aws_api_gateway_authorizer" "test" {
  name                   = %[1]q
  rest_api_id            = aws_api_gateway_rest_api.test.id
  authorizer_uri         = aws_lambda_function.test.invoke_arn
  authorizer_credentials = aws_iam_role.test.arn
}
`, rName)
}

func testAccAuthorizerConfig_lambdaUpdate(rName string) string {
	return testAccAuthorizerBaseConfig(rName) + fmt.Sprintf(`
resource "aws_api_gateway_authorizer" "test" {
  name                             = "%[1]s_modified"
  rest_api_id                      = aws_api_gateway_rest_api.test.id
  authorizer_uri                   = aws_lambda_function.test.invoke_arn
  authorizer_credentials           = aws_iam_role.test.arn
  authorizer_result_ttl_in_seconds = 360
  identity_validation_expression   = ".*"
}
`, rName)
}

func testAccAuthorizerConfig_lambdaNoCache(rName string) string {
	return testAccAuthorizerBaseConfig(rName) + fmt.Sprintf(`
resource "aws_api_gateway_authorizer" "test" {
  name                             = "%[1]s_modified"
  rest_api_id                      = aws_api_gateway_rest_api.test.id
  authorizer_uri                   = aws_lambda_function.test.invoke_arn
  authorizer_credentials           = aws_iam_role.test.arn
  authorizer_result_ttl_in_seconds = 0
  identity_validation_expression   = ".*"
}
`, rName)
}

func testAccAuthorizerConfig_cognito(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_cognito_user_pool" "test" {
  count = 2
  name  = "%[1]s-${count.index}"
}

resource "aws_api_gateway_authorizer" "test" {
  name          = %[1]q
  type          = "COGNITO_USER_POOLS"
  rest_api_id   = aws_api_gateway_rest_api.test.id
  provider_arns = aws_cognito_user_pool.test[*].arn
}
`, rName)
}

func testAccAuthorizerConfig_cognitoUpdate(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_cognito_user_pool" "test" {
  count = 3
  name  = "%[1]s-${count.index}"
}

resource "aws_api_gateway_authorizer" "test" {
  name          = %[1]q
  type          = "COGNITO_USER_POOLS"
  rest_api_id   = aws_api_gateway_rest_api.test.id
  provider_arns = aws_cognito_user_pool.test[*].arn
}
`, rName)
}

func testAccAuthorizerConfig_cognitoCredentials(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "lambda" {
  name = "%[1]s-lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.${data.aws_partition.current.dns_suffix}"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_cognito_user_pool" "test" {
  count = 2
  name  = "%[1]s-${count.index}"
}

resource "aws_api_gateway_authorizer" "test" {
  authorizer_credentials = aws_iam_role.lambda.arn
  name                   = %[1]q
  type                   = "COGNITO_USER_POOLS"
  rest_api_id            = aws_api_gateway_rest_api.test.id
  provider_arns          = aws_cognito_user_pool.test[*].arn
}
`, rName)
}

func testAccAuthorizerConfig_authTypeValidationDefaultToken(rName string) string {
	return testAccAuthorizerBaseConfig(rName) + fmt.Sprintf(`
resource "aws_api_gateway_authorizer" "test" {
  name                   = "%s"
  rest_api_id            = aws_api_gateway_rest_api.test.id
  authorizer_credentials = aws_iam_role.test.arn
}
`, rName)
}

func testAccAuthorizerConfig_authTypeValidationRequest(rName string) string {
	return testAccAuthorizerBaseConfig(rName) + fmt.Sprintf(`
resource "aws_api_gateway_authorizer" "test" {
  name                   = "%s"
  type                   = "REQUEST"
  rest_api_id            = aws_api_gateway_rest_api.test.id
  authorizer_credentials = aws_iam_role.test.arn
}
`, rName)
}

func testAccAuthorizerConfig_authTypeValidationCognito(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_cognito_user_pool" "test" {
  count = 2
  name  = "%[1]s-${count.index}"
}

resource "aws_api_gateway_authorizer" "test" {
  name        = %[1]q
  type        = "COGNITO_USER_POOLS"
  rest_api_id = aws_api_gateway_rest_api.test.id
}
`, rName)
}
