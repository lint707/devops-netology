package s3control_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfs3control "github.com/hashicorp/terraform-provider-aws/internal/service/s3control"
)

func TestAccS3ControlBucketLifecycleConfiguration_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_bucket_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3control.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleID(rName, "test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "aws_s3control_bucket.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"expiration.#":      "1",
						"expiration.0.days": "365",
						"id":                "test",
						"status":            s3control.ExpirationStatusEnabled,
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

func TestAccS3ControlBucketLifecycleConfiguration_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_bucket_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3control.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleID(rName, "test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfs3control.ResourceBucketLifecycleConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccS3ControlBucketLifecycleConfiguration_RuleAbortIncompleteMultipartUpload_daysAfterInitiation(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_bucket_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3control.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleAbortIncompleteMultipartUploadDaysAfterInitiation(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"abort_incomplete_multipart_upload.#":                       "1",
						"abort_incomplete_multipart_upload.0.days_after_initiation": "1",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleAbortIncompleteMultipartUploadDaysAfterInitiation(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"abort_incomplete_multipart_upload.#":                       "1",
						"abort_incomplete_multipart_upload.0.days_after_initiation": "2",
					}),
				),
			},
		},
	})
}

func TestAccS3ControlBucketLifecycleConfiguration_RuleExpiration_date(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_bucket_lifecycle_configuration.test"
	date1 := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	date2 := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3control.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleExpirationDate(rName, date1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"expiration.#":      "1",
						"expiration.0.date": date1,
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleExpirationDate(rName, date2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"expiration.#":      "1",
						"expiration.0.date": date2,
					}),
				),
			},
		},
	})
}

func TestAccS3ControlBucketLifecycleConfiguration_RuleExpiration_days(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_bucket_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3control.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleExpirationDays(rName, 7),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"expiration.#":      "1",
						"expiration.0.days": "7",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleExpirationDays(rName, 30),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"expiration.#":      "1",
						"expiration.0.days": "30",
					}),
				),
			},
		},
	})
}

func TestAccS3ControlBucketLifecycleConfiguration_RuleExpiration_expiredObjectDeleteMarker(t *testing.T) {
	acctest.Skip(t, "S3 on Outposts does not error or save it in the API when receiving this parameter")
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_bucket_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3control.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleExpirationExpiredObjectDeleteMarker(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"expiration.#": "1",
						"expiration.0.expired_object_delete_marker": "true",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleExpirationExpiredObjectDeleteMarker(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"expiration.#": "1",
						"expiration.0.expired_object_delete_marker": "false",
					}),
				),
			},
		},
	})
}

func TestAccS3ControlBucketLifecycleConfiguration_RuleFilter_prefix(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_bucket_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3control.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleFilterPrefix(rName, "test1/"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"filter.#":        "1",
						"filter.0.prefix": "test1/",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleFilterPrefix(rName, "test2/"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"filter.#":        "1",
						"filter.0.prefix": "test2/",
					}),
				),
			},
		},
	})
}

func TestAccS3ControlBucketLifecycleConfiguration_RuleFilter_tags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_bucket_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3control.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleFilterTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"filter.#":           "1",
						"filter.0.tags.%":    "1",
						"filter.0.tags.key1": "value1",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// There is currently an API model or AWS Go SDK bug where LifecycleFilter.And.Tags
			// does not get populated from the XML response. Reference:
			// https://github.com/aws/aws-sdk-go/issues/3591
			// {
			// 	Config: testAccBucketLifecycleConfigurationConfig_ruleFilterTags2(rName, "key1", "value1updated", "key2", "value2"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckBucketLifecycleConfigurationExists(resourceName),
			// 		resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
			// 		resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
			// 			"filter.#":           "1",
			// 			"filter.0.tags.%":    "2",
			// 			"filter.0.tags.key1": "value1updated",
			// 			"filter.0.tags.key2": "value2",
			// 		}),
			// 	),
			// },
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleFilterTags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"filter.#":           "1",
						"filter.0.tags.%":    "1",
						"filter.0.tags.key2": "value2",
					}),
				),
			},
		},
	})
}

func TestAccS3ControlBucketLifecycleConfiguration_Rule_id(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_bucket_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3control.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleID(rName, "test1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"id": "test1",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleID(rName, "test2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"id": "test2",
					}),
				),
			},
		},
	})
}

func TestAccS3ControlBucketLifecycleConfiguration_Rule_status(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_bucket_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3control.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleStatus(rName, s3control.ExpirationStatusDisabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"status": s3control.ExpirationStatusDisabled,
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketLifecycleConfigurationConfig_ruleStatus(rName, s3control.ExpirationStatusEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLifecycleConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"status": s3control.ExpirationStatusEnabled,
					}),
				),
			},
		},
	})
}

func testAccCheckBucketLifecycleConfigurationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).S3ControlConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_s3control_bucket_lifecycle_configuration" {
			continue
		}

		parsedArn, err := arn.Parse(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error parsing S3 Control Bucket ARN (%s): %w", rs.Primary.ID, err)
		}

		input := &s3control.GetBucketLifecycleConfigurationInput{
			AccountId: aws.String(parsedArn.AccountID),
			Bucket:    aws.String(rs.Primary.ID),
		}

		_, err = conn.GetBucketLifecycleConfiguration(input)

		if tfawserr.ErrCodeEquals(err, "NoSuchBucket") {
			continue
		}

		if tfawserr.ErrCodeEquals(err, "NoSuchLifecycleConfiguration") {
			continue
		}

		if tfawserr.ErrCodeEquals(err, "NoSuchOutpost") {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("S3 Control Bucket Lifecycle Configuration (%s) still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckBucketLifecycleConfigurationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).S3ControlConn

		parsedArn, err := arn.Parse(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error parsing S3 Control Bucket ARN (%s): %w", rs.Primary.ID, err)
		}

		input := &s3control.GetBucketLifecycleConfigurationInput{
			AccountId: aws.String(parsedArn.AccountID),
			Bucket:    aws.String(rs.Primary.ID),
		}

		_, err = conn.GetBucketLifecycleConfiguration(input)

		return err
	}
}

func testAccBucketLifecycleConfigurationConfig_ruleAbortIncompleteMultipartUploadDaysAfterInitiation(rName string, daysAfterInitiation int) string {
	return fmt.Sprintf(`
data "aws_outposts_outposts" "test" {}

data "aws_outposts_outpost" "test" {
  id = tolist(data.aws_outposts_outposts.test.ids)[0]
}

resource "aws_s3control_bucket" "test" {
  bucket     = %[1]q
  outpost_id = data.aws_outposts_outpost.test.id
}

resource "aws_s3control_bucket_lifecycle_configuration" "test" {
  bucket = aws_s3control_bucket.test.arn

  rule {
    abort_incomplete_multipart_upload {
      days_after_initiation = %[2]d
    }

    expiration {
      days = 365
    }

    id = "test"
  }
}
`, rName, daysAfterInitiation)
}

func testAccBucketLifecycleConfigurationConfig_ruleExpirationDate(rName string, date string) string {
	return fmt.Sprintf(`
data "aws_outposts_outposts" "test" {}

data "aws_outposts_outpost" "test" {
  id = tolist(data.aws_outposts_outposts.test.ids)[0]
}

resource "aws_s3control_bucket" "test" {
  bucket     = %[1]q
  outpost_id = data.aws_outposts_outpost.test.id
}

resource "aws_s3control_bucket_lifecycle_configuration" "test" {
  bucket = aws_s3control_bucket.test.arn

  rule {
    expiration {
      date = %[2]q
    }

    id = "test"
  }
}
`, rName, date)
}

func testAccBucketLifecycleConfigurationConfig_ruleExpirationDays(rName string, days int) string {
	return fmt.Sprintf(`
data "aws_outposts_outposts" "test" {}

data "aws_outposts_outpost" "test" {
  id = tolist(data.aws_outposts_outposts.test.ids)[0]
}

resource "aws_s3control_bucket" "test" {
  bucket     = %[1]q
  outpost_id = data.aws_outposts_outpost.test.id
}

resource "aws_s3control_bucket_lifecycle_configuration" "test" {
  bucket = aws_s3control_bucket.test.arn

  rule {
    expiration {
      days = %[2]d
    }

    id = "test"
  }
}
`, rName, days)
}

func testAccBucketLifecycleConfigurationConfig_ruleExpirationExpiredObjectDeleteMarker(rName string, expiredObjectDeleteMarker bool) string {
	return fmt.Sprintf(`
data "aws_outposts_outposts" "test" {}

data "aws_outposts_outpost" "test" {
  id = tolist(data.aws_outposts_outposts.test.ids)[0]
}

resource "aws_s3control_bucket" "test" {
  bucket     = %[1]q
  outpost_id = data.aws_outposts_outpost.test.id
}

resource "aws_s3control_bucket_lifecycle_configuration" "test" {
  bucket = aws_s3control_bucket.test.arn

  rule {
    expiration {
      days                         = %[2]t ? null : 365
      expired_object_delete_marker = %[2]t
    }

    id = "test"
  }
}
`, rName, expiredObjectDeleteMarker)
}

func testAccBucketLifecycleConfigurationConfig_ruleFilterPrefix(rName, prefix string) string {
	return fmt.Sprintf(`
data "aws_outposts_outposts" "test" {}

data "aws_outposts_outpost" "test" {
  id = tolist(data.aws_outposts_outposts.test.ids)[0]
}

resource "aws_s3control_bucket" "test" {
  bucket     = %[1]q
  outpost_id = data.aws_outposts_outpost.test.id
}

resource "aws_s3control_bucket_lifecycle_configuration" "test" {
  bucket = aws_s3control_bucket.test.arn

  rule {
    expiration {
      days = 365
    }

    filter {
      prefix = %[2]q
    }

    id = "test"
  }
}
`, rName, prefix)
}

func testAccBucketLifecycleConfigurationConfig_ruleFilterTags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
data "aws_outposts_outposts" "test" {}

data "aws_outposts_outpost" "test" {
  id = tolist(data.aws_outposts_outposts.test.ids)[0]
}

resource "aws_s3control_bucket" "test" {
  bucket     = %[1]q
  outpost_id = data.aws_outposts_outpost.test.id
}

resource "aws_s3control_bucket_lifecycle_configuration" "test" {
  bucket = aws_s3control_bucket.test.arn

  rule {
    expiration {
      days = 365
    }

    filter {
      tags = {
        %[2]q = %[3]q
      }
    }

    id = "test"
  }
}
`, rName, tagKey1, tagValue1)
}

// See TestAccS3ControlBucketLifecycleConfiguration_RuleFilter_tags note about XML handling bug.
// func testAccBucketLifecycleConfigurationConfig_ruleFilterTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
// 	return fmt.Sprintf(`
// data "aws_outposts_outposts" "test" {}

// data "aws_outposts_outpost" "test" {
//   id = tolist(data.aws_outposts_outposts.test.ids)[0]
// }

// resource "aws_s3control_bucket" "test" {
//   bucket     = %[1]q
//   outpost_id = data.aws_outposts_outpost.test.id
// }

// resource "aws_s3control_bucket_lifecycle_configuration" "test" {
//   bucket = aws_s3control_bucket.test.arn

//   rule {
//     expiration {
//       days = 365
//     }

//     filter {
//       tags = {
//         %[2]q = %[3]q
//         %[4]q = %[5]q
//       }
//     }

//     id = "test"
//   }
// }
// `, rName, tagKey1, tagValue1, tagKey2, tagValue2)
// }

func testAccBucketLifecycleConfigurationConfig_ruleID(rName, id string) string {
	return fmt.Sprintf(`
data "aws_outposts_outposts" "test" {}

data "aws_outposts_outpost" "test" {
  id = tolist(data.aws_outposts_outposts.test.ids)[0]
}

resource "aws_s3control_bucket" "test" {
  bucket     = %[1]q
  outpost_id = data.aws_outposts_outpost.test.id
}

resource "aws_s3control_bucket_lifecycle_configuration" "test" {
  bucket = aws_s3control_bucket.test.arn

  rule {
    expiration {
      days = 365
    }

    id = %[2]q
  }
}
`, rName, id)
}

func testAccBucketLifecycleConfigurationConfig_ruleStatus(rName, status string) string {
	return fmt.Sprintf(`
data "aws_outposts_outposts" "test" {}

data "aws_outposts_outpost" "test" {
  id = tolist(data.aws_outposts_outposts.test.ids)[0]
}

resource "aws_s3control_bucket" "test" {
  bucket     = %[1]q
  outpost_id = data.aws_outposts_outpost.test.id
}

resource "aws_s3control_bucket_lifecycle_configuration" "test" {
  bucket = aws_s3control_bucket.test.arn

  rule {
    expiration {
      days = 365
    }

    id     = "test"
    status = %[2]q
  }
}
`, rName, status)
}
