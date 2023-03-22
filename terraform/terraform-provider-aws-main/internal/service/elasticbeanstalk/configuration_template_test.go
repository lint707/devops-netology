package elasticbeanstalk_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccElasticBeanstalkConfigurationTemplate_Beanstalk_basic(t *testing.T) {
	var config elasticbeanstalk.ConfigurationSettingsDescription

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConfigurationTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigurationTemplateConfig_basic(sdkacctest.RandString(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationTemplateExists("aws_elastic_beanstalk_configuration_template.tf_template", &config),
				),
			},
		},
	})
}

func TestAccElasticBeanstalkConfigurationTemplate_Beanstalk_vpc(t *testing.T) {
	var config elasticbeanstalk.ConfigurationSettingsDescription

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConfigurationTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigurationTemplateConfig_vpc(sdkacctest.RandString(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationTemplateExists("aws_elastic_beanstalk_configuration_template.tf_template", &config),
				),
			},
		},
	})
}

func TestAccElasticBeanstalkConfigurationTemplate_Beanstalk_setting(t *testing.T) {
	var config elasticbeanstalk.ConfigurationSettingsDescription

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConfigurationTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigurationTemplateConfig_setting(sdkacctest.RandString(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationTemplateExists("aws_elastic_beanstalk_configuration_template.tf_template", &config),
					resource.TestCheckResourceAttr(
						"aws_elastic_beanstalk_configuration_template.tf_template", "setting.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("aws_elastic_beanstalk_configuration_template.tf_template", "setting.*", map[string]string{
						"value": "m1.small",
					}),
				),
			},
		},
	})
}

func testAccCheckConfigurationTemplateDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ElasticBeanstalkConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_elastic_beanstalk_configuration_template" {
			continue
		}

		// Try to find the Configuration Template
		opts := elasticbeanstalk.DescribeConfigurationSettingsInput{
			TemplateName:    aws.String(rs.Primary.ID),
			ApplicationName: aws.String(rs.Primary.Attributes["application"]),
		}
		resp, err := conn.DescribeConfigurationSettings(&opts)
		if err == nil {
			if len(resp.ConfigurationSettings) > 0 {
				return fmt.Errorf("Elastic Beanstalk Application still exists.")
			}

			return nil
		}

		// Verify the error is what we want
		ec2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}

		switch {
		case ec2err.Code() == "InvalidBeanstalkConfigurationTemplateID.NotFound":
			return nil
		// This error can be returned when the beanstalk application no longer exists.
		case ec2err.Code() == "InvalidParameterValue":
			return nil
		default:
			return err
		}
	}

	return nil
}

func testAccCheckConfigurationTemplateExists(n string, config *elasticbeanstalk.ConfigurationSettingsDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ElasticBeanstalkConn
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Elastic Beanstalk config ID is not set")
		}

		opts := elasticbeanstalk.DescribeConfigurationSettingsInput{
			TemplateName:    aws.String(rs.Primary.ID),
			ApplicationName: aws.String(rs.Primary.Attributes["application"]),
		}
		resp, err := conn.DescribeConfigurationSettings(&opts)
		if err != nil {
			return err
		}
		if len(resp.ConfigurationSettings) == 0 {
			return fmt.Errorf("Elastic Beanstalk Configurations not found.")
		}

		*config = *resp.ConfigurationSettings[0]

		return nil
	}
}

func testAccConfigurationTemplateConfig_basic(r string) string {
	return fmt.Sprintf(`
resource "aws_elastic_beanstalk_application" "tftest" {
  name        = "tf-test-%s"
  description = "tf-test-desc-%s"
}

resource "aws_elastic_beanstalk_configuration_template" "tf_template" {
  name                = "tf-test-template-config"
  application         = aws_elastic_beanstalk_application.tftest.name
  solution_stack_name = "64bit Amazon Linux running Python"
}
`, r, r)
}

func testAccConfigurationTemplateConfig_vpc(name string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "tf_b_test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-elastic-beanstalk-cfg-tpl-vpc"
  }
}

resource "aws_subnet" "main" {
  vpc_id     = aws_vpc.tf_b_test.id
  cidr_block = "10.0.0.0/24"

  tags = {
    Name = "tf-acc-elastic-beanstalk-cfg-tpl-vpc"
  }
}

resource "aws_elastic_beanstalk_application" "tftest" {
  name        = "tf-test-%s"
  description = "tf-test-desc"
}

resource "aws_elastic_beanstalk_configuration_template" "tf_template" {
  name        = "tf-test-%s"
  application = aws_elastic_beanstalk_application.tftest.name

  solution_stack_name = "64bit Amazon Linux running Python"

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.tf_b_test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.main.id
  }
}
`, name, name)
}

func testAccConfigurationTemplateConfig_setting(name string) string {
	return fmt.Sprintf(`
resource "aws_elastic_beanstalk_application" "tftest" {
  name        = "tf-test-%s"
  description = "tf-test-desc"
}

resource "aws_elastic_beanstalk_configuration_template" "tf_template" {
  name        = "tf-test-%s"
  application = aws_elastic_beanstalk_application.tftest.name

  solution_stack_name = "64bit Amazon Linux running Python"

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "InstanceType"
    value     = "m1.small"
  }
}
`, name, name)
}
