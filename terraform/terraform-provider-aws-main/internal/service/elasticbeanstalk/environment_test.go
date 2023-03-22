package elasticbeanstalk_test

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfelasticbeanstalk "github.com/hashicorp/terraform-provider-aws/internal/service/elasticbeanstalk"
)

func TestAccElasticBeanstalkEnvironment_BeanstalkEnv_basic(t *testing.T) {
	var app elasticbeanstalk.EnvironmentDescription

	resourceName := "aws_elastic_beanstalk_environment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	beanstalkAsgNameRegexp := regexp.MustCompile("awseb.+?AutoScalingGroup[^,]+")
	beanstalkElbNameRegexp := regexp.MustCompile("awseb.+?EBLoa[^,]+")
	beanstalkInstancesNameRegexp := regexp.MustCompile("i-([0-9a-fA-F]{8}|[0-9a-fA-F]{17})")
	beanstalkLcNameRegexp := regexp.MustCompile("awseb.+?AutoScalingLaunch[^,]+")
	beanstalkEndpointURL := regexp.MustCompile("awseb.+?EBLoa[^,].+?elb.amazonaws.com")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "elasticbeanstalk", fmt.Sprintf("environment/%s/%s", rName, rName)),
					resource.TestMatchResourceAttr(resourceName, "autoscaling_groups.0", beanstalkAsgNameRegexp),
					resource.TestMatchResourceAttr(resourceName, "endpoint_url", beanstalkEndpointURL),
					resource.TestMatchResourceAttr(resourceName, "instances.0", beanstalkInstancesNameRegexp),
					resource.TestMatchResourceAttr(resourceName, "launch_configurations.0", beanstalkLcNameRegexp),
					resource.TestMatchResourceAttr(resourceName, "load_balancers.0", beanstalkElbNameRegexp),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"setting",
					"wait_for_ready_timeout",
				},
			},
		},
	})
}

func TestAccElasticBeanstalkEnvironment_BeanstalkEnv_tier(t *testing.T) {
	var app elasticbeanstalk.EnvironmentDescription
	beanstalkQueuesNameRegexp := regexp.MustCompile("https://sqs.+?awseb[^,]+")

	resourceName := "aws_elastic_beanstalk_environment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_worker(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentTier(resourceName, &app),
					resource.TestMatchResourceAttr(resourceName, "queues.0", beanstalkQueuesNameRegexp),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"setting",
					"wait_for_ready_timeout",
				},
			},
		},
	})
}

func TestAccElasticBeanstalkEnvironment_BeanstalkEnvCNAME_prefix(t *testing.T) {
	var app elasticbeanstalk.EnvironmentDescription

	resourceName := "aws_elastic_beanstalk_environment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	beanstalkCnameRegexp := regexp.MustCompile("^" + rName + ".+?elasticbeanstalk.com$")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_cnamePrefix(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					resource.TestMatchResourceAttr(resourceName, "cname", beanstalkCnameRegexp),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"setting",
					"wait_for_ready_timeout",
				},
			},
		},
	})
}

func TestAccElasticBeanstalkEnvironment_beanstalkEnv(t *testing.T) {
	var app elasticbeanstalk.EnvironmentDescription

	resourceName := "aws_elastic_beanstalk_environment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_template(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					testAccCheckEnvironmentConfigValue(resourceName, "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"setting",
					"template_name",
					"wait_for_ready_timeout",
				},
			},
			{
				Config: testAccEnvironmentConfig_template(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					testAccCheckEnvironmentConfigValue(resourceName, "2"),
				),
			},
			{
				Config: testAccEnvironmentConfig_template(rName, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					testAccCheckEnvironmentConfigValue(resourceName, "3"),
				),
			},
		},
	})
}

func TestAccElasticBeanstalkEnvironment_BeanstalkEnv_resource(t *testing.T) {
	var app elasticbeanstalk.EnvironmentDescription

	resourceName := "aws_elastic_beanstalk_environment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_resourceOptionSetting(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"setting",
					"wait_for_ready_timeout",
				},
			},
		},
	})
}

func TestAccElasticBeanstalkEnvironment_BeanstalkEnv_tags(t *testing.T) {
	var app elasticbeanstalk.EnvironmentDescription

	resourceName := "aws_elastic_beanstalk_environment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_tagsTemplate(rName, "test1", "test2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					testAccCheckEnvironmentTagsMatch(&app, map[string]string{"firstTag": "test1", "secondTag": "test2"}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"setting",
					"wait_for_ready_timeout",
				},
			},
			{
				Config: testAccEnvironmentConfig_tagsTemplate(rName, "test2", "test1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					testAccCheckEnvironmentTagsMatch(&app, map[string]string{"firstTag": "test2", "secondTag": "test1"}),
				),
			},
			{
				Config: testAccEnvironmentConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					testAccCheckEnvironmentTagsMatch(&app, map[string]string{}),
				),
			},
		},
	})
}

func TestAccElasticBeanstalkEnvironment_BeanstalkEnvTemplate_change(t *testing.T) {
	var app elasticbeanstalk.EnvironmentDescription

	resourceName := "aws_elastic_beanstalk_environment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_templateChangeStack(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
				),
			},
			{
				Config: testAccEnvironmentConfig_templateChangeTemp(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
				),
			},
			{
				Config: testAccEnvironmentConfig_templateChangeStack(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
				),
			},
		},
	})
}

func TestAccElasticBeanstalkEnvironment_BeanstalkEnvSettings_update(t *testing.T) {
	var app elasticbeanstalk.EnvironmentDescription

	resourceName := "aws_elastic_beanstalk_environment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					testAccVerifyConfig(&app, []string{}),
				),
			},
			{
				Config: testAccEnvironmentConfig_settings(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					testAccVerifyConfig(&app, []string{"ENV_STATIC", "ENV_UPDATE"}),
				),
			},
			{
				Config: testAccEnvironmentConfig_settings(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					testAccVerifyConfig(&app, []string{"ENV_STATIC", "ENV_UPDATE"}),
				),
			},
			{
				Config: testAccEnvironmentConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					testAccVerifyConfig(&app, []string{}),
				),
			},
		},
	})
}

func TestAccElasticBeanstalkEnvironment_BeanstalkEnvVersion_label(t *testing.T) {
	var app elasticbeanstalk.EnvironmentDescription

	resourceName := "aws_elastic_beanstalk_environment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_applicationVersion(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationVersionDeployed(resourceName, &app),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"setting",
					"wait_for_ready_timeout",
				},
			},
			{
				Config: testAccEnvironmentConfig_applicationVersionUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationVersionDeployed(resourceName, &app),
				),
			},
		},
	})
}

func TestAccElasticBeanstalkEnvironment_BeanstalkEnv_settingWithJSONValue(t *testing.T) {
	var app elasticbeanstalk.EnvironmentDescription

	resourceName := "aws_elastic_beanstalk_environment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	publicKey, _, err := sdkacctest.RandSSHKeyPair(acctest.DefaultEmailAddress)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_settingJSONValue(rName, publicKey, acctest.DefaultEmailAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"setting",
					"wait_for_ready_timeout",
				},
			},
		},
	})
}

func TestAccElasticBeanstalkEnvironment_BeanstalkEnv_platformARN(t *testing.T) {
	var app elasticbeanstalk.EnvironmentDescription

	resourceName := "aws_elastic_beanstalk_environment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_platformARN(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &app),
					acctest.CheckResourceAttrRegionalARNNoAccount(resourceName, "platform_arn", "elasticbeanstalk", "platform/Python 3.6 running on 64bit Amazon Linux/2.9.6"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"setting",
					"wait_for_ready_timeout",
				},
			},
		},
	})
}

func testAccVerifyConfig(env *elasticbeanstalk.EnvironmentDescription, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if env == nil {
			return fmt.Errorf("Nil environment in testAccVerifyConfig")
		}
		conn := acctest.Provider.Meta().(*conns.AWSClient).ElasticBeanstalkConn

		resp, err := conn.DescribeConfigurationSettings(&elasticbeanstalk.DescribeConfigurationSettingsInput{
			ApplicationName: env.ApplicationName,
			EnvironmentName: env.EnvironmentName,
		})

		if err != nil {
			return fmt.Errorf("Error describing config settings in testAccVerifyConfig: %s", err)
		}

		// should only be 1 environment
		if len(resp.ConfigurationSettings) != 1 {
			return fmt.Errorf("Expected only 1 set of Configuration Settings in testAccVerifyConfig, got (%d)", len(resp.ConfigurationSettings))
		}

		cs := resp.ConfigurationSettings[0]

		var foundEnvs []string
		testStrings := []string{"ENV_STATIC", "ENV_UPDATE"}
		for _, os := range cs.OptionSettings {
			for _, k := range testStrings {
				if *os.OptionName == k {
					foundEnvs = append(foundEnvs, k)
				}
			}
		}

		// if expected is zero, then we should not have found any of the predefined
		// env vars
		if len(expected) == 0 {
			if len(foundEnvs) > 0 {
				return fmt.Errorf("Found configs we should not have: %#v", foundEnvs)
			}
			return nil
		}

		sort.Strings(testStrings)
		sort.Strings(expected)
		if !reflect.DeepEqual(testStrings, expected) {
			return fmt.Errorf("error matching strings, expected:\n\n%#v\n\ngot:\n\n%#v", testStrings, foundEnvs)
		}

		return nil
	}
}

func testAccCheckEnvironmentDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ElasticBeanstalkConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_elastic_beanstalk_environment" {
			continue
		}

		// Try to find the environment
		describeBeanstalkEnvOpts := &elasticbeanstalk.DescribeEnvironmentsInput{
			EnvironmentIds: []*string{aws.String(rs.Primary.ID)},
		}
		resp, err := conn.DescribeEnvironments(describeBeanstalkEnvOpts)
		if err == nil {
			switch {
			case len(resp.Environments) > 1:
				return fmt.Errorf("error %d environments match, expected 1", len(resp.Environments))
			case len(resp.Environments) == 1:
				if *resp.Environments[0].Status == "Terminated" {
					return nil
				}
				return fmt.Errorf("Elastic Beanstalk ENV still exists")
			default:
				return nil
			}
		}

		if !tfawserr.ErrCodeEquals(err, "InvalidBeanstalkEnvID.NotFound") {
			return err
		}
	}

	return nil
}

func testAccCheckEnvironmentExists(n string, app *elasticbeanstalk.EnvironmentDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Elastic Beanstalk ENV is not set")
		}

		env, err := describeEnvironment(acctest.Provider.Meta().(*conns.AWSClient).ElasticBeanstalkConn, aws.String(rs.Primary.ID))
		if err != nil {
			return err
		}

		*app = *env

		return nil
	}
}

func testAccCheckEnvironmentTier(n string, app *elasticbeanstalk.EnvironmentDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Elastic Beanstalk ENV is not set")
		}

		env, err := describeEnvironment(acctest.Provider.Meta().(*conns.AWSClient).ElasticBeanstalkConn, aws.String(rs.Primary.ID))
		if err != nil {
			return err
		}
		if *env.Tier.Name != "Worker" {
			return fmt.Errorf("Beanstalk Environment tier is %s, expected Worker", *env.Tier.Name)
		}

		*app = *env

		return nil
	}
}

func testAccCheckEnvironmentConfigValue(n string, expectedValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ElasticBeanstalkConn

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Elastic Beanstalk ENV is not set")
		}

		resp, err := conn.DescribeConfigurationOptions(&elasticbeanstalk.DescribeConfigurationOptionsInput{
			ApplicationName: aws.String(rs.Primary.Attributes["application"]),
			EnvironmentName: aws.String(rs.Primary.Attributes["name"]),
			Options: []*elasticbeanstalk.OptionSpecification{
				{
					Namespace:  aws.String("aws:elasticbeanstalk:application:environment"),
					OptionName: aws.String("TEMPLATE"),
				},
			},
		})
		if err != nil {
			return err
		}

		if len(resp.Options) != 1 {
			return fmt.Errorf("Found %d options, expected 1.", len(resp.Options))
		}

		log.Printf("[DEBUG] %d Elastic Beanstalk Option values returned.", len(resp.Options[0].ValueOptions))

		for _, value := range resp.Options[0].ValueOptions {
			if *value != expectedValue {
				return fmt.Errorf("Option setting value: %s. Expected %s", *value, expectedValue)
			}
		}

		return nil
	}
}

func testAccCheckEnvironmentTagsMatch(env *elasticbeanstalk.EnvironmentDescription, expectedValue map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if env == nil {
			return fmt.Errorf("Nil environment in testAccCheckEnvironmentTagsMatch")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ElasticBeanstalkConn

		tags, err := conn.ListTagsForResource(&elasticbeanstalk.ListTagsForResourceInput{
			ResourceArn: env.EnvironmentArn,
		})

		if err != nil {
			return err
		}

		foundTags := tfelasticbeanstalk.KeyValueTags(tags.ResourceTags).IgnoreElasticbeanstalk().Map()

		if !reflect.DeepEqual(foundTags, expectedValue) {
			return fmt.Errorf("Tag value: %s.  Expected %s", foundTags, expectedValue)
		}

		return nil
	}
}

func testAccCheckApplicationVersionDeployed(n string, app *elasticbeanstalk.EnvironmentDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Elastic Beanstalk ENV is not set")
		}

		env, err := describeEnvironment(acctest.Provider.Meta().(*conns.AWSClient).ElasticBeanstalkConn, aws.String(rs.Primary.ID))
		if err != nil {
			return err
		}

		if *env.VersionLabel != rs.Primary.Attributes["version_label"] {
			return fmt.Errorf("Elastic Beanstalk version deployed %s. Expected %s", *env.VersionLabel, rs.Primary.Attributes["version_label"])
		}

		*app = *env

		return nil
	}
}

func describeEnvironment(conn *elasticbeanstalk.ElasticBeanstalk,
	envID *string) (*elasticbeanstalk.EnvironmentDescription, error) {
	describeBeanstalkEnvOpts := &elasticbeanstalk.DescribeEnvironmentsInput{
		EnvironmentIds: []*string{envID},
	}

	log.Printf("[DEBUG] Elastic Beanstalk Environment TEST describe opts: %s", describeBeanstalkEnvOpts)

	resp, err := conn.DescribeEnvironments(describeBeanstalkEnvOpts)
	if err != nil {
		return &elasticbeanstalk.EnvironmentDescription{}, err
	}
	if len(resp.Environments) == 0 {
		return &elasticbeanstalk.EnvironmentDescription{}, fmt.Errorf("Elastic Beanstalk ENV not found")
	}
	if len(resp.Environments) > 1 {
		return &elasticbeanstalk.EnvironmentDescription{}, fmt.Errorf("found %d environments, expected 1", len(resp.Environments))
	}
	return resp.Environments[0], nil
}

func testAccEnvironmentConfig_base(rName string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  # Default instance type of t2.micro is not available in this Availability Zone
  # The failure will occur during Elastic Beanstalk CloudFormation Template handling
  # after waiting upwards of one hour to initialize the Auto Scaling Group.
  exclude_zone_ids = ["usw2-az4"]
  state            = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

data "aws_elastic_beanstalk_solution_stack" "test" {
  most_recent = true
  name_regex  = "64bit Amazon Linux .* running Python .*"
}

data "aws_partition" "current" {}

data "aws_region" "current" {}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-elastic-beanstalk-env-vpc"
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id
}

resource "aws_route" "test" {
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.test.id
  route_table_id         = aws_vpc.test.main_route_table_id
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = "tf-acc-elastic-beanstalk-env-vpc"
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id
}

resource "aws_elastic_beanstalk_application" "test" {
  description = "tf-test-desc"
  name        = %[1]q
}

# Create custom service role per test to remove dependency on
# Service-Linked Role existing.
resource "aws_iam_role" "service_role" {
  name = "%[1]s-service"
  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Condition = {
        StringEquals = {
          "sts:ExternalId" = "elasticbeanstalk"
        }
      }
      Effect = "Allow"
      Principal = {
        Service = "elasticbeanstalk.${data.aws_partition.current.dns_suffix}"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "service_role-AWSElasticBeanstalkEnhancedHealth" {
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AWSElasticBeanstalkEnhancedHealth"
  role       = aws_iam_role.service_role.id
}

resource "aws_iam_role_policy_attachment" "service_role-AWSElasticBeanstalkService" {
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AWSElasticBeanstalkService"
  role       = aws_iam_role.service_role.id
}

# Amazon Linux 2 environments require IAM Instance Profile.
resource "aws_iam_role" "instance_profile" {
  name = "%[1]s-instance"
  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.${data.aws_partition.current.dns_suffix}"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "instance_profile-AWSElasticBeanstalkWebTier" {
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AWSElasticBeanstalkWebTier"
  role       = aws_iam_role.instance_profile.id
}

resource "aws_iam_role_policy_attachment" "instance_profile-AWSElasticBeanstalkWorkerTier" {
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AWSElasticBeanstalkWorkerTier"
  role       = aws_iam_role.instance_profile.id
}

# Since the IAM Instance Profile is required anyways, also use it to
# ensure IAM Role permissions for both roles are attached.
resource "aws_iam_instance_profile" "test" {
  depends_on = [
    aws_iam_role_policy_attachment.instance_profile-AWSElasticBeanstalkWebTier,
    aws_iam_role_policy_attachment.instance_profile-AWSElasticBeanstalkWorkerTier,
    aws_iam_role_policy_attachment.service_role-AWSElasticBeanstalkEnhancedHealth,
    aws_iam_role_policy_attachment.service_role-AWSElasticBeanstalkService,
  ]

  name = %[1]q
  role = aws_iam_role.instance_profile.name
}
`, rName)
}

func testAccEnvironmentConfig_basic(rName string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_elastic_beanstalk_environment" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }
}
`, rName)
}

func testAccEnvironmentConfig_platformARN(rName string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_elastic_beanstalk_environment" "test" {
  application  = aws_elastic_beanstalk_application.test.name
  name         = %[1]q
  platform_arn = "arn:${data.aws_partition.current.partition}:elasticbeanstalk:${data.aws_region.current.name}::platform/Python 3.6 running on 64bit Amazon Linux/2.9.6"

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }
}
`, rName)
}

func testAccEnvironmentConfig_settings(rName string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_elastic_beanstalk_environment" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:application:environment"
    name      = "ENV_STATIC"
    value     = "true"
  }

  setting {
    namespace = "aws:elasticbeanstalk:application:environment"
    name      = "ENV_UPDATE"
    value     = "true"
  }

  setting {
    namespace = "aws:elasticbeanstalk:application:environment"
    name      = "ENV_REMOVE"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:scheduledaction"
    resource  = "ScheduledAction01"
    name      = "MinSize"
    value     = 2
  }

  setting {
    namespace = "aws:autoscaling:scheduledaction"
    resource  = "ScheduledAction01"
    name      = "MaxSize"
    value     = 3
  }

  setting {
    namespace = "aws:autoscaling:scheduledaction"
    resource  = "ScheduledAction01"
    name      = "StartTime"
    value     = "2016-07-28T04:07:02Z"
  }
}
`, rName)
}

func testAccEnvironmentConfig_worker(rName string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_elastic_beanstalk_environment" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name
  tier                = "Worker"

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }
}
`, rName)
}

func testAccEnvironmentConfig_cnamePrefix(rName string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_elastic_beanstalk_environment" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  cname_prefix        = %[1]q
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }
}
`, rName)
}

func testAccEnvironmentConfig_template(rName string, cfgTplValue int) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_elastic_beanstalk_environment" "test" {
  application   = aws_elastic_beanstalk_application.test.name
  name          = %[1]q
  template_name = aws_elastic_beanstalk_configuration_template.test.name
}

resource "aws_elastic_beanstalk_configuration_template" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:application:environment"
    name      = "TEMPLATE"
    value     = %[2]d
  }
}
`, rName, cfgTplValue)
}

func testAccEnvironmentConfig_resourceOptionSetting(rName string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_elastic_beanstalk_environment" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }

  setting {
    namespace = "aws:autoscaling:scheduledaction"
    resource  = "ScheduledAction01"
    name      = "MinSize"
    value     = "2"
  }

  setting {
    namespace = "aws:autoscaling:scheduledaction"
    resource  = "ScheduledAction01"
    name      = "MaxSize"
    value     = "6"
  }

  setting {
    namespace = "aws:autoscaling:scheduledaction"
    resource  = "ScheduledAction01"
    name      = "Recurrence"
    value     = "0 8 * * *"
  }
}
`, rName)
}

func testAccEnvironmentConfig_tagsTemplate(rName, firstTag, secondTag string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_elastic_beanstalk_environment" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }

  tags = {
    firstTag  = %[2]q
    secondTag = %[3]q
  }
}
`, rName, firstTag, secondTag)
}

func testAccEnvironmentConfig_templateChangeStack(rName string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_elastic_beanstalk_environment" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }
}

resource "aws_elastic_beanstalk_configuration_template" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name
}
`, rName)
}

func testAccEnvironmentConfig_templateChangeTemp(rName string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_elastic_beanstalk_environment" "test" {
  application   = aws_elastic_beanstalk_application.test.name
  name          = %[1]q
  template_name = aws_elastic_beanstalk_configuration_template.test.name

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }
}

resource "aws_elastic_beanstalk_configuration_template" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name
}
`, rName)
}

func testAccEnvironmentConfig_applicationVersion(rName string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_object" "test" {
  bucket = aws_s3_bucket.test.id
  key    = "python-v1.zip"
  source = "test-fixtures/python-v1.zip"
}

resource "aws_elastic_beanstalk_application_version" "test" {
  application = aws_elastic_beanstalk_application.test.name
  bucket      = aws_s3_bucket.test.id
  key         = aws_s3_object.test.id
  name        = "%[1]s-1"
}

resource "aws_elastic_beanstalk_environment" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name
  version_label       = aws_elastic_beanstalk_application_version.test.name

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }
}
`, rName)
}

func testAccEnvironmentConfig_applicationVersionUpdate(rName string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_object" "test" {
  bucket = aws_s3_bucket.test.id
  key    = "python-v1.zip"
  source = "test-fixtures/python-v1.zip"
}

resource "aws_elastic_beanstalk_application_version" "test" {
  application = aws_elastic_beanstalk_application.test.name
  bucket      = aws_s3_bucket.test.id
  key         = aws_s3_object.test.id
  name        = "%[1]s-2"
}

resource "aws_elastic_beanstalk_environment" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name
  version_label       = aws_elastic_beanstalk_application_version.test.name

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }
}
`, rName)
}

func testAccEnvironmentConfig_settingJSONValue(rName, publicKey, email string) string {
	return testAccEnvironmentConfig_base(rName) + fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name = %[1]q
}

resource "aws_key_pair" "test" {
  key_name   = %[1]q
  public_key = %[2]q
}

resource "aws_elastic_beanstalk_environment" "test" {
  application         = aws_elastic_beanstalk_application.test.name
  name                = %[1]q
  solution_stack_name = data.aws_elastic_beanstalk_solution_stack.test.name
  tier                = "Worker"

  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = aws_vpc.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = aws_subnet.test.id
  }

  setting {
    namespace = "aws:ec2:vpc"
    name      = "AssociatePublicIpAddress"
    value     = "true"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "SecurityGroups"
    value     = aws_security_group.test.id
  }

  setting {
    namespace = "aws:elasticbeanstalk:command"
    name      = "BatchSize"
    value     = "30"
  }

  setting {
    namespace = "aws:elasticbeanstalk:command"
    name      = "BatchSizeType"
    value     = "Percentage"
  }

  setting {
    namespace = "aws:elasticbeanstalk:command"
    name      = "DeploymentPolicy"
    value     = "Rolling"
  }

  setting {
    namespace = "aws:elasticbeanstalk:sns:topics"
    name      = "Notification Endpoint"
    value     = %[3]q
  }

  setting {
    namespace = "aws:elasticbeanstalk:sqsd"
    name      = "ErrorVisibilityTimeout"
    value     = "2"
  }

  setting {
    namespace = "aws:elasticbeanstalk:sqsd"
    name      = "HttpPath"
    value     = "/event-message"
  }

  setting {
    namespace = "aws:elasticbeanstalk:sqsd"
    name      = "WorkerQueueURL"
    value     = aws_sqs_queue.test.id
  }

  setting {
    namespace = "aws:elasticbeanstalk:sqsd"
    name      = "VisibilityTimeout"
    value     = "300"
  }

  setting {
    namespace = "aws:elasticbeanstalk:sqsd"
    name      = "HttpConnections"
    value     = "10"
  }

  setting {
    namespace = "aws:elasticbeanstalk:sqsd"
    name      = "InactivityTimeout"
    value     = "299"
  }

  setting {
    namespace = "aws:elasticbeanstalk:sqsd"
    name      = "MimeType"
    value     = "application/json"
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "EnvironmentType"
    value     = "LoadBalanced"
  }

  setting {
    namespace = "aws:elasticbeanstalk:application"
    name      = "Application Healthcheck URL"
    value     = "/health"
  }

  setting {
    namespace = "aws:elasticbeanstalk:healthreporting:system"
    name      = "SystemType"
    value     = "enhanced"
  }

  setting {
    namespace = "aws:elasticbeanstalk:healthreporting:system"
    name      = "HealthCheckSuccessThreshold"
    value     = "Ok"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.test.name
  }

  setting {
    namespace = "aws:elasticbeanstalk:environment"
    name      = "ServiceRole"
    value     = aws_iam_role.service_role.name
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "InstanceType"
    value     = "t2.micro"
  }

  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "EC2KeyName"
    value     = aws_key_pair.test.key_name
  }

  setting {
    namespace = "aws:autoscaling:updatepolicy:rollingupdate"
    name      = "RollingUpdateEnabled"
    value     = "false"
  }

  setting {
    namespace = "aws:elasticbeanstalk:healthreporting:system"
    name      = "ConfigDocument"

    value = <<EOF
{
	"Version": 1,
	"CloudWatchMetrics": {
		"Instance": {
			"ApplicationRequestsTotal": 60
		},
		"Environment": {
			"ApplicationRequests5xx": 60,
			"ApplicationRequests4xx": 60,
			"ApplicationRequests2xx": 60
		}
	}
}
EOF
  }
}
`, rName, publicKey, email)
}
