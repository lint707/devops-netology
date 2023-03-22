package s3_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfs3 "github.com/hashicorp/terraform-provider-aws/internal/service/s3"
)

func TestAccS3BucketInventory_basic(t *testing.T) {
	var conf s3.InventoryConfiguration
	rString := sdkacctest.RandString(8)
	resourceName := "aws_s3_bucket_inventory.test"

	bucketName := fmt.Sprintf("tf-acc-bucket-inventory-%s", rString)
	inventoryName := t.Name()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketInventoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketInventoryConfig_basic(bucketName, inventoryName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketInventoryExistsConfig(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "bucket", bucketName),
					resource.TestCheckNoResourceAttr(resourceName, "filter"),
					resource.TestCheckResourceAttr(resourceName, "name", inventoryName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "included_object_versions", "All"),

					resource.TestCheckResourceAttr(resourceName, "optional_fields.#", "2"),

					resource.TestCheckResourceAttr(resourceName, "schedule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "schedule.0.frequency", "Weekly"),

					resource.TestCheckResourceAttr(resourceName, "destination.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.bucket.#", "1"),
					acctest.CheckResourceAttrGlobalARNNoAccount(resourceName, "destination.0.bucket.0.bucket_arn", "s3", bucketName),
					acctest.CheckResourceAttrAccountID(resourceName, "destination.0.bucket.0.account_id"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.bucket.0.format", "ORC"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.bucket.0.prefix", "inventory"),
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

func TestAccS3BucketInventory_encryptWithSSES3(t *testing.T) {
	var conf s3.InventoryConfiguration
	rString := sdkacctest.RandString(8)
	resourceName := "aws_s3_bucket_inventory.test"

	bucketName := fmt.Sprintf("tf-acc-bucket-inventory-%s", rString)
	inventoryName := t.Name()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketInventoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketInventoryConfig_encryptSSE(bucketName, inventoryName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketInventoryExistsConfig(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "destination.0.bucket.0.encryption.0.sse_s3.#", "1"),
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

func TestAccS3BucketInventory_encryptWithSSEKMS(t *testing.T) {
	var conf s3.InventoryConfiguration
	rString := sdkacctest.RandString(8)
	resourceName := "aws_s3_bucket_inventory.test"

	bucketName := fmt.Sprintf("tf-acc-bucket-inventory-%s", rString)
	inventoryName := t.Name()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketInventoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketInventoryConfig_encryptSSEKMS(bucketName, inventoryName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketInventoryExistsConfig(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "destination.0.bucket.0.encryption.0.sse_kms.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "destination.0.bucket.0.encryption.0.sse_kms.0.key_id", regexp.MustCompile(fmt.Sprintf("^arn:%s:kms:", acctest.Partition()))),
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

func testAccCheckBucketInventoryExistsConfig(n string, res *s3.InventoryConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No S3 bucket inventory configuration ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn
		bucket, name, err := tfs3.BucketInventoryParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		input := &s3.GetBucketInventoryConfigurationInput{
			Bucket: aws.String(bucket),
			Id:     aws.String(name),
		}
		log.Printf("[DEBUG] Reading S3 bucket inventory configuration: %s", input)
		output, err := conn.GetBucketInventoryConfiguration(input)
		if err != nil {
			return err
		}

		*res = *output.InventoryConfiguration

		return nil
	}
}

func testAccCheckBucketInventoryDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_s3_bucket_inventory" {
			continue
		}

		bucket, name, err := tfs3.BucketInventoryParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			input := &s3.GetBucketInventoryConfigurationInput{
				Bucket: aws.String(bucket),
				Id:     aws.String(name),
			}
			log.Printf("[DEBUG] Reading S3 bucket inventory configuration: %s", input)
			output, err := conn.GetBucketInventoryConfiguration(input)
			if err != nil {
				if tfawserr.ErrCodeEquals(err, s3.ErrCodeNoSuchBucket) || tfawserr.ErrMessageContains(err, "NoSuchConfiguration", "The specified configuration does not exist.") {
					return nil
				}
				return resource.NonRetryableError(err)
			}
			if output.InventoryConfiguration != nil {
				return resource.RetryableError(fmt.Errorf("S3 bucket inventory configuration exists: %v", output))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccBucketInventoryBucketConfig(name string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_acl" "test" {
  bucket = aws_s3_bucket.test.id
  acl    = "private"
}
`, name)
}

func testAccBucketInventoryConfig_basic(bucketName, inventoryName string) string {
	return testAccBucketInventoryBucketConfig(bucketName) + fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_s3_bucket_inventory" "test" {
  bucket = aws_s3_bucket.test.id
  name   = %[1]q

  included_object_versions = "All"

  optional_fields = [
    "Size",
    "LastModifiedDate",
  ]

  filter {
    prefix = "documents/"
  }

  schedule {
    frequency = "Weekly"
  }

  destination {
    bucket {
      format     = "ORC"
      bucket_arn = aws_s3_bucket.test.arn
      account_id = data.aws_caller_identity.current.account_id
      prefix     = "inventory"
    }
  }
}
`, inventoryName)
}

func testAccBucketInventoryConfig_encryptSSE(bucketName, inventoryName string) string {
	return testAccBucketInventoryBucketConfig(bucketName) + fmt.Sprintf(`
resource "aws_s3_bucket_inventory" "test" {
  bucket = aws_s3_bucket.test.id
  name   = %[1]q

  included_object_versions = "Current"

  schedule {
    frequency = "Daily"
  }

  destination {
    bucket {
      format     = "CSV"
      bucket_arn = aws_s3_bucket.test.arn

      encryption {
        sse_s3 {}
      }
    }
  }
}
`, inventoryName)
}

func testAccBucketInventoryConfig_encryptSSEKMS(bucketName, inventoryName string) string {
	return testAccBucketInventoryBucketConfig(bucketName) + fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = "Terraform acc test S3 inventory SSE-KMS encryption: %[1]s"
  deletion_window_in_days = 7
}

resource "aws_s3_bucket_inventory" "test" {
  bucket = aws_s3_bucket.test.id
  name   = %[2]q

  included_object_versions = "Current"

  schedule {
    frequency = "Daily"
  }

  destination {
    bucket {
      format     = "Parquet"
      bucket_arn = aws_s3_bucket.test.arn

      encryption {
        sse_kms {
          key_id = aws_kms_key.test.arn
        }
      }
    }
  }
}
`, bucketName, inventoryName)
}
