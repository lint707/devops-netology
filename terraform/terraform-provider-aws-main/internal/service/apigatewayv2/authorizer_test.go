package apigatewayv2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfapigatewayv2 "github.com/hashicorp/terraform-provider-aws/internal/service/apigatewayv2"
)

func TestAccAPIGatewayV2Authorizer_basic(t *testing.T) {
	var apiId string
	var v apigatewayv2.GetAuthorizerOutput
	resourceName := "aws_apigatewayv2_authorizer.test"
	lambdaResourceName := "aws_lambda_function.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigatewayv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttr(resourceName, "authorizer_credentials_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_payload_format_version", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_type", "REQUEST"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttr(resourceName, "enable_simple_responses", "false"),
					resource.TestCheckResourceAttr(resourceName, "identity_sources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccAuthorizerImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAPIGatewayV2Authorizer_disappears(t *testing.T) {
	var apiId string
	var v apigatewayv2.GetAuthorizerOutput
	resourceName := "aws_apigatewayv2_authorizer.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigatewayv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					acctest.CheckResourceDisappears(acctest.Provider, tfapigatewayv2.ResourceAuthorizer(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAPIGatewayV2Authorizer_credentials(t *testing.T) {
	var apiId string
	var v apigatewayv2.GetAuthorizerOutput
	resourceName := "aws_apigatewayv2_authorizer.test"
	iamRoleResourceName := "aws_iam_role.test"
	lambdaResourceName := "aws_lambda_function.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigatewayv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_credentials(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_credentials_arn", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_payload_format_version", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_type", "REQUEST"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttr(resourceName, "enable_simple_responses", "false"),
					resource.TestCheckResourceAttr(resourceName, "identity_sources.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "route.request.header.Auth"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccAuthorizerImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAuthorizerConfig_credentialsUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_credentials_arn", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_type", "REQUEST"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_payload_format_version", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttr(resourceName, "enable_simple_responses", "false"),
					resource.TestCheckResourceAttr(resourceName, "identity_sources.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "route.request.header.Auth"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "route.request.querystring.Name"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-updated", rName)),
				),
			},
			{
				Config: testAccAuthorizerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttr(resourceName, "authorizer_credentials_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_type", "REQUEST"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_payload_format_version", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttr(resourceName, "enable_simple_responses", "false"),
					resource.TestCheckResourceAttr(resourceName, "identity_sources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccAPIGatewayV2Authorizer_jwt(t *testing.T) {
	var apiId string
	var v apigatewayv2.GetAuthorizerOutput
	resourceName := "aws_apigatewayv2_authorizer.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigatewayv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_jwt(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttr(resourceName, "authorizer_credentials_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_payload_format_version", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_type", "JWT"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_uri", ""),
					resource.TestCheckResourceAttr(resourceName, "enable_simple_responses", "false"),
					resource.TestCheckResourceAttr(resourceName, "identity_sources.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "$request.header.Authorization"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.0.audience.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "jwt_configuration.0.audience.*", "test"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccAuthorizerImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAuthorizerConfig_jwtUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttr(resourceName, "authorizer_credentials_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_payload_format_version", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_type", "JWT"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_uri", ""),
					resource.TestCheckResourceAttr(resourceName, "enable_simple_responses", "false"),
					resource.TestCheckResourceAttr(resourceName, "identity_sources.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "$request.header.Authorization"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.0.audience.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "jwt_configuration.0.audience.*", "test"),
					resource.TestCheckTypeSetElemAttr(resourceName, "jwt_configuration.0.audience.*", "testing"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccAPIGatewayV2Authorizer_HTTPAPILambdaRequestAuthorizer_initialMissingCacheTTL(t *testing.T) {
	var apiId string
	var v apigatewayv2.GetAuthorizerOutput
	resourceName := "aws_apigatewayv2_authorizer.test"
	lambdaResourceName := "aws_lambda_function.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigatewayv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_httpAPILambdaRequest(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttr(resourceName, "authorizer_credentials_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_payload_format_version", "2.0"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "300"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_type", "REQUEST"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttr(resourceName, "enable_simple_responses", "true"),
					resource.TestCheckResourceAttr(resourceName, "identity_sources.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "$request.header.Auth"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccAuthorizerImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAuthorizerConfig_httpAPILambdaRequestUpdated(rName, 3600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttr(resourceName, "authorizer_credentials_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_payload_format_version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_type", "REQUEST"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttr(resourceName, "enable_simple_responses", "false"),
					resource.TestCheckResourceAttr(resourceName, "identity_sources.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "$request.querystring.User"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "$context.routeKey"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccAuthorizerConfig_httpAPILambdaRequestUpdated(rName, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttr(resourceName, "authorizer_credentials_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_payload_format_version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_type", "REQUEST"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttr(resourceName, "enable_simple_responses", "false"),
					resource.TestCheckResourceAttr(resourceName, "identity_sources.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "$request.querystring.User"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "$context.routeKey"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccAPIGatewayV2Authorizer_HTTPAPILambdaRequestAuthorizer_initialZeroCacheTTL(t *testing.T) {
	var apiId string
	var v apigatewayv2.GetAuthorizerOutput
	resourceName := "aws_apigatewayv2_authorizer.test"
	lambdaResourceName := "aws_lambda_function.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigatewayv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerConfig_httpAPILambdaRequestUpdated(rName, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttr(resourceName, "authorizer_credentials_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_payload_format_version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_type", "REQUEST"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttr(resourceName, "enable_simple_responses", "false"),
					resource.TestCheckResourceAttr(resourceName, "identity_sources.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "$request.querystring.User"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "$context.routeKey"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccAuthorizerImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAuthorizerConfig_httpAPILambdaRequestUpdated(rName, 600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizerExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttr(resourceName, "authorizer_credentials_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "authorizer_payload_format_version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_result_ttl_in_seconds", "600"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_type", "REQUEST"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", lambdaResourceName, "invoke_arn"),
					resource.TestCheckResourceAttr(resourceName, "enable_simple_responses", "false"),
					resource.TestCheckResourceAttr(resourceName, "identity_sources.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "$request.querystring.User"),
					resource.TestCheckTypeSetElemAttr(resourceName, "identity_sources.*", "$context.routeKey"),
					resource.TestCheckResourceAttr(resourceName, "jwt_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckAuthorizerDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayV2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_apigatewayv2_authorizer" {
			continue
		}

		_, err := conn.GetAuthorizer(&apigatewayv2.GetAuthorizerInput{
			ApiId:        aws.String(rs.Primary.Attributes["api_id"]),
			AuthorizerId: aws.String(rs.Primary.ID),
		})
		if tfawserr.ErrCodeEquals(err, apigatewayv2.ErrCodeNotFoundException) {
			continue
		}
		if err != nil {
			return err
		}

		return fmt.Errorf("API Gateway v2 authorizer %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckAuthorizerExists(n string, vApiId *string, v *apigatewayv2.GetAuthorizerOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No API Gateway v2 authorizer ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayV2Conn

		apiId := aws.String(rs.Primary.Attributes["api_id"])
		resp, err := conn.GetAuthorizer(&apigatewayv2.GetAuthorizerInput{
			ApiId:        apiId,
			AuthorizerId: aws.String(rs.Primary.ID),
		})
		if err != nil {
			return err
		}

		*vApiId = *apiId
		*v = *resp

		return nil
	}
}

func testAccAuthorizerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["api_id"], rs.Primary.ID), nil
	}
}

func testAccAuthorizerConfig_apiWebSocket(rName string) string {
	return fmt.Sprintf(`
resource "aws_apigatewayv2_api" "test" {
  name                       = %[1]q
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}
`, rName)
}

func testAccAuthorizerConfig_apiHTTP(rName string) string {
	return fmt.Sprintf(`
resource "aws_apigatewayv2_api" "test" {
  name          = %[1]q
  protocol_type = "HTTP"
}
`, rName)
}

func testAccAuthorizerConfig_baseLambda(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigLambdaBase(rName, rName, rName), fmt.Sprintf(`
resource "aws_lambda_function" "test" {
  filename      = "test-fixtures/lambdatest.zip"
  function_name = %[1]q
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "index.handler"
  runtime       = "nodejs14.x"
}

resource "aws_iam_role" "test" {
  name = "%[1]s_auth_invocation_role"
  path = "/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Action": "sts:AssumeRole",
    "Principal": {"Service": "apigateway.amazonaws.com"},
    "Effect": "Allow"
  }]
}
EOF
}
`, rName))
}

func testAccAuthorizerConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccAuthorizerConfig_apiWebSocket(rName),
		testAccAuthorizerConfig_baseLambda(rName),
		fmt.Sprintf(`
resource "aws_apigatewayv2_authorizer" "test" {
  api_id          = aws_apigatewayv2_api.test.id
  authorizer_type = "REQUEST"
  authorizer_uri  = aws_lambda_function.test.invoke_arn
  name            = %[1]q
}
`, rName))
}

func testAccAuthorizerConfig_credentials(rName string) string {
	return acctest.ConfigCompose(
		testAccAuthorizerConfig_apiWebSocket(rName),
		testAccAuthorizerConfig_baseLambda(rName),
		fmt.Sprintf(`
resource "aws_apigatewayv2_authorizer" "test" {
  api_id           = aws_apigatewayv2_api.test.id
  authorizer_type  = "REQUEST"
  authorizer_uri   = aws_lambda_function.test.invoke_arn
  identity_sources = ["route.request.header.Auth"]
  name             = %[1]q

  authorizer_credentials_arn = aws_iam_role.test.arn
}
`, rName))
}

func testAccAuthorizerConfig_credentialsUpdated(rName string) string {
	return acctest.ConfigCompose(
		testAccAuthorizerConfig_apiWebSocket(rName),
		testAccAuthorizerConfig_baseLambda(rName),
		fmt.Sprintf(`
resource "aws_apigatewayv2_authorizer" "test" {
  api_id           = aws_apigatewayv2_api.test.id
  authorizer_type  = "REQUEST"
  authorizer_uri   = aws_lambda_function.test.invoke_arn
  identity_sources = ["route.request.header.Auth", "route.request.querystring.Name"]
  name             = "%[1]s-updated"

  authorizer_credentials_arn = aws_iam_role.test.arn
}
`, rName))
}

func testAccAuthorizerConfig_jwt(rName string) string {
	return acctest.ConfigCompose(
		testAccAuthorizerConfig_apiHTTP(rName),
		testAccAuthorizerConfig_baseLambda(rName),
		fmt.Sprintf(`
resource "aws_cognito_user_pool" "test" {
  name = %[1]q
}

resource "aws_apigatewayv2_authorizer" "test" {
  api_id           = aws_apigatewayv2_api.test.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = %[1]q

  jwt_configuration {
    audience = ["test"]
    issuer   = "https://${aws_cognito_user_pool.test.endpoint}"
  }
}
`, rName))
}

func testAccAuthorizerConfig_jwtUpdated(rName string) string {
	return acctest.ConfigCompose(
		testAccAuthorizerConfig_apiHTTP(rName),
		testAccAuthorizerConfig_baseLambda(rName),
		fmt.Sprintf(`
resource "aws_cognito_user_pool" "test" {
  name = %[1]q
}

resource "aws_apigatewayv2_authorizer" "test" {
  api_id           = aws_apigatewayv2_api.test.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = %[1]q

  jwt_configuration {
    audience = ["test", "testing"]
    issuer   = "https://${aws_cognito_user_pool.test.endpoint}"
  }
}
`, rName))
}

func testAccAuthorizerConfig_httpAPILambdaRequest(rName string) string {
	return acctest.ConfigCompose(
		testAccAuthorizerConfig_apiHTTP(rName),
		testAccAuthorizerConfig_baseLambda(rName),
		fmt.Sprintf(`
resource "aws_apigatewayv2_authorizer" "test" {
  api_id                            = aws_apigatewayv2_api.test.id
  authorizer_payload_format_version = "2.0"
  authorizer_type                   = "REQUEST"
  authorizer_uri                    = aws_lambda_function.test.invoke_arn
  enable_simple_responses           = true
  identity_sources                  = ["$request.header.Auth"]
  name                              = %[1]q
}
`, rName))
}

func testAccAuthorizerConfig_httpAPILambdaRequestUpdated(rName string, authorizerResultTtl int) string {
	return acctest.ConfigCompose(
		testAccAuthorizerConfig_apiHTTP(rName),
		testAccAuthorizerConfig_baseLambda(rName),
		fmt.Sprintf(`
resource "aws_apigatewayv2_authorizer" "test" {
  api_id                            = aws_apigatewayv2_api.test.id
  authorizer_payload_format_version = "1.0"
  authorizer_result_ttl_in_seconds  = %[2]d
  authorizer_type                   = "REQUEST"
  authorizer_uri                    = aws_lambda_function.test.invoke_arn
  enable_simple_responses           = false
  identity_sources                  = ["$request.querystring.User", "$context.routeKey"]
  name                              = %[1]q
}
`, rName, authorizerResultTtl))
}
