package opsworks_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/opsworks"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccOpsWorksCustomLayer_basic(t *testing.T) {
	var v opsworks.Layer
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_custom_layer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomLayerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLayerConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLayerExists(resourceName, &v),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "opsworks", regexp.MustCompile(`layer/.+`)),
					resource.TestCheckResourceAttr(resourceName, "auto_assign_elastic_ips", "false"),
					resource.TestCheckResourceAttr(resourceName, "auto_assign_public_ips", "false"),
					resource.TestCheckResourceAttr(resourceName, "auto_healing", "true"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_configure_recipes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_deploy_recipes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_instance_profile_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "custom_json", ""),
					resource.TestCheckResourceAttr(resourceName, "custom_security_group_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "custom_setup_recipes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_shutdown_recipes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_undeploy_recipes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "drain_elb_on_shutdown", "true"),
					resource.TestCheckResourceAttr(resourceName, "ebs_volume.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ebs_volume.*", map[string]string{
						"type":            "gp2",
						"number_of_disks": "2",
						"mount_point":     "/home",
						"size":            "100",
						"encrypted":       "false",
					}),
					resource.TestCheckResourceAttr(resourceName, "elastic_load_balancer", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_shutdown_timeout", "300"),
					resource.TestCheckResourceAttr(resourceName, "install_updates_on_boot", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "short_name", "tf-ops-acc-custom-layer"),
					resource.TestCheckResourceAttr(resourceName, "system_packages.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "system_packages.*", "git"),
					resource.TestCheckTypeSetElemAttr(resourceName, "system_packages.*", "golang"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "use_ebs_optimized_instances", "false"),
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

// _disappears and _tags for OpsWorks Layers are tested via aws_opsworks_rails_app_layer.

func TestAccOpsWorksCustomLayer_update(t *testing.T) {
	var v opsworks.Layer
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_custom_layer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomLayerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLayerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLayerExists(resourceName, &v),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCustomLayerConfig_update(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "opsworks", regexp.MustCompile(`layer/.+`)),
					resource.TestCheckResourceAttr(resourceName, "auto_assign_elastic_ips", "false"),
					resource.TestCheckResourceAttr(resourceName, "auto_assign_public_ips", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_healing", "true"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_configure_recipes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_deploy_recipes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_instance_profile_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "custom_json", testAccCustomJSON1),
					resource.TestCheckResourceAttr(resourceName, "custom_security_group_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_setup_recipes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_shutdown_recipes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_undeploy_recipes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "drain_elb_on_shutdown", "false"),
					resource.TestCheckResourceAttr(resourceName, "ebs_volume.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ebs_volume.*", map[string]string{
						"type":            "gp2",
						"number_of_disks": "2",
						"mount_point":     "/home",
						"size":            "100",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ebs_volume.*", map[string]string{
						"type":            "io1",
						"number_of_disks": "4",
						"mount_point":     "/var",
						"size":            "100",
						"raid_level":      "1",
						"iops":            "3000",
						"encrypted":       "true",
					}),
					resource.TestCheckResourceAttr(resourceName, "elastic_load_balancer", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_shutdown_timeout", "120"),
					resource.TestCheckResourceAttr(resourceName, "install_updates_on_boot", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "short_name", "tf-ops-acc-custom-layer"),
					resource.TestCheckResourceAttr(resourceName, "system_packages.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "system_packages.*", "git"),
					resource.TestCheckTypeSetElemAttr(resourceName, "system_packages.*", "golang"),
					resource.TestCheckTypeSetElemAttr(resourceName, "system_packages.*", "subversion"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "use_ebs_optimized_instances", "false"),
				),
			},
		},
	})
}

func TestAccOpsWorksCustomLayer_cloudWatch(t *testing.T) {
	var v opsworks.Layer
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_custom_layer.test"
	logGroupResourceName := "aws_cloudwatch_log_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomLayerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLayerConfig_cloudWatch(rName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLayerExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.batch_count", "1000"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.batch_size", "32768"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.buffer_duration", "5000"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.datetime_format", ""),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.encoding", "utf_8"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.file", "/var/log/system.log*"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.file_fingerprint_lines", "1"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.initial_position", "start_of_file"),
					resource.TestCheckResourceAttrPair(resourceName, "cloudwatch_configuration.0.log_streams.0.log_group_name", logGroupResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.multiline_start_pattern", ""),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.time_zone", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCustomLayerConfig_cloudWatch(rName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLayerExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.batch_count", "1000"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.batch_size", "32768"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.buffer_duration", "5000"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.datetime_format", ""),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.encoding", "utf_8"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.file", "/var/log/system.log*"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.file_fingerprint_lines", "1"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.initial_position", "start_of_file"),
					resource.TestCheckResourceAttrPair(resourceName, "cloudwatch_configuration.0.log_streams.0.log_group_name", logGroupResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.multiline_start_pattern", ""),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.time_zone", ""),
				),
			},
			{
				Config: testAccCustomLayerConfig_cloudWatchFull(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLayerExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.batch_count", "2000"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.batch_size", "50000"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.buffer_duration", "6000"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.encoding", "mac_turkish"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.file", "/var/log/system.lo*"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.file_fingerprint_lines", "2"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.initial_position", "end_of_file"),
					resource.TestCheckResourceAttrPair(resourceName, "cloudwatch_configuration.0.log_streams.0.log_group_name", logGroupResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.multiline_start_pattern", "test*"),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_configuration.0.log_streams.0.time_zone", "LOCAL"),
				),
			},
		},
	})
}

func testAccCheckCustomLayerDestroy(s *terraform.State) error {
	return testAccCheckLayerDestroy("aws_opsworks_custom_layer", s)
}

func testAccCustomLayerConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccLayerConfig_base(rName), fmt.Sprintf(`
resource "aws_opsworks_custom_layer" "test" {
  stack_id               = aws_opsworks_stack.test.id
  name                   = %[1]q
  short_name             = "tf-ops-acc-custom-layer"
  auto_assign_public_ips = false

  custom_security_group_ids = aws_security_group.test[*].id

  drain_elb_on_shutdown     = true
  instance_shutdown_timeout = 300

  system_packages = [
    "git",
    "golang",
  ]

  ebs_volume {
    type            = "gp2"
    number_of_disks = 2
    mount_point     = "/home"
    size            = 100
    raid_level      = 0
  }
}
`, rName))
}

func testAccCustomLayerConfig_update(rName string) string {
	return acctest.ConfigCompose(testAccLayerConfig_base(rName), fmt.Sprintf(`
resource "aws_security_group" "extra" {
  name   = "%[1]s-extra"
  vpc_id = aws_vpc.test.id

  ingress {
    from_port   = 8
    to_port     = -1
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_opsworks_custom_layer" "test" {
  stack_id               = aws_opsworks_stack.test.id
  name                   = %[1]q
  short_name             = "tf-ops-acc-custom-layer"
  auto_assign_public_ips = true

  custom_security_group_ids = concat(aws_security_group.test[*].id, [aws_security_group.extra.id])

  drain_elb_on_shutdown     = false
  instance_shutdown_timeout = 120

  system_packages = [
    "git",
    "golang",
    "subversion",
  ]

  ebs_volume {
    type            = "gp2"
    number_of_disks = 2
    mount_point     = "/home"
    size            = 100
    raid_level      = 0
    encrypted       = true
  }

  ebs_volume {
    type            = "io1"
    number_of_disks = 4
    mount_point     = "/var"
    size            = 100
    raid_level      = 1
    iops            = 3000
    encrypted       = true
  }

  custom_json = %[2]q
}
`, rName, testAccCustomJSON1))
}

func testAccCustomLayerConfig_cloudWatch(rName string, enabled bool) string {
	return acctest.ConfigCompose(testAccLayerConfig_base(rName), fmt.Sprintf(`
resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_opsworks_custom_layer" "test" {
  stack_id               = aws_opsworks_stack.test.id
  name                   = %[1]q
  short_name             = "tf-ops-acc-custom-layer"
  auto_assign_public_ips = true

  custom_security_group_ids = aws_security_group.test[*].id

  drain_elb_on_shutdown     = true
  instance_shutdown_timeout = 300

  cloudwatch_configuration {
    enabled = %[2]t

    log_streams {
      log_group_name = aws_cloudwatch_log_group.test.name
      file           = "/var/log/system.log*"
    }
  }
}
`, rName, enabled))
}

func testAccCustomLayerConfig_cloudWatchFull(rName string) string {
	return acctest.ConfigCompose(testAccLayerConfig_base(rName), fmt.Sprintf(`
resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_opsworks_custom_layer" "test" {
  stack_id               = aws_opsworks_stack.test.id
  name                   = %[1]q
  short_name             = "tf-ops-acc-custom-layer"
  auto_assign_public_ips = true

  custom_security_group_ids = aws_security_group.test[*].id

  drain_elb_on_shutdown     = true
  instance_shutdown_timeout = 300

  cloudwatch_configuration {
    enabled = true

    log_streams {
      log_group_name          = aws_cloudwatch_log_group.test.name
      file                    = "/var/log/system.lo*"
      batch_count             = 2000
      batch_size              = 50000
      buffer_duration         = 6000
      encoding                = "mac_turkish"
      file_fingerprint_lines  = "2"
      initial_position        = "end_of_file"
      multiline_start_pattern = "test*"
      time_zone               = "LOCAL"
    }
  }
}
`, rName))
}
