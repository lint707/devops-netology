package docdb_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/docdb"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccDocDBClusterInstance_basic(t *testing.T) {
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.cluster_instances"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(resourceName, &v),
					testAccCheckClusterInstanceAttributes(&v),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "rds", regexp.MustCompile(`db:.+`)),
					resource.TestCheckResourceAttr(resourceName, "auto_minor_version_upgrade", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "preferred_maintenance_window"),
					resource.TestCheckResourceAttrSet(resourceName, "preferred_backup_window"),
					resource.TestCheckResourceAttrSet(resourceName, "dbi_resource_id"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "engine_version"),
					resource.TestCheckResourceAttrSet(resourceName, "ca_cert_identifier"),
					resource.TestCheckResourceAttr(resourceName, "engine", "docdb"),
				),
			},
			{
				Config: testAccClusterInstanceConfig_modified(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(resourceName, &v),
					testAccCheckClusterInstanceAttributes(&v),
					resource.TestCheckResourceAttr(resourceName, "auto_minor_version_upgrade", "false"),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
				},
			},
		},
	})
}

func TestAccDocDBClusterInstance_az(t *testing.T) {
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.cluster_instances"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_az(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(resourceName, &v),
					testAccCheckClusterInstanceAttributes(&v),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone"),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
				},
			},
		},
	})
}

func TestAccDocDBClusterInstance_namePrefix(t *testing.T) {
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.test"
	rNamePrefix := acctest.ResourcePrefix
	rName := sdkacctest.RandomWithPrefix(rNamePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_namePrefix(rName, rNamePrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(resourceName, &v),
					testAccCheckClusterInstanceAttributes(&v),
					resource.TestCheckResourceAttr(resourceName, "db_subnet_group_name", rName),
					resource.TestMatchResourceAttr(resourceName, "identifier", regexp.MustCompile(fmt.Sprintf("^%s", rNamePrefix))),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
				},
			},
		},
	})
}

func TestAccDocDBClusterInstance_generatedName(t *testing.T) {
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_generatedName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(resourceName, &v),
					testAccCheckClusterInstanceAttributes(&v),
					resource.TestMatchResourceAttr(resourceName, "identifier", regexp.MustCompile("^tf-")),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
				},
			},
		},
	})
}

func TestAccDocDBClusterInstance_kmsKey(t *testing.T) {
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.cluster_instances"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_kmsKey(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_id", "aws_kms_key.foo", "arn"),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
				},
			},
		},
	})
}

// https://github.com/hashicorp/terraform/issues/5350
func TestAccDocDBClusterInstance_disappears(t *testing.T) {
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.cluster_instances"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(resourceName, &v),
					testAccClusterInstanceDisappears(&v),
				),
				// A non-empty plan is what we want. A crash is what we don't want. :)
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckClusterInstanceAttributes(v *docdb.DBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *v.Engine != "docdb" {
			return fmt.Errorf("bad engine, expected \"docdb\": %#v", *v.Engine)
		}

		if !strings.HasPrefix(*v.DBClusterIdentifier, acctest.ResourcePrefix) {
			return fmt.Errorf("Bad Cluster Identifier prefix:\nexpected: %s-*\ngot: %s", acctest.ResourcePrefix, *v.DBClusterIdentifier)
		}

		return nil
	}
}

func testAccClusterInstanceDisappears(v *docdb.DBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DocDBConn
		opts := &docdb.DeleteDBInstanceInput{
			DBInstanceIdentifier: v.DBInstanceIdentifier,
		}
		if _, err := conn.DeleteDBInstance(opts); err != nil {
			return err
		}
		return resource.Retry(40*time.Minute, func() *resource.RetryError {
			opts := &docdb.DescribeDBInstancesInput{
				DBInstanceIdentifier: v.DBInstanceIdentifier,
			}
			_, err := conn.DescribeDBInstances(opts)
			if err != nil {
				dbinstanceerr, ok := err.(awserr.Error)
				if ok && dbinstanceerr.Code() == "DBInstanceNotFound" {
					return nil
				}
				return resource.NonRetryableError(
					fmt.Errorf("Error retrieving DB Instances: %s", err))
			}
			return resource.RetryableError(fmt.Errorf(
				"Waiting for instance to be deleted: %v", v.DBInstanceIdentifier))
		})
	}
}

func testAccCheckClusterInstanceExists(n string, v *docdb.DBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Instance ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DocDBConn
		resp, err := conn.DescribeDBInstances(&docdb.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		for _, d := range resp.DBInstances {
			if *d.DBInstanceIdentifier == rs.Primary.ID {
				*v = *d
				return nil
			}
		}

		return fmt.Errorf("DB Cluster (%s) not found", rs.Primary.ID)
	}
}

// Add some random to the name, to avoid collision
func testAccClusterInstanceConfig_basic(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_docdb_cluster" "default" {
  cluster_identifier  = %[1]q
  availability_zones  = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]
  master_username     = "foo"
  master_password     = "mustbeeightcharaters"
  skip_final_snapshot = true
}

data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "cluster_instances" {
  identifier         = %[1]q
  cluster_identifier = aws_docdb_cluster.default.id
  instance_class     = data.aws_docdb_orderable_db_instance.test.instance_class
  promotion_tier     = "3"
}
`, rName))
}

func testAccClusterInstanceConfig_modified(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_docdb_cluster" "default" {
  cluster_identifier  = %[1]q
  availability_zones  = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]
  master_username     = "foo"
  master_password     = "mustbeeightcharaters"
  skip_final_snapshot = true
}

data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "cluster_instances" {
  identifier                 = %[1]q
  cluster_identifier         = aws_docdb_cluster.default.id
  instance_class             = data.aws_docdb_orderable_db_instance.test.instance_class
  auto_minor_version_upgrade = false
  promotion_tier             = "3"
}
`, rName))
}

func testAccClusterInstanceConfig_az(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_docdb_cluster" "default" {
  cluster_identifier  = %[1]q
  availability_zones  = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]
  master_username     = "foo"
  master_password     = "mustbeeightcharaters"
  skip_final_snapshot = true
}

data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "cluster_instances" {
  identifier         = %[1]q
  cluster_identifier = aws_docdb_cluster.default.id
  instance_class     = data.aws_docdb_orderable_db_instance.test.instance_class
  promotion_tier     = "3"
  availability_zone  = data.aws_availability_zones.available.names[0]
}
`, rName))
}

func testAccClusterInstanceConfig_namePrefix(rName, rNamePrefix string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "test" {
  identifier_prefix  = %[2]q
  cluster_identifier = aws_docdb_cluster.test.id
  instance_class     = data.aws_docdb_orderable_db_instance.test.instance_class
}

resource "aws_docdb_cluster" "test" {
  cluster_identifier   = %[1]q
  master_username      = "root"
  master_password      = "password"
  db_subnet_group_name = aws_docdb_subnet_group.test.name
  skip_final_snapshot  = true
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "a" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = "%[1]s-a"
  }
}

resource "aws_subnet" "b" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = data.aws_availability_zones.available.names[1]

  tags = {
    Name = "%[1]s-b"
  }
}

resource "aws_docdb_subnet_group" "test" {
  name       = %[1]q
  subnet_ids = [aws_subnet.a.id, aws_subnet.b.id]
}
`, rName, rNamePrefix))
}

func testAccClusterInstanceConfig_generatedName(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "test" {
  cluster_identifier = aws_docdb_cluster.test.id
  instance_class     = data.aws_docdb_orderable_db_instance.test.instance_class
}

resource "aws_docdb_cluster" "test" {
  cluster_identifier   = %[1]q
  master_username      = "root"
  master_password      = "password"
  db_subnet_group_name = aws_docdb_subnet_group.test.name
  skip_final_snapshot  = true
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "a" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = "%[1]s-a"
  }
}

resource "aws_subnet" "b" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = data.aws_availability_zones.available.names[1]

  tags = {
    Name = "%[1]s-b"
  }
}

resource "aws_docdb_subnet_group" "test" {
  name       = %[1]q
  subnet_ids = [aws_subnet.a.id, aws_subnet.b.id]
}
`, rName))
}

func testAccClusterInstanceConfig_kmsKey(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_kms_key" "foo" {
  description = "Terraform acc test %[1]s"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "kms-tf-1",
  "Statement": [
    {
      "Sid": "Enable IAM User Permissions",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": "kms:*",
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_docdb_cluster" "default" {
  cluster_identifier  = %[1]q
  availability_zones  = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]
  master_username     = "foo"
  master_password     = "mustbeeightcharaters"
  storage_encrypted   = true
  kms_key_id          = aws_kms_key.foo.arn
  skip_final_snapshot = true
}

data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "cluster_instances" {
  identifier         = %[1]q
  cluster_identifier = aws_docdb_cluster.default.id
  instance_class     = data.aws_docdb_orderable_db_instance.test.instance_class
}
`, rName))
}
