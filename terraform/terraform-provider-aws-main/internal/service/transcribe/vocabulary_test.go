package transcribe_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tftranscribe "github.com/hashicorp/terraform-provider-aws/internal/service/transcribe"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccTranscribeVocabulary_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var vocabulary transcribe.GetVocabularyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transcribe_vocabulary.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.TranscribeEndpointID, t)
			testAccVocabulariesPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.TranscribeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVocabularyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVocabularyConfig_basicFile(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &vocabulary),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "download_uri"),
					resource.TestCheckResourceAttr(resourceName, "language_code", "en-US"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vocabulary_file_uri", "download_uri"},
			},
		},
	})
}

func TestAccTranscribeVocabulary_basicPhrases(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var vocabulary transcribe.GetVocabularyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transcribe_vocabulary.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.TranscribeEndpointID, t)
			testAccVocabulariesPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.TranscribeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVocabularyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVocabularyConfig_basicPhrases(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &vocabulary),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "download_uri"),
					resource.TestCheckResourceAttr(resourceName, "language_code", "en-US"),
				),
			},
		},
	})
}

func TestAccTranscribeVocabulary_updateS3URI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var vocabulary transcribe.GetVocabularyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transcribe_vocabulary.test"

	file1 := "test1.txt"
	file2 := "test2.txt"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.TranscribeEndpointID, t)
			testAccVocabulariesPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.TranscribeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVocabularyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVocabularyConfig_updateFile(rName, file1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &vocabulary),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "vocabulary_file_uri", "s3://"+rName+"/transcribe/test1.txt"),
				),
			},
			{
				Config: testAccVocabularyConfig_updateFile(rName, file2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &vocabulary),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "vocabulary_file_uri", "s3://"+rName+"/transcribe/test2.txt"),
				),
			},
		},
	})
}

func TestAccTranscribeVocabulary_updateTags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var vocabulary transcribe.GetVocabularyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transcribe_vocabulary.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.TranscribeEndpointID, t)
			testAccVocabulariesPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.TranscribeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVocabularyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVocabularyConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &vocabulary),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccVocabularyConfig_tags2(rName, "key1", "value1", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &vocabulary),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccVocabularyConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &vocabulary),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccTranscribeVocabulary_disappears(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var vocabulary transcribe.GetVocabularyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transcribe_vocabulary.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.TranscribeEndpointID, t)
			testAccVocabulariesPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.TranscribeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVocabularyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVocabularyConfig_basicFile(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &vocabulary),
					acctest.CheckResourceDisappears(acctest.Provider, tftranscribe.ResourceVocabulary(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckVocabularyDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).TranscribeConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_transcribe_vocabulary" {
			continue
		}

		_, err := tftranscribe.FindVocabularyByName(context.Background(), conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return create.Error(names.Transcribe, create.ErrActionCheckingDestroyed, tftranscribe.ResNameVocabulary, rs.Primary.ID, errors.New("not destroyed"))
	}

	return nil
}

func testAccCheckVocabularyExists(name string, vocabulary *transcribe.GetVocabularyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.Transcribe, create.ErrActionCheckingExistence, tftranscribe.ResNameVocabulary, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.Transcribe, create.ErrActionCheckingExistence, tftranscribe.ResNameVocabulary, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).TranscribeConn
		resp, err := tftranscribe.FindVocabularyByName(context.Background(), conn, rs.Primary.ID)

		if err != nil {
			return create.Error(names.Transcribe, create.ErrActionCheckingExistence, tftranscribe.ResNameVocabulary, rs.Primary.ID, err)
		}

		*vocabulary = *resp

		return nil
	}
}

func testAccVocabulariesPreCheck(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).TranscribeConn

	input := &transcribe.ListVocabulariesInput{}

	_, err := conn.ListVocabularies(context.Background(), input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccVocabularyBaseConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_s3_object" "object1" {
  bucket = aws_s3_bucket.test.id
  key    = "transcribe/test1.txt"
  source = "test-fixtures/vocabulary_test1.txt"
}

resource "aws_s3_object" "object2" {
  bucket = aws_s3_bucket.test.id
  key    = "transcribe/test2.txt"
  source = "test-fixtures/vocabulary_test2.txt"
}

`, rName)
}

func testAccVocabularyConfig_basicFile(rName string) string {
	return acctest.ConfigCompose(
		testAccVocabularyBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_transcribe_vocabulary" "test" {
  vocabulary_name     = %[1]q
  language_code       = "en-US"
  vocabulary_file_uri = "s3://${aws_s3_bucket.test.id}/${aws_s3_object.object1.key}"

  tags = {
    tag1 = "value1"
    tag2 = "value3"
  }

  depends_on = [
    aws_s3_object.object1
  ]
}
`, rName))
}

func testAccVocabularyConfig_basicPhrases(rName string) string {
	return acctest.ConfigCompose(
		testAccVocabularyBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_transcribe_vocabulary" "test" {
  vocabulary_name = %[1]q
  language_code   = "en-US"
  phrases         = ["Los-Angeles", "CLI", "Eva-Maria"]

  tags = {
    tag1 = "value1"
    tag2 = "value3"
  }

  depends_on = [
    aws_s3_object.object1
  ]
}
`, rName))
}

func testAccVocabularyConfig_updateFile(rName, fileName string) string {
	return acctest.ConfigCompose(
		testAccVocabularyBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_transcribe_vocabulary" "test" {
  vocabulary_name     = %[1]q
  language_code       = "en-US"
  vocabulary_file_uri = "s3://${aws_s3_bucket.test.id}/transcribe/%[2]s"

  tags = {
    tag1 = "value1"
    tag2 = "value3"
  }

  depends_on = [
    aws_s3_object.object1,
    aws_s3_object.object2
  ]
}
`, rName, fileName))
}

func testAccVocabularyConfig_tags1(rName, key1, value1 string) string {
	return acctest.ConfigCompose(
		testAccVocabularyBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_transcribe_vocabulary" "test" {
  vocabulary_name     = %[1]q
  language_code       = "en-US"
  vocabulary_file_uri = "s3://${aws_s3_bucket.test.id}/${aws_s3_object.object1.key}"

  tags = {
    %[2]q = %[3]q
  }

  depends_on = [
    aws_s3_object.object1
  ]
}
`, rName, key1, value1))
}

func testAccVocabularyConfig_tags2(rName, key1, value1, key2, value2 string) string {
	return acctest.ConfigCompose(
		testAccVocabularyBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_transcribe_vocabulary" "test" {
  vocabulary_name     = %[1]q
  language_code       = "en-US"
  vocabulary_file_uri = "s3://${aws_s3_bucket.test.id}/${aws_s3_object.object1.key}"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }

  depends_on = [
    aws_s3_object.object1
  ]
}
`, rName, key1, value1, key2, value2))
}
