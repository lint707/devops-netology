package s3_test

// WARNING: This code is DEPRECATED and will be removed in a future release!!
// DO NOT apply fixes or enhancements to this file.
// INSTEAD, apply fixes and enhancements to "object_test.go".

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
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
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccS3BucketObject_noNameNoKey(t *testing.T) {
	bucketError := regexp.MustCompile(`bucket must not be empty`)
	keyError := regexp.MustCompile(`key must not be empty`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() {},
				Config:      testAccBucketObjectConfig_basic("", "a key"),
				ExpectError: bucketError,
			},
			{
				PreConfig:   func() {},
				Config:      testAccBucketObjectConfig_basic("a name", ""),
				ExpectError: keyError,
			},
		},
	})
}

func TestAccS3BucketObject_empty(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_empty(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectBody(&obj, ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/test-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_source(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	source := testAccBucketObjectCreateTempFile(t, "{anything will do }")
	defer os.Remove(source)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_source(rName, source),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectBody(&obj, "{anything will do }"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "source", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/test-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_content(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_content(rName, "some_bucket_content"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectBody(&obj, "some_bucket_content"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "content", "content_base64", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/test-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_etagEncryption(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	source := testAccBucketObjectCreateTempFile(t, "{anything will do }")
	defer os.Remove(source)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_etagEncryption(rName, source),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectBody(&obj, "{anything will do }"),
					resource.TestCheckResourceAttr(resourceName, "etag", "7b006ff4d70f68cc65061acf2f802e6f"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "source", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/test-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_contentBase64(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_contentBase64(rName, base64.StdEncoding.EncodeToString([]byte("some_bucket_content"))),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectBody(&obj, "some_bucket_content"),
				),
			},
		},
	})
}

func TestAccS3BucketObject_sourceHashTrigger(t *testing.T) {
	var obj, updated_obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	startingData := "Ebben!"
	changingData := "Ne andrò lontana"

	filename := testAccBucketObjectCreateTempFile(t, startingData)
	defer os.Remove(filename)

	rewriteFile := func(*terraform.State) error {
		if err := os.WriteFile(filename, []byte(changingData), 0644); err != nil {
			os.Remove(filename)
			t.Fatal(err)
		}
		return nil
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_sourceHashTrigger(rName, filename),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectBody(&obj, "Ebben!"),
					resource.TestCheckResourceAttr(resourceName, "source_hash", "7c7e02a79f28968882bb1426c8f8bfc6"),
					rewriteFile,
				),
				ExpectNonEmptyPlan: true,
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_sourceHashTrigger(rName, filename),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &updated_obj),
					testAccCheckBucketObjectBody(&updated_obj, "Ne andrò lontana"),
					resource.TestCheckResourceAttr(resourceName, "source_hash", "cffc5e20de2d21764145b1124c9b337b"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "content", "content_base64", "force_destroy", "source", "source_hash"},
				ImportStateId:           fmt.Sprintf("s3://%s/test-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_withContentCharacteristics(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	source := testAccBucketObjectCreateTempFile(t, "{anything will do }")
	defer os.Remove(source)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_contentCharacteristics(rName, source),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectBody(&obj, "{anything will do }"),
					resource.TestCheckResourceAttr(resourceName, "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr(resourceName, "website_redirect", "http://google.com"),
				),
			},
		},
	})
}

func TestAccS3BucketObject_nonVersioned(t *testing.T) {
	sourceInitial := testAccBucketObjectCreateTempFile(t, "initial object state")
	defer os.Remove(sourceInitial)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	var originalObj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckAssumeRoleARN(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_nonVersioned(rName, sourceInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &originalObj),
					testAccCheckBucketObjectBody(&originalObj, "initial object state"),
					resource.TestCheckResourceAttr(resourceName, "version_id", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "source", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/updateable-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_updates(t *testing.T) {
	var originalObj, modifiedObj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	sourceInitial := testAccBucketObjectCreateTempFile(t, "initial object state")
	defer os.Remove(sourceInitial)
	sourceModified := testAccBucketObjectCreateTempFile(t, "modified object")
	defer os.Remove(sourceInitial)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_updateable(rName, false, sourceInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &originalObj),
					testAccCheckBucketObjectBody(&originalObj, "initial object state"),
					resource.TestCheckResourceAttr(resourceName, "etag", "647d1d58e1011c743ec67d5e8af87b53"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", ""),
				),
			},
			{
				Config: testAccBucketObjectConfig_updateable(rName, false, sourceModified),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &modifiedObj),
					testAccCheckBucketObjectBody(&modifiedObj, "modified object"),
					resource.TestCheckResourceAttr(resourceName, "etag", "1c7fd13df1515c2a13ad9eb068931f09"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "source", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/updateable-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_updateSameFile(t *testing.T) {
	var originalObj, modifiedObj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	startingData := "lane 8"
	changingData := "chicane"

	filename := testAccBucketObjectCreateTempFile(t, startingData)
	defer os.Remove(filename)

	rewriteFile := func(*terraform.State) error {
		if err := os.WriteFile(filename, []byte(changingData), 0644); err != nil {
			os.Remove(filename)
			t.Fatal(err)
		}
		return nil
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_updateable(rName, false, filename),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &originalObj),
					testAccCheckBucketObjectBody(&originalObj, startingData),
					resource.TestCheckResourceAttr(resourceName, "etag", "aa48b42f36a2652cbee40c30a5df7d25"),
					rewriteFile,
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccBucketObjectConfig_updateable(rName, false, filename),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &modifiedObj),
					testAccCheckBucketObjectBody(&modifiedObj, changingData),
					resource.TestCheckResourceAttr(resourceName, "etag", "fafc05f8c4da0266a99154681ab86e8c"),
				),
			},
		},
	})
}

func TestAccS3BucketObject_updatesWithVersioning(t *testing.T) {
	var originalObj, modifiedObj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	sourceInitial := testAccBucketObjectCreateTempFile(t, "initial versioned object state")
	defer os.Remove(sourceInitial)
	sourceModified := testAccBucketObjectCreateTempFile(t, "modified versioned object")
	defer os.Remove(sourceInitial)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_updateable(rName, true, sourceInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &originalObj),
					testAccCheckBucketObjectBody(&originalObj, "initial versioned object state"),
					resource.TestCheckResourceAttr(resourceName, "etag", "cee4407fa91906284e2a5e5e03e86b1b"),
				),
			},
			{
				Config: testAccBucketObjectConfig_updateable(rName, true, sourceModified),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &modifiedObj),
					testAccCheckBucketObjectBody(&modifiedObj, "modified versioned object"),
					resource.TestCheckResourceAttr(resourceName, "etag", "00b8c73b1b50e7cc932362c7225b8e29"),
					testAccCheckBucketObjectVersionIdDiffers(&modifiedObj, &originalObj),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "source", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/updateable-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_updatesWithVersioningViaAccessPoint(t *testing.T) {
	var originalObj, modifiedObj s3.GetObjectOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_object.test"
	accessPointResourceName := "aws_s3_access_point.test"

	sourceInitial := testAccBucketObjectCreateTempFile(t, "initial versioned object state")
	defer os.Remove(sourceInitial)
	sourceModified := testAccBucketObjectCreateTempFile(t, "modified versioned object")
	defer os.Remove(sourceInitial)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_updateableViaAccessPoint(rName, s3.BucketVersioningStatusEnabled, sourceInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &originalObj),
					testAccCheckBucketObjectBody(&originalObj, "initial versioned object state"),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", accessPointResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "etag", "cee4407fa91906284e2a5e5e03e86b1b"),
				),
			},
			{
				Config: testAccBucketObjectConfig_updateableViaAccessPoint(rName, s3.BucketVersioningStatusEnabled, sourceModified),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &modifiedObj),
					testAccCheckBucketObjectBody(&modifiedObj, "modified versioned object"),
					resource.TestCheckResourceAttr(resourceName, "etag", "00b8c73b1b50e7cc932362c7225b8e29"),
					testAccCheckBucketObjectVersionIdDiffers(&modifiedObj, &originalObj),
				),
			},
		},
	})
}

func TestAccS3BucketObject_kms(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	source := testAccBucketObjectCreateTempFile(t, "{anything will do }")
	defer os.Remove(source)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_kmsID(rName, source),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectSSE(resourceName, "aws:kms"),
					testAccCheckBucketObjectBody(&obj, "{anything will do }"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "source", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/test-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_sse(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	source := testAccBucketObjectCreateTempFile(t, "{anything will do }")
	defer os.Remove(source)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_sse(rName, source),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectSSE(resourceName, "AES256"),
					testAccCheckBucketObjectBody(&obj, "{anything will do }"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "source", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/test-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_acl(t *testing.T) {
	var obj1, obj2, obj3 s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_acl(rName, "some_bucket_content", "private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj1),
					testAccCheckBucketObjectBody(&obj1, "some_bucket_content"),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					testAccCheckBucketObjectACL(resourceName, []string{"FULL_CONTROL"}),
				),
			},
			{
				Config: testAccBucketObjectConfig_acl(rName, "some_bucket_content", "public-read"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj2),
					testAccCheckBucketObjectVersionIdEquals(&obj2, &obj1),
					testAccCheckBucketObjectBody(&obj2, "some_bucket_content"),
					resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
					testAccCheckBucketObjectACL(resourceName, []string{"FULL_CONTROL", "READ"}),
				),
			},
			{
				Config: testAccBucketObjectConfig_acl(rName, "changed_some_bucket_content", "private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj3),
					testAccCheckBucketObjectVersionIdDiffers(&obj3, &obj2),
					testAccCheckBucketObjectBody(&obj3, "changed_some_bucket_content"),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					testAccCheckBucketObjectACL(resourceName, []string{"FULL_CONTROL"}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "content", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/test-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_metadata(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_metadata(rName, "key1", "value1", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key2", "value2"),
				),
			},
			{
				Config: testAccBucketObjectConfig_metadata(rName, "key1", "value1updated", "key3", "value3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key3", "value3"),
				),
			},
			{
				Config: testAccBucketObjectConfig_empty(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/test-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_storageClass(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_content(rName, "some_bucket_content"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "STANDARD"),
					testAccCheckBucketObjectStorageClass(resourceName, "STANDARD"),
				),
			},
			{
				Config: testAccBucketObjectConfig_storageClass(rName, "REDUCED_REDUNDANCY"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "REDUCED_REDUNDANCY"),
					testAccCheckBucketObjectStorageClass(resourceName, "REDUCED_REDUNDANCY"),
				),
			},
			{
				Config: testAccBucketObjectConfig_storageClass(rName, "GLACIER"),
				Check: resource.ComposeTestCheckFunc(
					// Can't GetObject on an object in Glacier without restoring it.
					resource.TestCheckResourceAttr(resourceName, "storage_class", "GLACIER"),
					testAccCheckBucketObjectStorageClass(resourceName, "GLACIER"),
				),
			},
			{
				Config: testAccBucketObjectConfig_storageClass(rName, "INTELLIGENT_TIERING"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "INTELLIGENT_TIERING"),
					testAccCheckBucketObjectStorageClass(resourceName, "INTELLIGENT_TIERING"),
				),
			},
			{
				Config: testAccBucketObjectConfig_storageClass(rName, "DEEP_ARCHIVE"),
				Check: resource.ComposeTestCheckFunc(
					// 	Can't GetObject on an object in DEEP_ARCHIVE without restoring it.
					resource.TestCheckResourceAttr(resourceName, "storage_class", "DEEP_ARCHIVE"),
					testAccCheckBucketObjectStorageClass(resourceName, "DEEP_ARCHIVE"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content", "acl", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/test-key", rName),
			},
		},
	})
}

func TestAccS3BucketObject_tags(t *testing.T) {
	var obj1, obj2, obj3, obj4 s3.GetObjectOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_object.object"
	key := "test-key"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_tags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj1),
					testAccCheckBucketObjectBody(&obj1, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "A@AA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_updatedTags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj2),
					testAccCheckBucketObjectVersionIdEquals(&obj2, &obj1),
					testAccCheckBucketObjectBody(&obj2, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "B@BB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "X X"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key4", "DDD"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key5", "E:/"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_noTags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj3),
					testAccCheckBucketObjectVersionIdEquals(&obj3, &obj2),
					testAccCheckBucketObjectBody(&obj3, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_tags(rName, key, "changed stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj4),
					testAccCheckBucketObjectVersionIdDiffers(&obj4, &obj3),
					testAccCheckBucketObjectBody(&obj4, "changed stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "A@AA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content", "acl", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/%s", rName, key),
			},
		},
	})
}

func TestAccS3BucketObject_tagsLeadingSingleSlash(t *testing.T) {
	var obj1, obj2, obj3, obj4 s3.GetObjectOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_object.object"
	key := "/test-key"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_tags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj1),
					testAccCheckBucketObjectBody(&obj1, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "A@AA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_updatedTags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj2),
					testAccCheckBucketObjectVersionIdEquals(&obj2, &obj1),
					testAccCheckBucketObjectBody(&obj2, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "B@BB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "X X"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key4", "DDD"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key5", "E:/"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_noTags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj3),
					testAccCheckBucketObjectVersionIdEquals(&obj3, &obj2),
					testAccCheckBucketObjectBody(&obj3, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_tags(rName, key, "changed stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj4),
					testAccCheckBucketObjectVersionIdDiffers(&obj4, &obj3),
					testAccCheckBucketObjectBody(&obj4, "changed stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "A@AA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content", "acl", "force_destroy"},
				ImportStateId:           fmt.Sprintf("s3://%s/%s", rName, key),
			},
		},
	})
}

func TestAccS3BucketObject_tagsLeadingMultipleSlashes(t *testing.T) {
	var obj1, obj2, obj3, obj4 s3.GetObjectOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_object.object"
	key := "/////test-key"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_tags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj1),
					testAccCheckBucketObjectBody(&obj1, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "A@AA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_updatedTags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj2),
					testAccCheckBucketObjectVersionIdEquals(&obj2, &obj1),
					testAccCheckBucketObjectBody(&obj2, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "B@BB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "X X"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key4", "DDD"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key5", "E:/"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_noTags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj3),
					testAccCheckBucketObjectVersionIdEquals(&obj3, &obj2),
					testAccCheckBucketObjectBody(&obj3, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_tags(rName, key, "changed stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj4),
					testAccCheckBucketObjectVersionIdDiffers(&obj4, &obj3),
					testAccCheckBucketObjectBody(&obj4, "changed stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "A@AA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
				),
			},
		},
	})
}

func TestAccS3BucketObject_tagsMultipleSlashes(t *testing.T) {
	var obj1, obj2, obj3, obj4 s3.GetObjectOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_object.object"
	key := "first//second///third//"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_tags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj1),
					testAccCheckBucketObjectBody(&obj1, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "A@AA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_updatedTags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj2),
					testAccCheckBucketObjectVersionIdEquals(&obj2, &obj1),
					testAccCheckBucketObjectBody(&obj2, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "B@BB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "X X"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key4", "DDD"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key5", "E:/"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_noTags(rName, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj3),
					testAccCheckBucketObjectVersionIdEquals(&obj3, &obj2),
					testAccCheckBucketObjectBody(&obj3, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				PreConfig: func() {},
				Config:    testAccBucketObjectConfig_tags(rName, key, "changed stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj4),
					testAccCheckBucketObjectVersionIdDiffers(&obj4, &obj3),
					testAccCheckBucketObjectBody(&obj4, "changed stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "A@AA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
				),
			},
		},
	})
}

func TestAccS3BucketObject_objectLockLegalHoldStartWithNone(t *testing.T) {
	var obj1, obj2, obj3 s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_noLockLegalHold(rName, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj1),
					testAccCheckBucketObjectBody(&obj1, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", ""),
				),
			},
			{
				Config: testAccBucketObjectConfig_lockLegalHold(rName, "stuff", "ON"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj2),
					testAccCheckBucketObjectVersionIdEquals(&obj2, &obj1),
					testAccCheckBucketObjectBody(&obj2, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", "ON"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", ""),
				),
			},
			// Remove legal hold but create a new object version to test force_destroy
			{
				Config: testAccBucketObjectConfig_lockLegalHold(rName, "changed stuff", "OFF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj3),
					testAccCheckBucketObjectVersionIdDiffers(&obj3, &obj2),
					testAccCheckBucketObjectBody(&obj3, "changed stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", ""),
				),
			},
		},
	})
}

func TestAccS3BucketObject_objectLockLegalHoldStartWithOn(t *testing.T) {
	var obj1, obj2 s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_lockLegalHold(rName, "stuff", "ON"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj1),
					testAccCheckBucketObjectBody(&obj1, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", "ON"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", ""),
				),
			},
			{
				Config: testAccBucketObjectConfig_lockLegalHold(rName, "stuff", "OFF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj2),
					testAccCheckBucketObjectVersionIdEquals(&obj2, &obj1),
					testAccCheckBucketObjectBody(&obj2, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", ""),
				),
			},
		},
	})
}

func TestAccS3BucketObject_objectLockRetentionStartWithNone(t *testing.T) {
	var obj1, obj2, obj3 s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	retainUntilDate := time.Now().UTC().AddDate(0, 0, 10).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_noLockRetention(rName, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj1),
					testAccCheckBucketObjectBody(&obj1, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", ""),
				),
			},
			{
				Config: testAccBucketObjectConfig_lockRetention(rName, "stuff", retainUntilDate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj2),
					testAccCheckBucketObjectVersionIdEquals(&obj2, &obj1),
					testAccCheckBucketObjectBody(&obj2, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", retainUntilDate),
				),
			},
			// Remove retention period but create a new object version to test force_destroy
			{
				Config: testAccBucketObjectConfig_noLockRetention(rName, "changed stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj3),
					testAccCheckBucketObjectVersionIdDiffers(&obj3, &obj2),
					testAccCheckBucketObjectBody(&obj3, "changed stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", ""),
				),
			},
		},
	})
}

func TestAccS3BucketObject_objectLockRetentionStartWithSet(t *testing.T) {
	var obj1, obj2, obj3, obj4 s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	retainUntilDate1 := time.Now().UTC().AddDate(0, 0, 20).Format(time.RFC3339)
	retainUntilDate2 := time.Now().UTC().AddDate(0, 0, 30).Format(time.RFC3339)
	retainUntilDate3 := time.Now().UTC().AddDate(0, 0, 10).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_lockRetention(rName, "stuff", retainUntilDate1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj1),
					testAccCheckBucketObjectBody(&obj1, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", retainUntilDate1),
				),
			},
			{
				Config: testAccBucketObjectConfig_lockRetention(rName, "stuff", retainUntilDate2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj2),
					testAccCheckBucketObjectVersionIdEquals(&obj2, &obj1),
					testAccCheckBucketObjectBody(&obj2, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", retainUntilDate2),
				),
			},
			{
				Config: testAccBucketObjectConfig_lockRetention(rName, "stuff", retainUntilDate3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj3),
					testAccCheckBucketObjectVersionIdEquals(&obj3, &obj2),
					testAccCheckBucketObjectBody(&obj3, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", retainUntilDate3),
				),
			},
			{
				Config: testAccBucketObjectConfig_noLockRetention(rName, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj4),
					testAccCheckBucketObjectVersionIdEquals(&obj4, &obj3),
					testAccCheckBucketObjectBody(&obj4, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "object_lock_retain_until_date", ""),
				),
			},
		},
	})
}

func TestAccS3BucketObject_objectBucketKeyEnabled(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_objectKeyEnabled(rName, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectBody(&obj, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "bucket_key_enabled", "true"),
				),
			},
		},
	})
}

func TestAccS3BucketObject_bucketBucketKeyEnabled(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_bucketKeyEnabled(rName, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectBody(&obj, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "bucket_key_enabled", "true"),
				),
			},
		},
	})
}

func TestAccS3BucketObject_defaultBucketSSE(t *testing.T) {
	var obj1 s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfig_defaultSSE(rName, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj1),
					testAccCheckBucketObjectBody(&obj1, "stuff"),
				),
			},
		},
	})
}

func TestAccS3BucketObject_ignoreTags(t *testing.T) {
	var obj s3.GetObjectOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_object.object"
	key := "test-key"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config: acctest.ConfigCompose(
					acctest.ConfigIgnoreTagsKeyPrefixes1("ignorekey"),
					testAccBucketObjectConfig_noTags(rName, key, "stuff")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectBody(&obj, "stuff"),
					testAccCheckBucketObjectUpdateTags(resourceName, nil, map[string]string{"ignorekey1": "ignorevalue1"}),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					testAccCheckBucketObjectCheckTags(resourceName, map[string]string{
						"ignorekey1": "ignorevalue1",
					}),
				),
			},
			{
				PreConfig: func() {},
				Config: acctest.ConfigCompose(
					acctest.ConfigIgnoreTagsKeyPrefixes1("ignorekey"),
					testAccBucketObjectConfig_tags(rName, key, "stuff")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketObjectExists(resourceName, &obj),
					testAccCheckBucketObjectBody(&obj, "stuff"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "A@AA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
					testAccCheckBucketObjectCheckTags(resourceName, map[string]string{
						"ignorekey1": "ignorevalue1",
						"Key1":       "A@AA",
						"Key2":       "BBB",
						"Key3":       "CCC",
					}),
				),
			},
		},
	})
}

func testAccCheckBucketObjectVersionIdDiffers(first, second *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if first.VersionId == nil {
			return fmt.Errorf("Expected first object to have VersionId: %s", first)
		}
		if second.VersionId == nil {
			return fmt.Errorf("Expected second object to have VersionId: %s", second)
		}

		if *first.VersionId == *second.VersionId {
			return fmt.Errorf("Expected Version IDs to differ, but they are equal (%s)", *first.VersionId)
		}

		return nil
	}
}

func testAccCheckBucketObjectVersionIdEquals(first, second *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if first.VersionId == nil {
			return fmt.Errorf("Expected first object to have VersionId: %s", first)
		}
		if second.VersionId == nil {
			return fmt.Errorf("Expected second object to have VersionId: %s", second)
		}

		if *first.VersionId != *second.VersionId {
			return fmt.Errorf("Expected Version IDs to be equal, but they differ (%s, %s)", *first.VersionId, *second.VersionId)
		}

		return nil
	}
}

func testAccCheckBucketObjectDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_s3_bucket_object" {
			continue
		}

		_, err := conn.HeadObject(
			&s3.HeadObjectInput{
				Bucket:  aws.String(rs.Primary.Attributes["bucket"]),
				Key:     aws.String(rs.Primary.Attributes["key"]),
				IfMatch: aws.String(rs.Primary.Attributes["etag"]),
			})
		if err == nil {
			return fmt.Errorf("AWS S3 Object still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckBucketObjectExists(n string, obj *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No S3 Object ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn

		input := &s3.GetObjectInput{
			Bucket:  aws.String(rs.Primary.Attributes["bucket"]),
			Key:     aws.String(rs.Primary.Attributes["key"]),
			IfMatch: aws.String(rs.Primary.Attributes["etag"]),
		}

		var out *s3.GetObjectOutput

		err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			var err error
			out, err = conn.GetObject(input)

			if tfawserr.ErrCodeEquals(err, s3.ErrCodeNoSuchKey) {
				return resource.RetryableError(
					fmt.Errorf("getting object %s, retrying: %w", rs.Primary.Attributes["bucket"], err),
				)
			}

			if err != nil {
				return resource.NonRetryableError(err)
			}

			return nil
		})
		if tfresource.TimedOut(err) {
			out, err = conn.GetObject(input)
		}

		if err != nil {
			return fmt.Errorf("S3 Object error: %s", err)
		}

		*obj = *out

		return nil
	}
}

func testAccCheckBucketObjectBody(obj *s3.GetObjectOutput, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		body, err := io.ReadAll(obj.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %s", err)
		}
		obj.Body.Close()

		if got := string(body); got != want {
			return fmt.Errorf("wrong result body %q; want %q", got, want)
		}

		return nil
	}
}

func testAccCheckBucketObjectACL(n string, expectedPerms []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn

		out, err := conn.GetObjectAcl(&s3.GetObjectAclInput{
			Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			Key:    aws.String(rs.Primary.Attributes["key"]),
		})

		if err != nil {
			return fmt.Errorf("GetObjectAcl error: %v", err)
		}

		var perms []string
		for _, v := range out.Grants {
			perms = append(perms, *v.Permission)
		}
		sort.Strings(perms)

		if !reflect.DeepEqual(perms, expectedPerms) {
			return fmt.Errorf("Expected ACL permissions to be %v, got %v", expectedPerms, perms)
		}

		return nil
	}
}

func testAccCheckBucketObjectStorageClass(n, expectedClass string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn

		out, err := conn.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			Key:    aws.String(rs.Primary.Attributes["key"]),
		})

		if err != nil {
			return fmt.Errorf("HeadObject error: %v", err)
		}

		// The "STANDARD" (which is also the default) storage
		// class when set would not be included in the results.
		storageClass := s3.StorageClassStandard
		if out.StorageClass != nil {
			storageClass = *out.StorageClass
		}

		if storageClass != expectedClass {
			return fmt.Errorf("Expected Storage Class to be %v, got %v",
				expectedClass, storageClass)
		}

		return nil
	}
}

func testAccCheckBucketObjectSSE(n, expectedSSE string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn

		out, err := conn.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			Key:    aws.String(rs.Primary.Attributes["key"]),
		})

		if err != nil {
			return fmt.Errorf("HeadObject error: %v", err)
		}

		if out.ServerSideEncryption == nil {
			return fmt.Errorf("Expected a non %v Server Side Encryption.", out.ServerSideEncryption)
		}

		sse := *out.ServerSideEncryption
		if sse != expectedSSE {
			return fmt.Errorf("Expected Server Side Encryption %v, got %v.",
				expectedSSE, sse)
		}

		return nil
	}
}

func testAccBucketObjectCreateTempFile(t *testing.T, data string) string {
	tmpFile, err := os.CreateTemp("", "tf-acc-s3-obj")
	if err != nil {
		t.Fatal(err)
	}
	filename := tmpFile.Name()

	err = os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		os.Remove(filename)
		t.Fatal(err)
	}

	return filename
}

func testAccCheckBucketObjectUpdateTags(n string, oldTags, newTags map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn

		return tfs3.ObjectUpdateTags(conn, rs.Primary.Attributes["bucket"], rs.Primary.Attributes["key"], oldTags, newTags)
	}
}

func testAccCheckBucketObjectCheckTags(n string, expectedTags map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn

		got, err := tfs3.ObjectListTags(conn, rs.Primary.Attributes["bucket"], rs.Primary.Attributes["key"])
		if err != nil {
			return err
		}

		want := tftags.New(expectedTags)
		if !reflect.DeepEqual(want, got) {
			return fmt.Errorf("Incorrect tags, want: %v got: %v", want, got)
		}

		return nil
	}
}

func testAccBucketObjectConfig_basic(bucket, key string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket_object" "object" {
  bucket = %[1]q
  key    = %[2]q
}
`, bucket, key)
}

func testAccBucketObjectConfig_empty(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket = aws_s3_bucket.test.bucket
  key    = "test-key"
}
`, rName)
}

func testAccBucketObjectConfig_source(rName string, source string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket       = aws_s3_bucket.test.bucket
  key          = "test-key"
  source       = %[2]q
  content_type = "binary/octet-stream"
}
`, rName, source)
}

func testAccBucketObjectConfig_contentCharacteristics(rName string, source string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket           = aws_s3_bucket.test.bucket
  key              = "test-key"
  source           = %[2]q
  content_language = "en"
  content_type     = "binary/octet-stream"
  website_redirect = "http://google.com"
}
`, rName, source)
}

func testAccBucketObjectConfig_content(rName string, content string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket  = aws_s3_bucket.test.bucket
  key     = "test-key"
  content = %[2]q
}
`, rName, content)
}

func testAccBucketObjectConfig_etagEncryption(rName string, source string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket                 = aws_s3_bucket.test.bucket
  key                    = "test-key"
  server_side_encryption = "AES256"
  source                 = %[2]q
  etag                   = filemd5(%[2]q)
}
`, rName, source)
}

func testAccBucketObjectConfig_contentBase64(rName string, contentBase64 string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket         = aws_s3_bucket.test.bucket
  key            = "test-key"
  content_base64 = %[2]q
}
`, rName, contentBase64)
}

func testAccBucketObjectConfig_sourceHashTrigger(rName string, source string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket      = aws_s3_bucket.test.bucket
  key         = "test-key"
  source      = %[2]q
  source_hash = filemd5(%[2]q)
}
`, rName, source)
}

func testAccBucketObjectConfig_updateable(rName string, bucketVersioning bool, source string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "object_bucket_3" {
  bucket = %[1]q
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.object_bucket_3.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_object" "object" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket = aws_s3_bucket.object_bucket_3.bucket
  key    = "updateable-key"
  source = %[3]q
  etag   = filemd5(%[3]q)
}
`, rName, bucketVersioning, source)
}

func testAccBucketObjectConfig_updateableViaAccessPoint(rName, bucketVersioning, source string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = %[2]q
  }
}

resource "aws_s3_access_point" "test" {
  bucket = aws_s3_bucket.test.bucket
  name   = %[1]q
}

resource "aws_s3_bucket_object" "test" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket = aws_s3_access_point.test.arn
  key    = "updateable-key"
  source = %[3]q
  etag   = filemd5(%[3]q)
}
`, rName, bucketVersioning, source)
}

func testAccBucketObjectConfig_kmsID(rName string, source string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "kms_key_1" {}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket     = aws_s3_bucket.test.bucket
  key        = "test-key"
  source     = %[2]q
  kms_key_id = aws_kms_key.kms_key_1.arn
}
`, rName, source)
}

func testAccBucketObjectConfig_sse(rName string, source string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket                 = aws_s3_bucket.test.bucket
  key                    = "test-key"
  source                 = %[2]q
  server_side_encryption = "AES256"
}
`, rName, source)
}

func testAccBucketObjectConfig_acl(rName string, content, acl string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_object" "object" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket  = aws_s3_bucket.test.bucket
  key     = "test-key"
  content = %[2]q
  acl     = %[3]q
}
`, rName, content, acl)
}

func testAccBucketObjectConfig_storageClass(rName string, storage_class string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket        = aws_s3_bucket.test.bucket
  key           = "test-key"
  content       = "some_bucket_content"
  storage_class = %[2]q
}
`, rName, storage_class)
}

func testAccBucketObjectConfig_tags(rName, key, content string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_object" "object" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket  = aws_s3_bucket.test.bucket
  key     = %[2]q
  content = %[3]q

  tags = {
    Key1 = "A@AA"
    Key2 = "BBB"
    Key3 = "CCC"
  }
}
`, rName, key, content)
}

func testAccBucketObjectConfig_updatedTags(rName, key, content string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_object" "object" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket  = aws_s3_bucket.test.bucket
  key     = %[2]q
  content = %[3]q

  tags = {
    Key2 = "B@BB"
    Key3 = "X X"
    Key4 = "DDD"
    Key5 = "E:/"
  }
}
`, rName, key, content)
}

func testAccBucketObjectConfig_noTags(rName, key, content string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_object" "object" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket  = aws_s3_bucket.test.bucket
  key     = %[2]q
  content = %[3]q
}
`, rName, key, content)
}

func testAccBucketObjectConfig_metadata(rName string, metadataKey1, metadataValue1, metadataKey2, metadataValue2 string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket = aws_s3_bucket.test.bucket
  key    = "test-key"

  metadata = {
    %[2]s = %[3]q
    %[4]s = %[5]q
  }
}
`, rName, metadataKey1, metadataValue1, metadataKey2, metadataValue2)
}

func testAccBucketObjectConfig_noLockLegalHold(rName string, content string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  object_lock_enabled = true
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_object" "object" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket        = aws_s3_bucket.test.bucket
  key           = "test-key"
  content       = %[2]q
  force_destroy = true
}
`, rName, content)
}

func testAccBucketObjectConfig_lockLegalHold(rName string, content, legalHoldStatus string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  object_lock_enabled = true
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_object" "object" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket                        = aws_s3_bucket.test.bucket
  key                           = "test-key"
  content                       = %[2]q
  object_lock_legal_hold_status = %[3]q
  force_destroy                 = true
}
`, rName, content, legalHoldStatus)
}

func testAccBucketObjectConfig_noLockRetention(rName string, content string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  object_lock_enabled = true
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_object" "object" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket        = aws_s3_bucket.test.bucket
  key           = "test-key"
  content       = %[2]q
  force_destroy = true
}
`, rName, content)
}

func testAccBucketObjectConfig_lockRetention(rName string, content, retainUntilDate string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  object_lock_enabled = true
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_object" "object" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.test]

  bucket                        = aws_s3_bucket.test.bucket
  key                           = "test-key"
  content                       = %[2]q
  force_destroy                 = true
  object_lock_mode              = "GOVERNANCE"
  object_lock_retain_until_date = %[3]q
}
`, rName, content, retainUntilDate)
}

func testAccBucketObjectConfig_nonVersioned(rName string, source string) string {
	policy := `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowYeah",
      "Effect": "Allow",
      "Action": "s3:*",
      "Resource": "*"
    },
    {
      "Sid": "DenyStm1",
      "Effect": "Deny",
      "Action": [
        "s3:GetObjectVersion*",
        "s3:ListBucketVersions"
      ],
      "Resource": "*"
    }
  ]
}`

	return acctest.ConfigAssumeRolePolicy(policy) + fmt.Sprintf(`
resource "aws_s3_bucket" "object_bucket_3" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket = aws_s3_bucket.object_bucket_3.bucket
  key    = "updateable-key"
  source = %[2]q
  etag   = filemd5(%[2]q)
}
`, rName, source)
}

func testAccBucketObjectConfig_objectKeyEnabled(rName string, content string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = "Encrypts test objects"
  deletion_window_in_days = 7
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_object" "object" {
  bucket             = aws_s3_bucket.test.bucket
  key                = "test-key"
  content            = %q
  kms_key_id         = aws_kms_key.test.arn
  bucket_key_enabled = true
}
`, rName, content)
}

func testAccBucketObjectConfig_bucketKeyEnabled(rName string, content string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = "Encrypts test objects"
  deletion_window_in_days = 7
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_server_side_encryption_configuration" "test" {
  bucket = aws_s3_bucket.test.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.test.arn
      sse_algorithm     = "aws:kms"
    }
    bucket_key_enabled = true
  }
}

resource "aws_s3_bucket_object" "object" {
  # Must have bucket SSE enabled first
  depends_on = [aws_s3_bucket_server_side_encryption_configuration.test]

  bucket  = aws_s3_bucket.test.bucket
  key     = "test-key"
  content = %q
}
`, rName, content)
}

func testAccBucketObjectConfig_defaultSSE(rName string, content string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = "Encrypts test objects"
  deletion_window_in_days = 7
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_server_side_encryption_configuration" "test" {
  bucket = aws_s3_bucket.test.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.test.arn
      sse_algorithm     = "aws:kms"
    }
  }
}

resource "aws_s3_bucket_object" "object" {
  # Must have bucket SSE enabled first
  depends_on = [aws_s3_bucket_server_side_encryption_configuration.test]

  bucket  = aws_s3_bucket.test.bucket
  key     = "test-key"
  content = %[2]q
}
`, rName, content)
}
