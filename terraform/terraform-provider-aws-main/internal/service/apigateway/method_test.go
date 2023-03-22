package apigateway_test

import (
	"fmt"
	"regexp"
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

func TestAccAPIGatewayMethod_basic(t *testing.T) {
	var conf apigateway.Method
	rInt := sdkacctest.RandInt()
	resourceName := "aws_api_gateway_method.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMethodConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodExists(resourceName, &conf),
					testAccCheckMethodAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "authorization", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "request_models.application/json", "Error"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccMethodImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},

			{
				Config: testAccMethodConfig_update(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodExists(resourceName, &conf),
					testAccCheckMethodAttributesUpdate(&conf),
				),
			},
		},
	})
}

func TestAccAPIGatewayMethod_customAuthorizer(t *testing.T) {
	var conf apigateway.Method
	rInt := sdkacctest.RandInt()
	resourceName := "aws_api_gateway_method.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMethodConfig_customAuthorizer(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodExists(resourceName, &conf),
					testAccCheckMethodAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "authorization", "CUSTOM"),
					resource.TestMatchResourceAttr(resourceName, "authorizer_id", regexp.MustCompile("^[a-z0-9]{6}$")),
					resource.TestCheckResourceAttr(resourceName, "request_models.application/json", "Error"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccMethodImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},

			{
				Config: testAccMethodConfig_update(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodExists(resourceName, &conf),
					testAccCheckMethodAttributesUpdate(&conf),
					resource.TestCheckResourceAttr(resourceName, "authorization", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_id", ""),
				),
			},
		},
	})
}

func TestAccAPIGatewayMethod_cognitoAuthorizer(t *testing.T) {
	var conf apigateway.Method
	rInt := sdkacctest.RandInt()
	resourceName := "aws_api_gateway_method.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMethodConfig_cognitoAuthorizer(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodExists(resourceName, &conf),
					testAccCheckMethodAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "authorization", "COGNITO_USER_POOLS"),
					resource.TestMatchResourceAttr(resourceName, "authorizer_id", regexp.MustCompile("^[a-z0-9]{6}$")),
					resource.TestCheckResourceAttr(resourceName, "request_models.application/json", "Error"),
					resource.TestCheckResourceAttr(resourceName, "authorization_scopes.#", "2"),
				),
			},

			{
				Config: testAccMethodConfig_cognitoAuthorizerUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodExists(resourceName, &conf),
					testAccCheckMethodAttributesUpdate(&conf),
					resource.TestCheckResourceAttr(resourceName, "authorization", "COGNITO_USER_POOLS"),
					resource.TestMatchResourceAttr(resourceName, "authorizer_id", regexp.MustCompile("^[a-z0-9]{6}$")),
					resource.TestCheckResourceAttr(resourceName, "request_models.application/json", "Error"),
					resource.TestCheckResourceAttr(resourceName, "authorization_scopes.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccMethodImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAPIGatewayMethod_customRequestValidator(t *testing.T) {
	var conf apigateway.Method
	rInt := sdkacctest.RandInt()
	resourceName := "aws_api_gateway_method.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMethodConfig_customRequestValidator(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodExists(resourceName, &conf),
					testAccCheckMethodAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "authorization", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "request_models.application/json", "Error"),
					resource.TestMatchResourceAttr(resourceName, "request_validator_id", regexp.MustCompile("^[a-z0-9]{6}$")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccMethodImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},

			{
				Config: testAccMethodConfig_customRequestValidatorUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodExists(resourceName, &conf),
					testAccCheckMethodAttributesUpdate(&conf),
					resource.TestCheckResourceAttr(resourceName, "request_validator_id", ""),
				),
			},
		},
	})
}

func TestAccAPIGatewayMethod_disappears(t *testing.T) {
	var conf apigateway.Method
	rInt := sdkacctest.RandInt()
	resourceName := "aws_api_gateway_method.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMethodConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodExists(resourceName, &conf),
					acctest.CheckResourceDisappears(acctest.Provider, tfapigateway.ResourceMethod(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAPIGatewayMethod_operationName(t *testing.T) {
	var conf apigateway.Method
	rInt := sdkacctest.RandInt()
	resourceName := "aws_api_gateway_method.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMethodConfig_operationName(rInt, "getTest"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "operation_name", "getTest"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccMethodImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				Config: testAccMethodConfig_operationName(rInt, "describeTest"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "operation_name", "describeTest"),
				),
			},
		},
	})
}

func testAccCheckMethodAttributes(conf *apigateway.Method) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *conf.HttpMethod != "GET" {
			return fmt.Errorf("Wrong HttpMethod: %q", *conf.HttpMethod)
		}
		if *conf.AuthorizationType != "NONE" && *conf.AuthorizationType != "CUSTOM" && *conf.AuthorizationType != "COGNITO_USER_POOLS" {
			return fmt.Errorf("Wrong Authorization: %q", *conf.AuthorizationType)
		}

		if val, ok := conf.RequestParameters["method.request.header.Content-Type"]; !ok {
			return fmt.Errorf("missing Content-Type RequestParameters")
		} else {
			if *val {
				return fmt.Errorf("wrong Content-Type RequestParameters value")
			}
		}
		if val, ok := conf.RequestParameters["method.request.querystring.page"]; !ok {
			return fmt.Errorf("missing page RequestParameters")
		} else {
			if !*val {
				return fmt.Errorf("wrong query string RequestParameters value")
			}
		}

		return nil
	}
}

func testAccCheckMethodAttributesUpdate(conf *apigateway.Method) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *conf.HttpMethod != "GET" {
			return fmt.Errorf("Wrong HttpMethod: %q", *conf.HttpMethod)
		}
		if conf.RequestParameters["method.request.header.Content-Type"] != nil {
			return fmt.Errorf("Content-Type RequestParameters shouldn't exist")
		}
		if val, ok := conf.RequestParameters["method.request.querystring.page"]; !ok {
			return fmt.Errorf("missing updated page RequestParameters")
		} else {
			if *val {
				return fmt.Errorf("wrong query string RequestParameters updated value")
			}
		}

		return nil
	}
}

func testAccCheckMethodExists(n string, res *apigateway.Method) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No API Gateway Method ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayConn

		req := &apigateway.GetMethodInput{
			HttpMethod: aws.String("GET"),
			ResourceId: aws.String(s.RootModule().Resources["aws_api_gateway_resource.test"].Primary.ID),
			RestApiId:  aws.String(s.RootModule().Resources["aws_api_gateway_rest_api.test"].Primary.ID),
		}
		describe, err := conn.GetMethod(req)
		if err != nil {
			return err
		}

		*res = *describe

		return nil
	}
}

func testAccCheckMethodDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_api_gateway_method" {
			continue
		}

		req := &apigateway.GetMethodInput{
			HttpMethod: aws.String("GET"),
			ResourceId: aws.String(s.RootModule().Resources["aws_api_gateway_resource.test"].Primary.ID),
			RestApiId:  aws.String(s.RootModule().Resources["aws_api_gateway_rest_api.test"].Primary.ID),
		}
		_, err := conn.GetMethod(req)

		if err == nil {
			return fmt.Errorf("API Gateway Method still exists")
		}

		aws2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if aws2err.Code() != "NotFoundException" {
			return err
		}

		return nil
	}

	return nil
}

func testAccMethodImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["rest_api_id"], rs.Primary.Attributes["resource_id"], rs.Primary.Attributes["http_method"]), nil
	}
}

func testAccMethodConfig_customAuthorizer(rInt int) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = "tf-acc-test-custom-auth-%d"
}

resource "aws_iam_role" "invocation_role" {
  name = "tf_acc_api_gateway_auth_invocation_role-%d"
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

resource "aws_iam_role_policy" "invocation_policy" {
  name = "tf-acc-api-gateway-%d"
  role = aws_iam_role.invocation_role.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "lambda:InvokeFunction",
      "Effect": "Allow",
      "Resource": "${aws_lambda_function.authorizer.arn}"
    }
  ]
}
EOF
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "tf_acc_iam_for_lambda_api_gateway_authorizer-%d"

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

resource "aws_lambda_function" "authorizer" {
  filename         = "test-fixtures/lambdatest.zip"
  source_code_hash = filebase64sha256("test-fixtures/lambdatest.zip")
  function_name    = "tf_acc_api_gateway_authorizer_%d"
  role             = aws_iam_role.iam_for_lambda.arn
  handler          = "exports.example"
  runtime          = "nodejs12.x"
}

resource "aws_api_gateway_authorizer" "test" {
  name                   = "tf-acc-test-authorizer"
  rest_api_id            = aws_api_gateway_rest_api.test.id
  authorizer_uri         = aws_lambda_function.authorizer.invoke_arn
  authorizer_credentials = aws_iam_role.invocation_role.arn
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
  authorization = "CUSTOM"
  authorizer_id = aws_api_gateway_authorizer.test.id

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.header.Content-Type" = false
    "method.request.querystring.page"    = true
  }
}
`, rInt, rInt, rInt, rInt, rInt)
}

func testAccMethodConfig_cognitoAuthorizer(rInt int) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = "tf-acc-test-cognito-auth-%d"
}

resource "aws_iam_role" "invocation_role" {
  name = "tf_acc_api_gateway_auth_invocation_role-%d"
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

resource "aws_iam_role" "iam_for_lambda" {
  name = "tf_acc_iam_for_lambda_api_gateway_authorizer-%d"

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

resource "aws_cognito_user_pool" "pool" {
  name = "tf-acc-test-cognito-pool-%d"
}

resource "aws_api_gateway_authorizer" "test" {
  name            = "tf-acc-test-cognito-authorizer"
  rest_api_id     = aws_api_gateway_rest_api.test.id
  identity_source = "method.request.header.Authorization"
  provider_arns   = [aws_cognito_user_pool.pool.arn]
  type            = "COGNITO_USER_POOLS"
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  parent_id   = aws_api_gateway_rest_api.test.root_resource_id
  path_part   = "test"
}

resource "aws_api_gateway_method" "test" {
  rest_api_id          = aws_api_gateway_rest_api.test.id
  resource_id          = aws_api_gateway_resource.test.id
  http_method          = "GET"
  authorization        = "COGNITO_USER_POOLS"
  authorizer_id        = aws_api_gateway_authorizer.test.id
  authorization_scopes = ["test.read", "test.write"]

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.header.Content-Type" = false
    "method.request.querystring.page"    = true
  }
}
`, rInt, rInt, rInt, rInt)
}

func testAccMethodConfig_cognitoAuthorizerUpdate(rInt int) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = "tf-acc-test-cognito-auth-%d"
}

resource "aws_iam_role" "invocation_role" {
  name = "tf_acc_api_gateway_auth_invocation_role-%d"
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

resource "aws_iam_role" "iam_for_lambda" {
  name = "tf_acc_iam_for_lambda_api_gateway_authorizer-%d"

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

resource "aws_cognito_user_pool" "pool" {
  name = "tf-acc-test-cognito-pool-%d"
}

resource "aws_api_gateway_authorizer" "test" {
  name            = "tf-acc-test-cognito-authorizer"
  rest_api_id     = aws_api_gateway_rest_api.test.id
  identity_source = "method.request.header.Authorization"
  provider_arns   = [aws_cognito_user_pool.pool.arn]
  type            = "COGNITO_USER_POOLS"
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  parent_id   = aws_api_gateway_rest_api.test.root_resource_id
  path_part   = "test"
}

resource "aws_api_gateway_method" "test" {
  rest_api_id          = aws_api_gateway_rest_api.test.id
  resource_id          = aws_api_gateway_resource.test.id
  http_method          = "GET"
  authorization        = "COGNITO_USER_POOLS"
  authorizer_id        = aws_api_gateway_authorizer.test.id
  authorization_scopes = ["test.read", "test.write", "test.delete"]

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.querystring.page" = false
  }
}
`, rInt, rInt, rInt, rInt)
}

func testAccMethodConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = "tf-acc-test-apig-method-%d"
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

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.header.Content-Type" = false
    "method.request.querystring.page"    = true
  }
}
`, rInt)
}

func testAccMethodConfig_update(rInt int) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = "tf-acc-test-apig-method-%d"
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

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.querystring.page" = false
  }
}
`, rInt)
}

func testAccMethodConfig_customRequestValidator(rInt int) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = "tf-acc-test-apig-method-custom-req-validator-%d"
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  parent_id   = aws_api_gateway_rest_api.test.root_resource_id
  path_part   = "test"
}

resource "aws_api_gateway_request_validator" "validator" {
  rest_api_id                 = aws_api_gateway_rest_api.test.id
  name                        = "paramsValidator"
  validate_request_parameters = true
}

resource "aws_api_gateway_method" "test" {
  rest_api_id   = aws_api_gateway_rest_api.test.id
  resource_id   = aws_api_gateway_resource.test.id
  http_method   = "GET"
  authorization = "NONE"

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.header.Content-Type" = false
    "method.request.querystring.page"    = true
  }

  request_validator_id = aws_api_gateway_request_validator.validator.id
}
`, rInt)
}

func testAccMethodConfig_customRequestValidatorUpdate(rInt int) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = "tf-acc-test-apig-method-custom-req-validator-%d"
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  parent_id   = aws_api_gateway_rest_api.test.root_resource_id
  path_part   = "test"
}

resource "aws_api_gateway_request_validator" "validator" {
  rest_api_id                 = aws_api_gateway_rest_api.test.id
  name                        = "paramsValidator"
  validate_request_parameters = true
}

resource "aws_api_gateway_method" "test" {
  rest_api_id   = aws_api_gateway_rest_api.test.id
  resource_id   = aws_api_gateway_resource.test.id
  http_method   = "GET"
  authorization = "NONE"

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.querystring.page" = false
  }
}
`, rInt)
}

func testAccMethodConfig_operationName(rInt int, operationName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = "tf-acc-test-apig-method-custom-op-name-%[1]d"
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  parent_id   = aws_api_gateway_rest_api.test.root_resource_id
  path_part   = "test"
}

resource "aws_api_gateway_method" "test" {
  authorization  = "NONE"
  http_method    = "GET"
  operation_name = %[2]q
  resource_id    = aws_api_gateway_resource.test.id
  rest_api_id    = aws_api_gateway_rest_api.test.id

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.header.Content-Type" = false
    "method.request.querystring.page"    = true
  }
}
`, rInt, operationName)
}
