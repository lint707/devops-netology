package mwaa_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mwaa"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfmwaa "github.com/hashicorp/terraform-provider-aws/internal/service/mwaa"
)

func TestAccMWAAEnvironment_basic(t *testing.T) {
	var environment mwaa.GetEnvironmentOutput

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_mwaa_environment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, mwaa.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttrSet(resourceName, "airflow_version"),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "airflow", "environment/"+rName),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "dag_s3_path", "dags/"),
					resource.TestCheckResourceAttr(resourceName, "environment_class", "mw1.small"),
					acctest.CheckResourceAttrGlobalARN(resourceName, "execution_role_arn", "iam", "role/service-role/"+rName),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.0.log_level", "INFO"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.0.log_level", "INFO"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_configuration.0.task_logs.0.cloud_watch_log_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.0.log_level", "INFO"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.0.log_level", "INFO"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.0.log_level", "INFO"),
					resource.TestCheckResourceAttr(resourceName, "max_workers", "10"),
					resource.TestCheckResourceAttr(resourceName, "min_workers", "1"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "network_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_configuration.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_configuration.0.subnet_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "schedulers", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "service_role_arn"),
					acctest.CheckResourceAttrGlobalARNNoAccount(resourceName, "source_bucket_arn", "s3", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "webserver_access_mode", mwaa.WebserverAccessModePrivateOnly),
					resource.TestCheckResourceAttrSet(resourceName, "webserver_url"),
					resource.TestCheckResourceAttrSet(resourceName, "weekly_maintenance_window_start"),
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

func TestAccMWAAEnvironment_disappears(t *testing.T) {
	var environment mwaa.GetEnvironmentOutput

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_mwaa_environment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, mwaa.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &environment),
					acctest.CheckResourceDisappears(acctest.Provider, tfmwaa.ResourceEnvironment(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccMWAAEnvironment_airflowOptions(t *testing.T) {
	var environment mwaa.GetEnvironmentOutput

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_mwaa_environment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, mwaa.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_airflowOptions(rName, "1", "16"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "airflow_configuration_options.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "airflow_configuration_options.core.default_task_retries", "1"),
					resource.TestCheckResourceAttr(resourceName, "airflow_configuration_options.core.parallelism", "16"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEnvironmentConfig_airflowOptions(rName, "2", "32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "airflow_configuration_options.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "airflow_configuration_options.core.default_task_retries", "2"),
					resource.TestCheckResourceAttr(resourceName, "airflow_configuration_options.core.parallelism", "32"),
				),
			},
			{
				Config: testAccEnvironmentConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "airflow_configuration_options.%", "0"),
				),
			},
		},
	})
}

func TestAccMWAAEnvironment_log(t *testing.T) {
	var environment mwaa.GetEnvironmentOutput

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_mwaa_environment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, mwaa.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_logging(rName, "true", mwaa.LoggingLevelCritical),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.#", "1"),

					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_configuration.0.dag_processing_logs.0.cloud_watch_log_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.0.log_level", mwaa.LoggingLevelCritical),

					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_configuration.0.scheduler_logs.0.cloud_watch_log_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.0.log_level", mwaa.LoggingLevelCritical),

					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_configuration.0.task_logs.0.cloud_watch_log_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.0.log_level", mwaa.LoggingLevelCritical),

					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_configuration.0.webserver_logs.0.cloud_watch_log_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.0.log_level", mwaa.LoggingLevelCritical),

					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_configuration.0.worker_logs.0.cloud_watch_log_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.0.log_level", mwaa.LoggingLevelCritical),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEnvironmentConfig_logging(rName, "false", mwaa.LoggingLevelInfo),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.#", "1"),

					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.0.cloud_watch_log_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.0.log_level", mwaa.LoggingLevelInfo),

					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.0.cloud_watch_log_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.0.log_level", mwaa.LoggingLevelInfo),

					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.0.cloud_watch_log_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.0.log_level", mwaa.LoggingLevelInfo),

					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.0.cloud_watch_log_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.0.log_level", mwaa.LoggingLevelInfo),

					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.0.cloud_watch_log_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.0.log_level", mwaa.LoggingLevelInfo),
				),
			},
		},
	})
}

func TestAccMWAAEnvironment_full(t *testing.T) {
	var environment mwaa.GetEnvironmentOutput

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_mwaa_environment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, mwaa.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_full(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "airflow_configuration_options.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "airflow_configuration_options.core.default_task_retries", "1"),
					resource.TestCheckResourceAttr(resourceName, "airflow_configuration_options.core.parallelism", "16"),
					resource.TestCheckResourceAttr(resourceName, "airflow_version", "1.10.12"),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "airflow", "environment/"+rName),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "dag_s3_path", "dags/"),
					resource.TestCheckResourceAttr(resourceName, "environment_class", "mw1.medium"),
					acctest.CheckResourceAttrGlobalARN(resourceName, "execution_role_arn", "iam", "role/service-role/"+rName),
					resource.TestCheckResourceAttrSet(resourceName, "kms_key"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_configuration.0.dag_processing_logs.0.cloud_watch_log_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.dag_processing_logs.0.log_level", "INFO"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_configuration.0.scheduler_logs.0.cloud_watch_log_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.scheduler_logs.0.log_level", "WARNING"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_configuration.0.task_logs.0.cloud_watch_log_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.task_logs.0.log_level", "ERROR"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_configuration.0.webserver_logs.0.cloud_watch_log_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.webserver_logs.0.log_level", "CRITICAL"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_configuration.0.worker_logs.0.cloud_watch_log_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "logging_configuration.0.worker_logs.0.log_level", "WARNING"),
					resource.TestCheckResourceAttr(resourceName, "max_workers", "20"),
					resource.TestCheckResourceAttr(resourceName, "min_workers", "15"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "network_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_configuration.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_configuration.0.subnet_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "plugins_s3_path", "plugins.zip"),
					resource.TestCheckResourceAttr(resourceName, "requirements_s3_path", "requirements.txt"),
					resource.TestCheckResourceAttr(resourceName, "schedulers", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "service_role_arn"),
					acctest.CheckResourceAttrGlobalARNNoAccount(resourceName, "source_bucket_arn", "s3", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "webserver_access_mode", mwaa.WebserverAccessModePublicOnly),
					resource.TestCheckResourceAttrSet(resourceName, "webserver_url"),
					resource.TestCheckResourceAttr(resourceName, "weekly_maintenance_window_start", "SAT:03:00"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "production"),
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

func TestAccMWAAEnvironment_pluginsS3ObjectVersion(t *testing.T) {
	var environment mwaa.GetEnvironmentOutput

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_mwaa_environment.test"
	s3ObjectResourceName := "aws_s3_object.plugins"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, mwaa.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_pluginsS3ObjectVersion(rName, "test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttrPair(resourceName, "plugins_s3_object_version", s3ObjectResourceName, "version_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEnvironmentConfig_pluginsS3ObjectVersion(rName, "test-updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttrPair(resourceName, "plugins_s3_object_version", s3ObjectResourceName, "version_id"),
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

func testAccCheckEnvironmentExists(resourceName string, environment *mwaa.GetEnvironmentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MWAA Environment ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).MWAAConn
		resp, err := conn.GetEnvironment(&mwaa.GetEnvironmentInput{
			Name: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("Error getting MWAA Environment: %s", err.Error())
		}

		*environment = *resp

		return nil
	}
}

func testAccCheckEnvironmentDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).MWAAConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_mwaa_environment" {
			continue
		}

		_, err := conn.GetEnvironment(&mwaa.GetEnvironmentInput{
			Name: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if tfawserr.ErrCodeEquals(err, mwaa.ErrCodeResourceNotFoundException) {
				continue
			}
			return err
		}

		return fmt.Errorf("Expected MWAA Environment to be destroyed, %s found", rs.Primary.ID)
	}

	return nil
}

func testAccEnvironmentBase(rName string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

data "aws_partition" "current" {}

resource "aws_vpc" "test" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.test.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.test.id
  }
}

resource "aws_main_route_table_association" "test" {
  route_table_id = aws_route_table.public.id
  vpc_id         = aws_vpc.test.id
}

resource "aws_subnet" "private" {
  count = 2

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(aws_vpc.test.cidr_block, 8, count.index + 2)
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = "%[1]s-private-${count.index}"
  }
}

resource "aws_eip" "private" {
  count = 2

  vpc = true
}

resource "aws_nat_gateway" "private" {
  count = 2

  allocation_id = aws_eip.private[count.index].id
  subnet_id     = aws_subnet.public[count.index].id
}

resource "aws_route_table" "private" {
  count = 2

  vpc_id = aws_vpc.test.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.private[count.index].id
  }
}

resource "aws_route_table_association" "private" {
  count = 2

  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}

resource "aws_subnet" "public" {
  count = 2

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(aws_vpc.test.cidr_block, 8, count.index)
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = "%[1]s-public-${count.index}"
  }
}

resource "aws_security_group" "test" {
  vpc_id = aws_vpc.test.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_acl" "test" {
  bucket = aws_s3_bucket.test.id
  acl    = "private"
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_public_access_block" "test" {
  bucket = aws_s3_bucket.test.bucket

  block_public_acls   = true
  block_public_policy = true
}

resource "aws_s3_object" "dags" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket       = aws_s3_bucket.test.id
  acl          = "private"
  key          = "dags/"
  content_type = "application/x-directory"
}

resource "aws_iam_role" "test" {
  name = %[1]q
  path = "/service-role/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "airflow.${data.aws_partition.current.dns_suffix}",
          "airflow-env.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "test" {
  role = aws_iam_role.test.name

  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "*",
            "Resource": "*"
        }
    ]
}
POLICY
}


`, rName)
}

func testAccEnvironmentConfig_basic(rName string) string {
	return testAccEnvironmentBase(rName) + fmt.Sprintf(`
resource "aws_mwaa_environment" "test" {
  dag_s3_path        = aws_s3_object.dags.key
  execution_role_arn = aws_iam_role.test.arn
  name               = %[1]q

  network_configuration {
    security_group_ids = [aws_security_group.test.id]
    subnet_ids         = aws_subnet.private[*].id
  }

  source_bucket_arn = aws_s3_bucket.test.arn
}
`, rName)
}

func testAccEnvironmentConfig_airflowOptions(rName, retries, parallelism string) string {
	return testAccEnvironmentBase(rName) + fmt.Sprintf(`
resource "aws_mwaa_environment" "test" {
  airflow_configuration_options = {
    "core.default_task_retries" = %[2]q
    "core.parallelism"          = %[3]q
  }

  dag_s3_path        = aws_s3_object.dags.key
  execution_role_arn = aws_iam_role.test.arn
  name               = %[1]q

  network_configuration {
    security_group_ids = [aws_security_group.test.id]
    subnet_ids         = aws_subnet.private[*].id
  }

  source_bucket_arn = aws_s3_bucket.test.arn
}
`, rName, retries, parallelism)
}

func testAccEnvironmentConfig_logging(rName, logEnabled, logLevel string) string {
	return testAccEnvironmentBase(rName) + fmt.Sprintf(`
resource "aws_mwaa_environment" "test" {
  dag_s3_path        = aws_s3_object.dags.key
  execution_role_arn = aws_iam_role.test.arn

  logging_configuration {
    dag_processing_logs {
      enabled   = %[2]s
      log_level = %[3]q
    }

    scheduler_logs {
      enabled   = %[2]s
      log_level = %[3]q
    }

    task_logs {
      enabled   = %[2]s
      log_level = %[3]q
    }

    webserver_logs {
      enabled   = %[2]s
      log_level = %[3]q
    }

    worker_logs {
      enabled   = %[2]s
      log_level = %[3]q
    }
  }

  name = %[1]q

  network_configuration {
    security_group_ids = [aws_security_group.test.id]
    subnet_ids         = aws_subnet.private[*].id
  }

  source_bucket_arn = aws_s3_bucket.test.arn
}
`, rName, logEnabled, logLevel)
}

func testAccEnvironmentConfig_full(rName string) string {
	return testAccEnvironmentBase(rName) + fmt.Sprintf(`
resource "aws_mwaa_environment" "test" {
  airflow_configuration_options = {
    "core.default_task_retries" = 1
    "core.parallelism"          = 16
  }

  airflow_version    = "1.10.12"
  dag_s3_path        = aws_s3_object.dags.key
  environment_class  = "mw1.medium"
  execution_role_arn = aws_iam_role.test.arn
  kms_key            = aws_kms_key.test.arn

  logging_configuration {
    dag_processing_logs {
      enabled   = true
      log_level = "INFO"
    }

    scheduler_logs {
      enabled   = true
      log_level = "WARNING"
    }

    task_logs {
      enabled   = true
      log_level = "ERROR"
    }

    webserver_logs {
      enabled   = true
      log_level = "CRITICAL"
    }

    worker_logs {
      enabled   = true
      log_level = "WARNING"
    }
  }

  max_workers = 20
  min_workers = 15
  name        = %[1]q

  network_configuration {
    security_group_ids = [aws_security_group.test.id]
    subnet_ids         = aws_subnet.private[*].id
  }

  plugins_s3_path                 = aws_s3_object.plugins.key
  requirements_s3_path            = aws_s3_object.requirements.key
  schedulers                      = 1
  source_bucket_arn               = aws_s3_bucket.test.arn
  webserver_access_mode           = "PUBLIC_ONLY"
  weekly_maintenance_window_start = "SAT:03:00"

  tags = {
    Name        = %[1]q
    Environment = "production"
  }
}

data "aws_region" "current" {}

resource "aws_kms_key" "test" {
  description = "Key for a Terraform ACC test"
  key_usage   = "ENCRYPT_DECRYPT"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": "kms:*",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "logs.${data.aws_region.current.name}.${data.aws_partition.current.dns_suffix}"
      },
      "Action": "kms:*",
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_s3_object" "plugins" {
  bucket  = aws_s3_bucket.test.id
  acl     = "private"
  key     = "plugins.zip"
  content = ""
}

resource "aws_s3_object" "requirements" {
  bucket  = aws_s3_bucket.test.id
  acl     = "private"
  key     = "requirements.txt"
  content = ""
}

`, rName)
}

func testAccEnvironmentConfig_pluginsS3ObjectVersion(rName, content string) string {
	return testAccEnvironmentBase(rName) + fmt.Sprintf(`
resource "aws_mwaa_environment" "test" {
  dag_s3_path        = aws_s3_object.dags.key
  execution_role_arn = aws_iam_role.test.arn
  name               = %[1]q

  network_configuration {
    security_group_ids = [aws_security_group.test.id]
    subnet_ids         = aws_subnet.private[*].id
  }

  plugins_s3_path           = aws_s3_object.plugins.key
  plugins_s3_object_version = aws_s3_object.plugins.version_id

  source_bucket_arn = aws_s3_bucket.test.arn
}

resource "aws_s3_object" "plugins" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket  = aws_s3_bucket.test.id
  acl     = "private"
  key     = "plugins.zip"
  content = %q
}
`, rName, content)
}
