package redshift_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/redshift"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccRedshiftClusterDataSource_basic(t *testing.T) {
	dataSourceName := "data.aws_redshift_cluster.test"
	resourceName := "aws_redshift_cluster.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_nodes.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_nodes.0.public_ip_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "allow_version_upgrade"),
					resource.TestCheckResourceAttrSet(dataSourceName, "automated_snapshot_retention_period"),
					resource.TestCheckResourceAttrSet(dataSourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_identifier"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_parameter_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_public_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_revision_number"),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_type", "single-node"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "database_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "encrypted"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "master_username"),
					resource.TestCheckResourceAttrSet(dataSourceName, "node_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "number_of_nodes"),
					resource.TestCheckResourceAttrSet(dataSourceName, "port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "preferred_maintenance_window"),
					resource.TestCheckResourceAttrSet(dataSourceName, "manual_snapshot_retention_period"),
					resource.TestCheckResourceAttrSet(dataSourceName, "maintenance_track_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "arn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "publicly_accessible"),
					resource.TestCheckResourceAttrPair(dataSourceName, "availability_zone_relocation_enabled", resourceName, "availability_zone_relocation_enabled"),
				),
			},
		},
	})
}

func TestAccRedshiftClusterDataSource_vpc(t *testing.T) {
	dataSourceName := "data.aws_redshift_cluster.test"
	subnetGroupResourceName := "aws_redshift_subnet_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig_vpc(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_id"),
					resource.TestCheckResourceAttr(dataSourceName, "vpc_security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_type", "multi-node"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cluster_subnet_group_name", subnetGroupResourceName, "name"),
				),
			},
		},
	})
}

func TestAccRedshiftClusterDataSource_logging(t *testing.T) {
	dataSourceName := "data.aws_redshift_cluster.test"
	bucketResourceName := "aws_s3_bucket.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig_logging(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "enable_logging", "true"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bucket_name", bucketResourceName, "bucket"),
					resource.TestCheckResourceAttr(dataSourceName, "s3_key_prefix", "cluster-logging/"),
				),
			},
		},
	})
}

func TestAccRedshiftClusterDataSource_availabilityZoneRelocationEnabled(t *testing.T) {
	dataSourceName := "data.aws_redshift_cluster.test"
	resourceName := "aws_redshift_cluster.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig_availabilityZoneRelocationEnabled(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "availability_zone_relocation_enabled", resourceName, "availability_zone_relocation_enabled"),
				),
			},
		},
	})
}

func testAccClusterDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_redshift_cluster" "test" {
  cluster_identifier = %[1]q

  database_name       = "testdb"
  master_username     = "foo"
  master_password     = "Password1"
  node_type           = "dc2.large"
  cluster_type        = "single-node"
  skip_final_snapshot = true
}

data "aws_redshift_cluster" "test" {
  cluster_identifier = aws_redshift_cluster.test.cluster_identifier
}
`, rName)
}

func testAccClusterDataSourceConfig_vpc(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 2), fmt.Sprintf(`
resource "aws_redshift_subnet_group" "test" {
  name       = %[1]q
  subnet_ids = aws_subnet.test[*].id
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_redshift_cluster" "test" {
  cluster_identifier = %[1]q

  database_name             = "testdb"
  master_username           = "foo"
  master_password           = "Password1"
  node_type                 = "dc2.large"
  cluster_type              = "multi-node"
  number_of_nodes           = 2
  publicly_accessible       = false
  cluster_subnet_group_name = aws_redshift_subnet_group.test.name
  vpc_security_group_ids    = [aws_security_group.test.id]
  skip_final_snapshot       = true
}

data "aws_redshift_cluster" "test" {
  cluster_identifier = aws_redshift_cluster.test.cluster_identifier
}
`, rName))
}

func testAccClusterDataSourceConfig_logging(rName string) string {
	return fmt.Sprintf(`
data "aws_redshift_service_account" "test" {}

resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

data "aws_iam_policy_document" "test" {
  statement {
    actions   = ["s3:PutObject"]
    resources = ["${aws_s3_bucket.test.arn}/*"]

    principals {
      identifiers = [data.aws_redshift_service_account.test.arn]
      type        = "AWS"
    }
  }

  statement {
    actions   = ["s3:GetBucketAcl"]
    resources = [aws_s3_bucket.test.arn]

    principals {
      identifiers = [data.aws_redshift_service_account.test.arn]
      type        = "AWS"
    }
  }
}

resource "aws_s3_bucket_policy" "test" {
  bucket = aws_s3_bucket.test.bucket
  policy = data.aws_iam_policy_document.test.json
}

resource "aws_redshift_cluster" "test" {
  depends_on = [aws_s3_bucket_policy.test]

  cluster_identifier  = %[1]q
  cluster_type        = "single-node"
  database_name       = "testdb"
  master_password     = "Password1"
  master_username     = "foo"
  node_type           = "dc2.large"
  skip_final_snapshot = true

  logging {
    bucket_name   = aws_s3_bucket.test.id
    enable        = true
    s3_key_prefix = "cluster-logging/"
  }
}

data "aws_redshift_cluster" "test" {
  cluster_identifier = aws_redshift_cluster.test.cluster_identifier
}
`, rName)
}

func testAccClusterDataSourceConfig_availabilityZoneRelocationEnabled(rName string) string {
	return fmt.Sprintf(`
resource "aws_redshift_cluster" "test" {
  cluster_identifier = %[1]q

  database_name       = "testdb"
  master_username     = "foo"
  master_password     = "Password1"
  node_type           = "ra3.xlplus"
  cluster_type        = "single-node"
  skip_final_snapshot = true
  publicly_accessible = false

  availability_zone_relocation_enabled = true
}

data "aws_redshift_cluster" "test" {
  cluster_identifier = aws_redshift_cluster.test.cluster_identifier
}
`, rName)
}
