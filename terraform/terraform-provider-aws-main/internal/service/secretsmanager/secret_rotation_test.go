package secretsmanager_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccSecretsManagerSecretRotation_basic(t *testing.T) {
	var secret secretsmanager.DescribeSecretOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_secretsmanager_secret_rotation.test"
	lambdaFunctionResourceName := "aws_lambda_function.test1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, secretsmanager.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecretRotationDestroy,
		Steps: []resource.TestStep{
			// Test creating secret rotation resource
			{
				Config: testAccSecretRotationConfig_basic(rName, 7),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretRotationExists(resourceName, &secret),
					resource.TestCheckResourceAttr(resourceName, "rotation_enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "rotation_lambda_arn", lambdaFunctionResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rules.0.automatically_after_days", "7"),
				),
			},
			// Test updating rotation
			// We need a valid rotation function for this testing
			// InvalidRequestException: A previous rotation isn’t complete. That rotation will be reattempted.
			/*
				{
					Config: testAccSecretRotationConfig_managerUpdated(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSecretRotationExists(resourceName, &secret),
						resource.TestCheckResourceAttr(resourceName, "rotation_enabled", "true"),
						resource.TestMatchResourceAttr(resourceName, "rotation_lambda_arn", regexp.MustCompile(fmt.Sprintf("^arn:[^:]+:lambda:[^:]+:[^:]+:function:%s-2$", rName))),
					),
				},
			*/
			// Test importing secret rotation
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSecretRotationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SecretsManagerConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_secretsmanager_secret_rotation" {
			continue
		}

		input := &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(rs.Primary.ID),
		}

		output, err := conn.DescribeSecret(input)

		if err != nil {
			if tfawserr.ErrCodeEquals(err, secretsmanager.ErrCodeResourceNotFoundException) {
				continue
			}
			return err
		}

		if output != nil && aws.BoolValue(output.RotationEnabled) {
			return fmt.Errorf("Secret rotation for %q still enabled", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckSecretRotationExists(resourceName string, secret *secretsmanager.DescribeSecretOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SecretsManagerConn
		input := &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(rs.Primary.ID),
		}

		output, err := conn.DescribeSecret(input)

		if err != nil {
			return err
		}

		if output == nil {
			return fmt.Errorf("Secret %q does not exist", rs.Primary.ID)
		}

		if !aws.BoolValue(output.RotationEnabled) {
			return fmt.Errorf("Secret rotation %q is not enabled", rs.Primary.ID)
		}

		*secret = *output

		return nil
	}
}

func testAccSecretRotationConfig_basic(rName string, automaticallyAfterDays int) string {
	return acctest.ConfigLambdaBase(rName, rName, rName) + fmt.Sprintf(`
# Not a real rotation function
resource "aws_lambda_function" "test1" {
  filename      = "test-fixtures/lambdatest.zip"
  function_name = "%[1]s-1"
  handler       = "exports.example"
  role          = aws_iam_role.iam_for_lambda.arn
  runtime       = "nodejs12.x"
}

resource "aws_lambda_permission" "test1" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.test1.function_name
  principal     = "secretsmanager.amazonaws.com"
  statement_id  = "AllowExecutionFromSecretsManager1"
}

# Not a real rotation function
resource "aws_lambda_function" "test2" {
  filename      = "test-fixtures/lambdatest.zip"
  function_name = "%[1]s-2"
  handler       = "exports.example"
  role          = aws_iam_role.iam_for_lambda.arn
  runtime       = "nodejs12.x"
}

resource "aws_lambda_permission" "test2" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.test2.function_name
  principal     = "secretsmanager.amazonaws.com"
  statement_id  = "AllowExecutionFromSecretsManager2"
}

resource "aws_secretsmanager_secret" "test" {
  name = "%[1]s"
}

resource "aws_secretsmanager_secret_rotation" "test" {
  secret_id           = aws_secretsmanager_secret.test.id
  rotation_lambda_arn = aws_lambda_function.test1.arn

  rotation_rules {
    automatically_after_days = %[2]d
  }

  depends_on = [aws_lambda_permission.test1]
}
`, rName, automaticallyAfterDays)
}
