package opsworks_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/opsworks"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfopsworks "github.com/hashicorp/terraform-provider-aws/internal/service/opsworks"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccOpsWorksStack_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_stack.test"
	var v opsworks.Stack

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t)
			testAccPreCheckStacks(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStackConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "agent_version"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "opsworks", regexp.MustCompile(`stack/.+/`)),
					resource.TestCheckResourceAttr(resourceName, "berkshelf_version", "3.2.0"),
					resource.TestCheckResourceAttr(resourceName, "color", ""),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_name", "Chef"),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_version", "11.10"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.#", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "custom_json"),
					resource.TestCheckResourceAttrPair(resourceName, "default_availability_zone", "data.aws_availability_zones.available", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "default_instance_profile_arn"),
					resource.TestCheckResourceAttr(resourceName, "default_os", "Ubuntu 12.04 LTS"),
					resource.TestCheckResourceAttr(resourceName, "default_root_device_type", "instance-store"),
					resource.TestCheckResourceAttr(resourceName, "default_ssh_key_name", ""),
					resource.TestCheckResourceAttrPair(resourceName, "default_subnet_id", "aws_subnet.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "hostname_theme", "Layer_Dependent"),
					resource.TestCheckResourceAttr(resourceName, "manage_berkshelf", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "region", acctest.Region()),
					resource.TestCheckResourceAttrSet(resourceName, "service_role_arn"),
					resource.TestCheckResourceAttr(resourceName, "stack_endpoint", acctest.Region()),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "use_custom_cookbooks", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_opsworks_security_groups", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "aws_vpc.test", "id"),
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

func TestAccOpsWorksStack_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_stack.test"
	var v opsworks.Stack

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t)
			testAccPreCheckStacks(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStackConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					acctest.CheckResourceDisappears(acctest.Provider, tfopsworks.ResourceStack(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccOpsWorksStack_tags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_stack.test"
	var v opsworks.Stack

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t)
			testAccPreCheckStacks(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStackConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStackConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccStackConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccOpsWorksStack_classic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_stack.test"
	var v opsworks.Stack

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckEC2Classic(t)
			acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t)
			testAccPreCheckStacks(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStackClassicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStackConfig_classic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckStackClassicExists(resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "default_availability_zone", "data.aws_availability_zones.available", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "default_subnet_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", ""),
				),
			},
			{
				Config:            testAccStackConfig_classic(rName),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOpsWorksStack_tagsAlternateRegion(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_stack.test"
	var v opsworks.Stack

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t)
			testAccPreCheckStacks(t)
			// This test requires a very particular AWS Region configuration
			// in order to exercise the OpsWorks classic endpoint functionality.
			acctest.PreCheckMultipleRegion(t, 2)
			acctest.PreCheckRegion(t, endpoints.UsEast1RegionID)
			acctest.PreCheckAlternateRegionIs(t, endpoints.UsWest1RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesMultipleRegions(t, 2),
		CheckDestroy:             testAccCheckStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStackConfig_tags1AlternateRegion(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttrWith(resourceName, "arn", func(value string) error {
						if !regexp.MustCompile(arn.ARN{
							Partition: acctest.Partition(),
							Service:   opsworks.ServiceName,
							Region:    acctest.AlternateRegion(),
							AccountID: acctest.AccountID(),
							Resource:  `stack/.+/`,
						}.String()).MatchString(value) {
							return fmt.Errorf("%s doesn't match ARN pattern", value)
						}

						return nil
					}),
					resource.TestCheckResourceAttr(resourceName, "region", acctest.AlternateRegion()),
					// "In this case, the actual API endpoint of the stack is in us-east-1."
					resource.TestCheckResourceAttr(resourceName, "stack_endpoint", endpoints.UsEast1RegionID),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStackConfig_tags2AlternateRegion(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccStackConfig_tags1AlternateRegion(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccOpsWorksStack_allAttributes(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_stack.test"
	var v opsworks.Stack

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t)
			testAccPreCheckStacks(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStackConfig_allAttributes(rName, "4039-20200430042739", "rgb(186, 65, 50)", "main", testAccCustomJSON1, "test1", "Baked_Goods"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "agent_version", "4039-20200430042739"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "opsworks", regexp.MustCompile(`stack/.+/`)),
					resource.TestCheckResourceAttr(resourceName, "berkshelf_version", "3.2.0"),
					resource.TestCheckResourceAttr(resourceName, "color", "rgb(186, 65, 50)"),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_name", "Chef"),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_version", "12"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.password", "avoid-plaintext-passwords"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.revision", "main"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.ssh_key", ""),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.type", "git"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.url", "https://github.com/aws/opsworks-example-cookbooks.git"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.username", "tfacctest"),
					resource.TestCheckResourceAttr(resourceName, "custom_json", testAccCustomJSON1),
					resource.TestCheckResourceAttrPair(resourceName, "default_availability_zone", "data.aws_availability_zones.available", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "default_instance_profile_arn"),
					resource.TestCheckResourceAttr(resourceName, "default_os", "Amazon Linux 2"),
					resource.TestCheckResourceAttr(resourceName, "default_root_device_type", "ebs"),
					resource.TestCheckResourceAttr(resourceName, "default_ssh_key_name", "test1"),
					resource.TestCheckResourceAttrPair(resourceName, "default_subnet_id", "aws_subnet.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "hostname_theme", "Baked_Goods"),
					resource.TestCheckResourceAttr(resourceName, "manage_berkshelf", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "region", acctest.Region()),
					resource.TestCheckResourceAttrSet(resourceName, "service_role_arn"),
					resource.TestCheckResourceAttr(resourceName, "stack_endpoint", acctest.Region()),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "use_custom_cookbooks", "true"),
					resource.TestCheckResourceAttr(resourceName, "use_opsworks_security_groups", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "aws_vpc.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"custom_cookbooks_source.0.password",
				},
			},
			{
				Config: testAccStackConfig_allAttributes(rName, "4038-20200305044341", "rgb(186, 65, 50)", "main", testAccCustomJSON1, "test2", "Scottish_Islands"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "agent_version", "4038-20200305044341"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "opsworks", regexp.MustCompile(`stack/.+/`)),
					resource.TestCheckResourceAttr(resourceName, "berkshelf_version", "3.2.0"),
					resource.TestCheckResourceAttr(resourceName, "color", "rgb(186, 65, 50)"),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_name", "Chef"),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_version", "12"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.password", "avoid-plaintext-passwords"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.revision", "main"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.ssh_key", ""),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.type", "git"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.url", "https://github.com/aws/opsworks-example-cookbooks.git"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.username", "tfacctest"),
					resource.TestCheckResourceAttr(resourceName, "custom_json", testAccCustomJSON1),
					resource.TestCheckResourceAttrPair(resourceName, "default_availability_zone", "data.aws_availability_zones.available", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "default_instance_profile_arn"),
					resource.TestCheckResourceAttr(resourceName, "default_os", "Amazon Linux 2"),
					resource.TestCheckResourceAttr(resourceName, "default_root_device_type", "ebs"),
					resource.TestCheckResourceAttr(resourceName, "default_ssh_key_name", "test2"),
					resource.TestCheckResourceAttrPair(resourceName, "default_subnet_id", "aws_subnet.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "hostname_theme", "Scottish_Islands"),
					resource.TestCheckResourceAttr(resourceName, "manage_berkshelf", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "region", acctest.Region()),
					resource.TestCheckResourceAttrSet(resourceName, "service_role_arn"),
					resource.TestCheckResourceAttr(resourceName, "stack_endpoint", acctest.Region()),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "use_custom_cookbooks", "true"),
					resource.TestCheckResourceAttr(resourceName, "use_opsworks_security_groups", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "aws_vpc.test", "id"),
				),
			},
			{
				Config: testAccStackConfig_allAttributes(rName, "4038-20200305044341", "rgb(209, 105, 41)", "dev", testAccCustomJSON2, "test2", "Scottish_Islands"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "agent_version", "4038-20200305044341"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "opsworks", regexp.MustCompile(`stack/.+/`)),
					resource.TestCheckResourceAttr(resourceName, "berkshelf_version", "3.2.0"),
					resource.TestCheckResourceAttr(resourceName, "color", "rgb(209, 105, 41)"),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_name", "Chef"),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_version", "12"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.password", "avoid-plaintext-passwords"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.revision", "dev"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.ssh_key", ""),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.type", "git"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.url", "https://github.com/aws/opsworks-example-cookbooks.git"),
					resource.TestCheckResourceAttr(resourceName, "custom_cookbooks_source.0.username", "tfacctest"),
					resource.TestCheckResourceAttr(resourceName, "custom_json", testAccCustomJSON2),
					resource.TestCheckResourceAttrPair(resourceName, "default_availability_zone", "data.aws_availability_zones.available", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "default_instance_profile_arn"),
					resource.TestCheckResourceAttr(resourceName, "default_os", "Amazon Linux 2"),
					resource.TestCheckResourceAttr(resourceName, "default_root_device_type", "ebs"),
					resource.TestCheckResourceAttr(resourceName, "default_ssh_key_name", "test2"),
					resource.TestCheckResourceAttrPair(resourceName, "default_subnet_id", "aws_subnet.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "hostname_theme", "Scottish_Islands"),
					resource.TestCheckResourceAttr(resourceName, "manage_berkshelf", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "region", acctest.Region()),
					resource.TestCheckResourceAttrSet(resourceName, "service_role_arn"),
					resource.TestCheckResourceAttr(resourceName, "stack_endpoint", acctest.Region()),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "use_custom_cookbooks", "true"),
					resource.TestCheckResourceAttr(resourceName, "use_opsworks_security_groups", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "aws_vpc.test", "id"),
				),
			},
		},
	})
}

func TestAccOpsWorksStack_windows(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opsworks_stack.test"
	var v opsworks.Stack

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(opsworks.EndpointsID, t)
			testAccPreCheckStacks(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, opsworks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStackConfig_windows(rName, "Microsoft Windows Server 2012 R2 Base"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "agent_version"),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_name", "Chef"),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_version", "12.2"),
					resource.TestCheckResourceAttr(resourceName, "default_os", "Microsoft Windows Server 2012 R2 Base"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStackConfig_windows(rName, "Microsoft Windows Server 2012 R2 with SQL Server Standard"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckStackExists(resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "agent_version"),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_name", "Chef"),
					resource.TestCheckResourceAttr(resourceName, "configuration_manager_version", "12.2"),
					resource.TestCheckResourceAttr(resourceName, "default_os", "Microsoft Windows Server 2012 R2 with SQL Server Standard"),
				),
			},
		},
	})
}

func testAccPreCheckStacks(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).OpsWorksConn

	input := &opsworks.DescribeStacksInput{}

	_, err := conn.DescribeStacks(input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccCheckStackExists(n string, v *opsworks.Stack) resource.TestCheckFunc {
	return testAccCheckStackExistsWithProvider(n, v, func() *schema.Provider { return acctest.Provider })
}

func testAccCheckStackClassicExists(n string, v *opsworks.Stack) resource.TestCheckFunc {
	return testAccCheckStackExistsWithProvider(n, v, func() *schema.Provider { return acctest.ProviderEC2Classic })
}

func testAccCheckStackExistsWithProvider(n string, v *opsworks.Stack, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OpsWorks Stack ID is set")
		}

		conn := providerF().Meta().(*conns.AWSClient).OpsWorksConn

		output, err := tfopsworks.FindStackByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckStackDestroy(s *terraform.State) error {
	return testAccCheckStackDestroyWithProvider(s, func() *schema.Provider { return acctest.Provider })
}

func testAccCheckStackClassicDestroy(s *terraform.State) error {
	return testAccCheckStackDestroyWithProvider(s, func() *schema.Provider { return acctest.ProviderEC2Classic })
}

func testAccCheckStackDestroyWithProvider(s *terraform.State, providerF func() *schema.Provider) error {
	conn := providerF().Meta().(*conns.AWSClient).OpsWorksConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_opsworks_stack" {
			continue
		}

		_, err := tfopsworks.FindStackByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("OpsWorks Stack %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccStackConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "opsworks_service" {
  name = "%[1]s-service"

  assume_role_policy = <<EOT
{
  "Version": "2008-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {
      "Service": "opsworks.amazonaws.com"
    },
    "Action": "sts:AssumeRole"
  }]
}
EOT
}

resource "aws_iam_role_policy" "opsworks_service" {
  name = "%[1]s-service"
  role = aws_iam_role.opsworks_service.id

  policy = <<EOT
{
  "Statement": [{
    "Action": [
      "ec2:*",
      "iam:PassRole",
      "cloudwatch:GetMetricStatistics",
      "elasticloadbalancing:*",
      "rds:*",
      "ecs:*"
    ],
    "Effect": "Allow",
    "Resource": ["*"]
  }]
}
EOT
}

resource "aws_iam_role" "opsworks_instance" {
  name = "%[1]s-instance"

  assume_role_policy = <<EOT
{
  "Version": "2008-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {
      "Service": "ec2.amazonaws.com"
    },
    "Action": "sts:AssumeRole"
  }]
}
EOT
}

resource "aws_iam_instance_profile" "opsworks_instance" {
  name = "%[1]s-instance"
  role = aws_iam_role.opsworks_instance.name
}
`, rName)
}

func testAccStackConfig_baseVPC(rName string) string {
	return acctest.ConfigCompose(testAccStackConfig_base(rName), acctest.ConfigVPCWithSubnets(rName, 2))
}

func testAccStackConfig_baseVPCAlternateRegion(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigMultipleRegionProvider(2),
		testAccStackConfig_base(rName),
		fmt.Sprintf(`
# The VPC (and subnets) must be in the target (alternate) AWS Region.
data "aws_availability_zones" "available" {
  provider = "awsalternate"

  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "test" {
  provider = "awsalternate"

  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  provider = "awsalternate"

  count = 2

  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(aws_vpc.test.cidr_block, 8, count.index)

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccStackConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccStackConfig_baseVPC(rName), fmt.Sprintf(`
resource "aws_opsworks_stack" "test" {
  name                         = %[1]q
  region                       = %[2]q
  service_role_arn             = aws_iam_role.opsworks_service.arn
  default_instance_profile_arn = aws_iam_instance_profile.opsworks_instance.arn
  default_subnet_id            = aws_subnet.test[0].id
  vpc_id                       = aws_vpc.test.id
  use_opsworks_security_groups = false
}
`, rName, acctest.Region()))
}

func testAccStackConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccStackConfig_baseVPC(rName), fmt.Sprintf(`
resource "aws_opsworks_stack" "test" {
  name                         = %[1]q
  region                       = %[2]q
  service_role_arn             = aws_iam_role.opsworks_service.arn
  default_instance_profile_arn = aws_iam_instance_profile.opsworks_instance.arn
  default_subnet_id            = aws_subnet.test[0].id
  vpc_id                       = aws_vpc.test.id
  use_opsworks_security_groups = false

  tags = {
    %[3]q = %[4]q
  }
}
`, rName, acctest.Region(), tagKey1, tagValue1))
}

func testAccStackConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccStackConfig_baseVPC(rName), fmt.Sprintf(`
resource "aws_opsworks_stack" "test" {
  name                         = %[1]q
  region                       = %[2]q
  service_role_arn             = aws_iam_role.opsworks_service.arn
  default_instance_profile_arn = aws_iam_instance_profile.opsworks_instance.arn
  default_subnet_id            = aws_subnet.test[0].id
  vpc_id                       = aws_vpc.test.id
  use_opsworks_security_groups = false

  tags = {
    %[3]q = %[4]q
    %[5]q = %[6]q
  }
}
`, rName, acctest.Region(), tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccStackConfig_classic(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigEC2ClassicRegionProvider(),
		acctest.ConfigAvailableAZsNoOptIn(),
		testAccStackConfig_base(rName),
		fmt.Sprintf(`
resource "aws_opsworks_stack" "test" {
  name                         = %[1]q
  region                       = %[2]q
  service_role_arn             = aws_iam_role.opsworks_service.arn
  default_instance_profile_arn = aws_iam_instance_profile.opsworks_instance.arn
  default_availability_zone    = data.aws_availability_zones.available.names[0]
  use_opsworks_security_groups = false
}
`, rName, acctest.EC2ClassicRegion()))
}

func testAccStackConfig_tags1AlternateRegion(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccStackConfig_baseVPCAlternateRegion(rName), fmt.Sprintf(`
resource "aws_opsworks_stack" "test" {
  name                         = %[1]q
  region                       = %[2]q
  service_role_arn             = aws_iam_role.opsworks_service.arn
  default_instance_profile_arn = aws_iam_instance_profile.opsworks_instance.arn
  default_subnet_id            = aws_subnet.test[0].id
  vpc_id                       = aws_vpc.test.id
  use_opsworks_security_groups = false

  tags = {
    %[3]q = %[4]q
  }
}
`, rName, acctest.AlternateRegion(), tagKey1, tagValue1))
}

func testAccStackConfig_tags2AlternateRegion(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccStackConfig_baseVPCAlternateRegion(rName), fmt.Sprintf(`
resource "aws_opsworks_stack" "test" {
  name                         = %[1]q
  region                       = %[2]q
  service_role_arn             = aws_iam_role.opsworks_service.arn
  default_instance_profile_arn = aws_iam_instance_profile.opsworks_instance.arn
  default_subnet_id            = aws_subnet.test[0].id
  vpc_id                       = aws_vpc.test.id
  use_opsworks_security_groups = false

  tags = {
    %[3]q = %[4]q
    %[5]q = %[6]q
  }
}
`, rName, acctest.AlternateRegion(), tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccStackConfig_allAttributes(rName, agentVersion, color, customCookbookRevision, customJSON, defaultSSHKeyName, hostnameTheme string) string {
	return acctest.ConfigCompose(testAccStackConfig_baseVPC(rName), fmt.Sprintf(`
resource "aws_opsworks_stack" "test" {
  name                         = %[1]q
  region                       = %[2]q
  service_role_arn             = aws_iam_role.opsworks_service.arn
  default_instance_profile_arn = aws_iam_instance_profile.opsworks_instance.arn
  default_subnet_id            = aws_subnet.test[0].id
  vpc_id                       = aws_vpc.test.id
  use_opsworks_security_groups = false

  agent_version                 = %[3]q
  color                         = %[4]q
  configuration_manager_name    = "Chef"
  configuration_manager_version = "12"
  custom_json                   = %[6]q
  default_os                    = "Amazon Linux 2"
  default_root_device_type      = "ebs"
  default_ssh_key_name          = %[7]q
  hostname_theme                = %[8]q
  manage_berkshelf              = false

  use_custom_cookbooks = true
  custom_cookbooks_source {
    type     = "git"
    revision = %[5]q
    url      = "https://github.com/aws/opsworks-example-cookbooks.git"
    password = "avoid-plaintext-passwords"
    username = "tfacctest"
  }
}
`, rName, acctest.Region(), agentVersion, color, customCookbookRevision, customJSON, defaultSSHKeyName, hostnameTheme))
}

func testAccStackConfig_windows(rName, defaultOS string) string {
	return acctest.ConfigCompose(testAccStackConfig_baseVPC(rName), fmt.Sprintf(`
resource "aws_opsworks_stack" "test" {
  name                         = %[1]q
  region                       = %[2]q
  service_role_arn             = aws_iam_role.opsworks_service.arn
  default_instance_profile_arn = aws_iam_instance_profile.opsworks_instance.arn
  default_subnet_id            = aws_subnet.test[0].id
  vpc_id                       = aws_vpc.test.id
  use_opsworks_security_groups = false

  default_os                    = %[3]q
  configuration_manager_version = "12.2"
}
`, rName, acctest.Region(), defaultOS))
}

// Layers
func testAccStackConfig_vpcCreate(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigAvailableAZsNoOptIn(),
		fmt.Sprintf(`
data "aws_partition" "current" {}

data "aws_region" "current" {}

resource "aws_vpc" "test" {
  cidr_block = "10.3.5.0/24"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = aws_vpc.test.cidr_block
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = %[1]q
  }
}

resource "aws_opsworks_stack" "test" {
  name                         = %[1]q
  region                       = data.aws_region.current.name
  vpc_id                       = aws_vpc.test.id
  default_subnet_id            = aws_subnet.test.id
  service_role_arn             = aws_iam_role.opsworks_service.arn
  default_instance_profile_arn = aws_iam_instance_profile.opsworks_instance.arn
  default_os                   = "Amazon Linux 2016.09"
  default_root_device_type     = "ebs"

  custom_json = <<EOF
{
  "key": "value"
}
EOF

  configuration_manager_version = "11.10"
  use_opsworks_security_groups  = false
}

resource "aws_iam_role" "opsworks_service" {
  name = %[1]q

  assume_role_policy = <<EOT
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "opsworks.${data.aws_partition.current.dns_suffix}"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOT
}

resource "aws_iam_role_policy" "opsworks_service" {
  name = %[1]q
  role = aws_iam_role.opsworks_service.id

  policy = <<EOT
{
  "Statement": [
    {
      "Action": [
        "ec2:*",
        "iam:PassRole",
        "cloudwatch:GetMetricStatistics",
        "elasticloadbalancing:*",
        "rds:*",
        "ecs:*"
      ],
      "Effect": "Allow",
      "Resource": [
        "*"
      ]
    }
  ]
}
EOT
}

resource "aws_iam_role" "opsworks_instance" {
  name = "%[1]s-instance"

  assume_role_policy = <<EOT
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.${data.aws_partition.current.dns_suffix}"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOT
}

resource "aws_iam_instance_profile" "opsworks_instance" {
  name = %[1]q
  role = aws_iam_role.opsworks_instance.name
}
`, rName))
}

const (
	testAccCustomJSON1 = `{"key1":"value1"}`
	testAccCustomJSON2 = `{"key2":"value2"}`
)
