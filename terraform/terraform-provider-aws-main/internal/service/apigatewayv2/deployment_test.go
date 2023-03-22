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
)

func TestAccAPIGatewayV2Deployment_basic(t *testing.T) {
	var apiId string
	var v apigatewayv2.GetDeploymentOutput
	resourceName := "aws_apigatewayv2_deployment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigatewayv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentConfig_basic(rName, "Test description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_deployed", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test description"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccDeploymentImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDeploymentConfig_basic(rName, "Test description updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentExists(resourceName, &apiId, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_deployed", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test description updated"),
				),
			},
		},
	})
}

func TestAccAPIGatewayV2Deployment_disappears(t *testing.T) {
	var apiId string
	var v apigatewayv2.GetDeploymentOutput
	resourceName := "aws_apigatewayv2_deployment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigatewayv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentConfig_basic(rName, "Test description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentExists(resourceName, &apiId, &v),
					testAccCheckDeploymentDisappears(&apiId, &v),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAPIGatewayV2Deployment_triggers(t *testing.T) {
	var apiId string
	var deployment1, deployment2, deployment3, deployment4 apigatewayv2.GetDeploymentOutput
	resourceName := "aws_apigatewayv2_deployment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigatewayv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentConfig_triggers(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentExists(resourceName, &apiId, &deployment1),
				),
				// Due to how the Terraform state is handled for resources during creation,
				// any SHA1 of whole resources will change after first apply, then stabilize.
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccDeploymentConfig_triggers(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentExists(resourceName, &apiId, &deployment2),
					testAccCheckDeploymentRecreated(&deployment1, &deployment2),
				),
			},
			{
				Config: testAccDeploymentConfig_triggers(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentExists(resourceName, &apiId, &deployment3),
					testAccCheckDeploymentNotRecreated(&deployment2, &deployment3),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccDeploymentImportStateIdFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"triggers"},
			},
			{
				Config: testAccDeploymentConfig_triggers(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentExists(resourceName, &apiId, &deployment4),
					testAccCheckDeploymentRecreated(&deployment3, &deployment4),
				),
			},
		},
	})
}

func testAccCheckDeploymentDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayV2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_apigatewayv2_deployment" {
			continue
		}

		_, err := conn.GetDeployment(&apigatewayv2.GetDeploymentInput{
			ApiId:        aws.String(rs.Primary.Attributes["api_id"]),
			DeploymentId: aws.String(rs.Primary.ID),
		})
		if tfawserr.ErrCodeEquals(err, apigatewayv2.ErrCodeNotFoundException) {
			continue
		}
		if err != nil {
			return err
		}

		return fmt.Errorf("API Gateway v2 deployment %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckDeploymentDisappears(apiId *string, v *apigatewayv2.GetDeploymentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayV2Conn

		_, err := conn.DeleteDeployment(&apigatewayv2.DeleteDeploymentInput{
			ApiId:        apiId,
			DeploymentId: v.DeploymentId,
		})

		return err
	}
}

func testAccCheckDeploymentExists(n string, vApiId *string, v *apigatewayv2.GetDeploymentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No API Gateway v2 deployment ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayV2Conn

		apiId := aws.String(rs.Primary.Attributes["api_id"])
		resp, err := conn.GetDeployment(&apigatewayv2.GetDeploymentInput{
			ApiId:        apiId,
			DeploymentId: aws.String(rs.Primary.ID),
		})
		if err != nil {
			return err
		}

		*vApiId = *apiId
		*v = *resp

		return nil
	}
}

func testAccCheckDeploymentNotRecreated(i, j *apigatewayv2.GetDeploymentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !aws.TimeValue(i.CreatedDate).Equal(aws.TimeValue(j.CreatedDate)) {
			return fmt.Errorf("API Gateway V2 Deployment recreated")
		}

		return nil
	}
}

func testAccCheckDeploymentRecreated(i, j *apigatewayv2.GetDeploymentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.TimeValue(i.CreatedDate).Equal(aws.TimeValue(j.CreatedDate)) {
			return fmt.Errorf("API Gateway V2 Deployment not recreated")
		}

		return nil
	}
}

func testAccDeploymentImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["api_id"], rs.Primary.ID), nil
	}
}

func testAccDeploymentConfig_basic(rName, description string) string {
	return testAccRouteConfig_target(rName) + fmt.Sprintf(`
resource "aws_apigatewayv2_deployment" "test" {
  api_id      = aws_apigatewayv2_api.test.id
  description = %[1]q

  depends_on = [aws_apigatewayv2_route.test]
}
`, description)
}

func testAccDeploymentConfig_triggers(rName string, apiKeyRequired bool) string {
	return fmt.Sprintf(`
resource "aws_apigatewayv2_api" "test" {
  name                       = %[1]q
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_integration" "test" {
  api_id           = aws_apigatewayv2_api.test.id
  integration_type = "MOCK"
}

resource "aws_apigatewayv2_route" "test" {
  api_id           = aws_apigatewayv2_api.test.id
  api_key_required = %[2]t
  route_key        = "$default"
  target           = "integrations/${aws_apigatewayv2_integration.test.id}"
}

resource "aws_apigatewayv2_deployment" "test" {
  api_id = aws_apigatewayv2_api.test.id

  triggers = {
    redeployment = sha1(jsonencode([
      aws_apigatewayv2_integration.test,
      aws_apigatewayv2_route.test,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}
`, rName, apiKeyRequired)
}
