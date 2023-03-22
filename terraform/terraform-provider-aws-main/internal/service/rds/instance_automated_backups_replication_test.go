package rds_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfrds "github.com/hashicorp/terraform-provider-aws/internal/service/rds"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccRDSInstanceAutomatedBackupsReplication_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_db_instance_automated_backups_replication.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(t),
		CheckDestroy:             testAccCheckInstanceAutomatedBackupsReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceAutomatedBackupsReplicationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceAutomatedBackupsReplicationExist(resourceName),
					resource.TestCheckResourceAttr(resourceName, "retention_period", "7"),
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

func TestAccRDSInstanceAutomatedBackupsReplication_retentionPeriod(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_db_instance_automated_backups_replication.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(t),
		CheckDestroy:             testAccCheckInstanceAutomatedBackupsReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceAutomatedBackupsReplicationConfig_retentionPeriod(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceAutomatedBackupsReplicationExist(resourceName),
					resource.TestCheckResourceAttr(resourceName, "retention_period", "14"),
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

func TestAccRDSInstanceAutomatedBackupsReplication_kmsEncrypted(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_db_instance_automated_backups_replication.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(t),
		CheckDestroy:             testAccCheckInstanceAutomatedBackupsReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceAutomatedBackupsReplicationConfig_kmsEncrypted(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceAutomatedBackupsReplicationExist(resourceName),
					resource.TestCheckResourceAttr(resourceName, "retention_period", "7"),
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
func testAccInstanceAutomatedBackupsReplicationConfig_base(rName string, storageEncrypted bool) string {
	return acctest.ConfigCompose(acctest.ConfigMultipleRegionProvider(2), fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }

  provider = "awsalternate"
}

resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }

  provider = "awsalternate"
}

resource "aws_subnet" "test" {
  count = 2

  cidr_block        = "10.1.${count.index}.0/24"
  availability_zone = data.aws_availability_zones.available.names[count.index]
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  provider = "awsalternate"
}

resource "aws_db_subnet_group" "test" {
  name       = %[1]q
  subnet_ids = aws_subnet.test[*].id

  tags = {
    Name = %[1]q
  }

  provider = "awsalternate"
}

resource "aws_db_instance" "test" {
  allocated_storage       = 10
  identifier              = %[1]q
  engine                  = "postgres"
  engine_version          = "13.4"
  instance_class          = "db.t3.micro"
  name                    = "mydb"
  username                = "masterusername"
  password                = "mustbeeightcharacters"
  backup_retention_period = 7
  skip_final_snapshot     = true
  storage_encrypted       = %[2]t
  db_subnet_group_name    = aws_db_subnet_group.test.name

  provider = "awsalternate"
}
`, rName, storageEncrypted))
}

func testAccInstanceAutomatedBackupsReplicationConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccInstanceAutomatedBackupsReplicationConfig_base(rName, false), `
resource "aws_db_instance_automated_backups_replication" "test" {
  source_db_instance_arn = aws_db_instance.test.arn
}
`)
}

func testAccInstanceAutomatedBackupsReplicationConfig_retentionPeriod(rName string) string {
	return acctest.ConfigCompose(testAccInstanceAutomatedBackupsReplicationConfig_base(rName, false), `
resource "aws_db_instance_automated_backups_replication" "test" {
  source_db_instance_arn = aws_db_instance.test.arn
  retention_period       = 14
}
`)
}

func testAccInstanceAutomatedBackupsReplicationConfig_kmsEncrypted(rName string) string {
	return acctest.ConfigCompose(testAccInstanceAutomatedBackupsReplicationConfig_base(rName, true), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description = %[1]q
}

resource "aws_db_instance_automated_backups_replication" "test" {
  source_db_instance_arn = aws_db_instance.test.arn
  kms_key_id             = aws_kms_key.test.arn
}
`, rName))
}

func testAccCheckInstanceAutomatedBackupsReplicationExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No RDS instance automated backups replication ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).RDSConn

		_, err := tfrds.FindDBInstanceAutomatedBackupByARN(conn, rs.Primary.ID)

		return err
	}
}

func testAccCheckInstanceAutomatedBackupsReplicationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).RDSConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_db_instance_automated_backups_replication" {
			continue
		}

		_, err := tfrds.FindDBInstanceAutomatedBackupByARN(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("RDS instance automated backups replication %s still exists", rs.Primary.ID)
	}

	return nil
}
