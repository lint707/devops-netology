package configservice_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/configservice"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfconfigservice "github.com/hashicorp/terraform-provider-aws/internal/service/configservice"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func testAccRemediationConfiguration_basic(t *testing.T) {
	var rc configservice.RemediationConfiguration
	resourceName := "aws_config_remediation_configuration.test"
	rInt := sdkacctest.RandInt()
	automatic := "false"
	rAttempts := sdkacctest.RandIntRange(1, 25)
	rSeconds := sdkacctest.RandIntRange(1, 2678000)
	rExecPct := sdkacctest.RandIntRange(1, 100)
	rErrorPct := sdkacctest.RandIntRange(1, 100)
	prefix := "Original"
	sseAlgorithm := "AES256"
	expectedName := fmt.Sprintf("%s-tf-acc-test-%d", prefix, rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, configservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRemediationConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRemediationConfigurationConfig_basic(prefix, sseAlgorithm, rInt, rAttempts, rSeconds, rExecPct, rErrorPct, automatic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemediationConfigurationExists(resourceName, &rc),
					resource.TestCheckResourceAttr(resourceName, "config_rule_name", expectedName),
					resource.TestCheckResourceAttr(resourceName, "target_id", "AWS-EnableS3BucketEncryption"),
					resource.TestCheckResourceAttr(resourceName, "target_type", "SSM_DOCUMENT"),
					resource.TestCheckResourceAttr(resourceName, "parameter.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "automatic", automatic),
					resource.TestCheckResourceAttr(resourceName, "maximum_automatic_attempts", strconv.Itoa(rAttempts)),
					resource.TestCheckResourceAttr(resourceName, "retry_attempt_seconds", strconv.Itoa(rSeconds)),
					resource.TestCheckResourceAttr(resourceName, "execution_controls.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "execution_controls.*.ssm_controls.*", map[string]string{
						"concurrent_execution_rate_percentage": strconv.Itoa(rExecPct),
						"error_percentage":                     strconv.Itoa(rErrorPct),
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

func testAccRemediationConfiguration_basicBackwardCompatible(t *testing.T) {
	var rc configservice.RemediationConfiguration
	resourceName := "aws_config_remediation_configuration.test"
	rInt := sdkacctest.RandInt()
	prefix := "Original"
	sseAlgorithm := "AES256"
	expectedName := fmt.Sprintf("%s-tf-acc-test-%d", prefix, rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, configservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRemediationConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRemediationConfigurationConfig_olderSchema(prefix, sseAlgorithm, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemediationConfigurationExists(resourceName, &rc),
					resource.TestCheckResourceAttr(resourceName, "config_rule_name", expectedName),
					resource.TestCheckResourceAttr(resourceName, "target_id", "AWS-EnableS3BucketEncryption"),
					resource.TestCheckResourceAttr(resourceName, "target_type", "SSM_DOCUMENT"),
					resource.TestCheckResourceAttr(resourceName, "parameter.#", "3"),
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

func testAccRemediationConfiguration_disappears(t *testing.T) {
	var rc configservice.RemediationConfiguration
	resourceName := "aws_config_remediation_configuration.test"
	rInt := sdkacctest.RandInt()
	automatic := "false"
	rAttempts := sdkacctest.RandIntRange(1, 25)
	rSeconds := sdkacctest.RandIntRange(1, 2678000)
	rExecPct := sdkacctest.RandIntRange(1, 100)
	rErrorPct := sdkacctest.RandIntRange(1, 100)
	prefix := "original"
	sseAlgorithm := "AES256"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, configservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRemediationConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRemediationConfigurationConfig_basic(prefix, sseAlgorithm, rInt, rAttempts, rSeconds, rExecPct, rErrorPct, automatic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemediationConfigurationExists(resourceName, &rc),
					acctest.CheckResourceDisappears(acctest.Provider, tfconfigservice.ResourceRemediationConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRemediationConfiguration_recreates(t *testing.T) {
	var original configservice.RemediationConfiguration
	var updated configservice.RemediationConfiguration
	resourceName := "aws_config_remediation_configuration.test"
	rInt := sdkacctest.RandInt()
	automatic := "false"
	rAttempts := sdkacctest.RandIntRange(1, 25)
	rSeconds := sdkacctest.RandIntRange(1, 2678000)
	rExecPct := sdkacctest.RandIntRange(1, 100)
	rErrorPct := sdkacctest.RandIntRange(1, 100)

	originalName := "Original"
	updatedName := "Updated"
	sseAlgorithm := "AES256"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, configservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRemediationConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRemediationConfigurationConfig_basic(originalName, sseAlgorithm, rInt, rAttempts, rSeconds, rExecPct, rErrorPct, automatic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemediationConfigurationExists(resourceName, &original),
					resource.TestCheckResourceAttr(resourceName, "config_rule_name", fmt.Sprintf("%s-tf-acc-test-%d", originalName, rInt)),
				),
			},
			{
				Config: testAccRemediationConfigurationConfig_basic(updatedName, sseAlgorithm, rInt, rAttempts, rSeconds, rExecPct, rErrorPct, automatic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemediationConfigurationExists(resourceName, &updated),
					testAccCheckRemediationConfigurationRecreated(&original, &updated),
					resource.TestCheckResourceAttr(resourceName, "config_rule_name", fmt.Sprintf("%s-tf-acc-test-%d", updatedName, rInt)),
				),
			},
		},
	})
}

func testAccRemediationConfiguration_updates(t *testing.T) {
	var original configservice.RemediationConfiguration
	var updated configservice.RemediationConfiguration
	resourceName := "aws_config_remediation_configuration.test"
	rInt := sdkacctest.RandInt()
	automatic := "false"
	rAttempts := sdkacctest.RandIntRange(1, 25)
	rSeconds := sdkacctest.RandIntRange(1, 2678000)
	rExecPct := sdkacctest.RandIntRange(1, 100)
	rErrorPct := sdkacctest.RandIntRange(1, 100)
	uAutomatic := "true"
	uAttempts := sdkacctest.RandIntRange(1, 25)
	uSeconds := sdkacctest.RandIntRange(1, 2678000)
	uExecPct := sdkacctest.RandIntRange(1, 100)
	uErrorPct := sdkacctest.RandIntRange(1, 100)

	name := "Original"
	originalSseAlgorithm := "AES256"
	updatedSseAlgorithm := "aws:kms"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, configservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRemediationConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRemediationConfigurationConfig_basic(name, originalSseAlgorithm, rInt, rAttempts, rSeconds, rExecPct, rErrorPct, automatic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemediationConfigurationExists(resourceName, &original),
					resource.TestCheckResourceAttr(resourceName, "parameter.2.static_value", originalSseAlgorithm),
					resource.TestCheckResourceAttr(resourceName, "automatic", automatic),
					resource.TestCheckResourceAttr(resourceName, "maximum_automatic_attempts", strconv.Itoa(rAttempts)),
					resource.TestCheckResourceAttr(resourceName, "retry_attempt_seconds", strconv.Itoa(rSeconds)),
					resource.TestCheckResourceAttr(resourceName, "execution_controls.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "execution_controls.*.ssm_controls.*", map[string]string{
						"concurrent_execution_rate_percentage": strconv.Itoa(rExecPct),
						"error_percentage":                     strconv.Itoa(rErrorPct),
					}),
				),
			},
			{
				Config: testAccRemediationConfigurationConfig_basic(name, updatedSseAlgorithm, rInt, uAttempts, uSeconds, uExecPct, uErrorPct, uAutomatic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemediationConfigurationExists(resourceName, &updated),
					testAccCheckRemediationConfigurationNotRecreated(&original, &updated),
					resource.TestCheckResourceAttr(resourceName, "parameter.2.static_value", updatedSseAlgorithm),
					resource.TestCheckResourceAttr(resourceName, "automatic", uAutomatic),
					resource.TestCheckResourceAttr(resourceName, "maximum_automatic_attempts", strconv.Itoa(uAttempts)),
					resource.TestCheckResourceAttr(resourceName, "retry_attempt_seconds", strconv.Itoa(uSeconds)),
					resource.TestCheckResourceAttr(resourceName, "execution_controls.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "execution_controls.*.ssm_controls.*", map[string]string{
						"concurrent_execution_rate_percentage": strconv.Itoa(uExecPct),
						"error_percentage":                     strconv.Itoa(uErrorPct),
					}),
				),
			},
		},
	})
}

func testAccRemediationConfiguration_values(t *testing.T) {
	var rc configservice.RemediationConfiguration
	resourceName := "aws_config_remediation_configuration.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	automatic := "false"
	rAttempts := sdkacctest.RandIntRange(1, 25)
	rSeconds := sdkacctest.RandIntRange(1, 2678000)
	rExecPct := sdkacctest.RandIntRange(1, 100)
	rErrorPct := sdkacctest.RandIntRange(1, 100)
	sseAlgorithm := "AES256"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, configservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRemediationConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRemediationConfigurationConfig_values(rName, sseAlgorithm, rAttempts, rSeconds, rExecPct, rErrorPct, automatic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemediationConfigurationExists(resourceName, &rc),
					resource.TestCheckResourceAttr(resourceName, "target_id", "AWS-EnableS3BucketEncryption"),
					resource.TestCheckResourceAttr(resourceName, "target_type", "SSM_DOCUMENT"),
					resource.TestCheckResourceAttr(resourceName, "parameter.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "automatic", automatic),
					resource.TestCheckResourceAttr(resourceName, "maximum_automatic_attempts", strconv.Itoa(rAttempts)),
					resource.TestCheckResourceAttr(resourceName, "retry_attempt_seconds", strconv.Itoa(rSeconds)),
					resource.TestCheckResourceAttr(resourceName, "execution_controls.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "execution_controls.*.ssm_controls.*", map[string]string{
						"concurrent_execution_rate_percentage": strconv.Itoa(rExecPct),
						"error_percentage":                     strconv.Itoa(rErrorPct),
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

func testAccCheckRemediationConfigurationExists(n string, obj *configservice.RemediationConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return create.Error(names.ConfigService, create.ErrActionCheckingExistence, tfconfigservice.ResNameRemediationConfiguration, n, errors.New("not found in state"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.ConfigService, create.ErrActionCheckingExistence, tfconfigservice.ResNameRemediationConfiguration, n, errors.New("ID not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ConfigServiceConn
		out, err := conn.DescribeRemediationConfigurations(&configservice.DescribeRemediationConfigurationsInput{
			ConfigRuleNames: []*string{aws.String(rs.Primary.Attributes["config_rule_name"])},
		})
		if err != nil {
			return create.Error(names.ConfigService, create.ErrActionCheckingExistence, tfconfigservice.ResNameRemediationConfiguration, n, err)
		}
		if len(out.RemediationConfigurations) < 1 {
			return create.Error(names.ConfigService, create.ErrActionCheckingExistence, tfconfigservice.ResNameRemediationConfiguration, n, errors.New("not found"))
		}

		rc := out.RemediationConfigurations[0]
		*obj = *rc

		return nil
	}
}

func testAccCheckRemediationConfigurationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ConfigServiceConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_config_remediation_configuration" {
			continue
		}

		resp, err := conn.DescribeRemediationConfigurations(&configservice.DescribeRemediationConfigurationsInput{
			ConfigRuleNames: []*string{aws.String(rs.Primary.Attributes["config_rule_name"])},
		})

		if err == nil {
			if len(resp.RemediationConfigurations) != 0 &&
				aws.StringValue(resp.RemediationConfigurations[0].ConfigRuleName) == rs.Primary.Attributes["name"] {
				return create.Error(names.ConfigService, create.ErrActionCheckingDestroyed, tfconfigservice.ResNameRemediationConfiguration, rs.Primary.Attributes["name"], errors.New("still exists"))
			}
		}
	}

	return nil
}

func testAccCheckRemediationConfigurationNotRecreated(before, after *configservice.RemediationConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(before.Arn) != aws.StringValue(after.Arn) {
			return create.Error(names.ConfigService, create.ErrActionCheckingNotRecreated, tfconfigservice.ResNameRemediationConfiguration, aws.StringValue(before.Arn), fmt.Errorf("ARNs changed, new: %s", aws.StringValue(after.Arn)))
		}
		return nil
	}
}

func testAccCheckRemediationConfigurationRecreated(before, after *configservice.RemediationConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(before.Arn) == aws.StringValue(after.Arn) {
			return create.Error(names.ConfigService, create.ErrActionCheckingRecreated, tfconfigservice.ResNameRemediationConfiguration, aws.StringValue(before.Arn), fmt.Errorf("wasn't recreated, new: %s", aws.StringValue(after.Arn)))
		}
		return nil
	}
}

func testAccRemediationConfigurationConfig_olderSchema(namePrefix, sseAlgorithm string, randInt int) string {
	return fmt.Sprintf(`
resource "aws_config_remediation_configuration" "test" {
  config_rule_name = aws_config_config_rule.test.name

  resource_type  = "AWS::S3::Bucket"
  target_id      = "AWS-EnableS3BucketEncryption"
  target_type    = "SSM_DOCUMENT"
  target_version = "1"

  parameter {
    name         = "AutomationAssumeRole"
    static_value = aws_iam_role.test.arn
  }
  parameter {
    name           = "BucketName"
    resource_value = "RESOURCE_ID"
  }
  parameter {
    name         = "SSEAlgorithm"
    static_value = "%[2]s"
  }
}

resource "aws_sns_topic" "test" {
  name = "sns_topic_name"
}

resource "aws_config_config_rule" "test" {
  name = "%[1]s-tf-acc-test-%[3]d"

  source {
    owner             = "AWS"
    source_identifier = "S3_BUCKET_VERSIONING_ENABLED"
  }

  depends_on = [aws_config_configuration_recorder.test]
}

resource "aws_config_configuration_recorder" "test" {
  name     = "%[1]s-tf-acc-test-%[3]d"
  role_arn = aws_iam_role.test.arn
}

resource "aws_iam_role" "test" {
  name = "%[1]s-tf-acc-test-awsconfig-%[3]d"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "config.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "test" {
  name = "%[1]s-tf-acc-test-awsconfig-%[3]d"
  role = aws_iam_role.test.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
        "Action": "config:Put*",
        "Effect": "Allow",
        "Resource": "*"

    }
  ]
}
EOF
}
`, namePrefix, sseAlgorithm, randInt)
}

func testAccRemediationConfigurationConfig_basic(namePrefix, sseAlgorithm string, randInt int, randAttempts int, randSeconds int, randExecPct int, randErrorPct int, automatic string) string {
	return fmt.Sprintf(`
resource "aws_config_remediation_configuration" "test" {
  config_rule_name = aws_config_config_rule.test.name

  resource_type  = "AWS::S3::Bucket"
  target_id      = "AWS-EnableS3BucketEncryption"
  target_type    = "SSM_DOCUMENT"
  target_version = "1"

  parameter {
    name         = "AutomationAssumeRole"
    static_value = aws_iam_role.test.arn
  }
  parameter {
    name           = "BucketName"
    resource_value = "RESOURCE_ID"
  }
  parameter {
    name         = "SSEAlgorithm"
    static_value = "%[2]s"
  }
  automatic                  = %[8]s
  maximum_automatic_attempts = %[4]d
  retry_attempt_seconds      = %[5]d
  execution_controls {
    ssm_controls {
      concurrent_execution_rate_percentage = %[6]d
      error_percentage                     = %[7]d
    }
  }
}

resource "aws_sns_topic" "test" {
  name = "sns_topic_name"
}

resource "aws_config_config_rule" "test" {
  name = "%[1]s-tf-acc-test-%[3]d"

  source {
    owner             = "AWS"
    source_identifier = "S3_BUCKET_VERSIONING_ENABLED"
  }

  depends_on = [aws_config_configuration_recorder.test]
}

resource "aws_config_configuration_recorder" "test" {
  name     = "%[1]s-tf-acc-test-%[3]d"
  role_arn = aws_iam_role.test.arn
}

resource "aws_iam_role" "test" {
  name = "%[1]s-tf-acc-test-awsconfig-%[3]d"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "config.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "test" {
  name = "%[1]s-tf-acc-test-awsconfig-%[3]d"
  role = aws_iam_role.test.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
        "Action": "config:Put*",
        "Effect": "Allow",
        "Resource": "*"

    }
  ]
}
EOF
}
`, namePrefix, sseAlgorithm, randInt, randAttempts, randSeconds, randExecPct, randErrorPct, automatic)
}

func testAccRemediationConfigurationConfig_values(rName, sseAlgorithm string, randAttempts int, randSeconds int, randExecPct int, randErrorPct int, automatic string) string {
	return fmt.Sprintf(`
resource "aws_config_remediation_configuration" "test" {
  config_rule_name = aws_config_config_rule.test.name

  resource_type  = "AWS::S3::Bucket"
  target_id      = "AWS-EnableS3BucketEncryption"
  target_type    = "SSM_DOCUMENT"
  target_version = "1"

  parameter {
    name          = "AutomationAssumeRole"
    static_values = [aws_iam_role.test.arn, aws_iam_role.test2.arn]
  }

  parameter {
    name           = "BucketName"
    resource_value = "RESOURCE_ID"
  }

  parameter {
    name         = "SSEAlgorithm"
    static_value = "%[2]s"
  }

  automatic                  = %[7]s
  maximum_automatic_attempts = %[3]d
  retry_attempt_seconds      = %[4]d

  execution_controls {
    ssm_controls {
      concurrent_execution_rate_percentage = %[5]d
      error_percentage                     = %[6]d
    }
  }
}

resource "aws_sns_topic" "test" {
  name = "sns_topic_name"
}

resource "aws_config_config_rule" "test" {
  name = %[1]q

  source {
    owner             = "AWS"
    source_identifier = "S3_BUCKET_VERSIONING_ENABLED"
  }

  depends_on = [aws_config_configuration_recorder.test]
}

resource "aws_config_configuration_recorder" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn
}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "config.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role" "test2" {
  name = "%[1]s-2"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "config.amazonaws.com"
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
        "Action": "config:Put*",
        "Effect": "Allow",
        "Resource": "*"

    }
  ]
}
EOF
}
`, rName, sseAlgorithm, randAttempts, randSeconds, randExecPct, randErrorPct, automatic)
}
