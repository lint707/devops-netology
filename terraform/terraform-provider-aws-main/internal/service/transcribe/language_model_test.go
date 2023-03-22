package transcribe_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftranscribe "github.com/hashicorp/terraform-provider-aws/internal/service/transcribe"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccTranscribeLanguageModel_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var languageModel types.LanguageModel
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transcribe_language_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.TranscribeEndpointID, t)
			testAccLanguageModelsPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.TranscribeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLanguageModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageModelConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanguageModelExists(resourceName, &languageModel),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "base_model_name", "NarrowBand"),
					resource.TestCheckResourceAttr(resourceName, "language_code", "en-US"),
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

func TestAccTranscribeLanguageModel_updateTags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var languageModel types.LanguageModel
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transcribe_language_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.TranscribeEndpointID, t)
			testAccLanguageModelsPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.TranscribeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLanguageModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageModelConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanguageModelExists(resourceName, &languageModel),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccLanguageModelConfig_tags2(rName, "key1", "value1", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanguageModelExists(resourceName, &languageModel),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccLanguageModelConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanguageModelExists(resourceName, &languageModel),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccTranscribeLanguageModel_disappears(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var languageModel types.LanguageModel
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transcribe_language_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.TranscribeEndpointID, t)
			testAccLanguageModelsPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.TranscribeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLanguageModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageModelConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanguageModelExists(resourceName, &languageModel),
					acctest.CheckResourceDisappears(acctest.Provider, tftranscribe.ResourceLanguageModel(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckLanguageModelDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).TranscribeConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_transcribe_language_model" {
			continue
		}

		_, err := tftranscribe.FindLanguageModelByName(context.Background(), conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckLanguageModelExists(name string, languageModel *types.LanguageModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Transcribe LanguageModel is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).TranscribeConn
		resp, err := tftranscribe.FindLanguageModelByName(context.Background(), conn, rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error describing Transcribe LanguageModel: %s", err.Error())
		}

		*languageModel = *resp

		return nil
	}
}

func testAccLanguageModelsPreCheck(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).TranscribeConn

	input := &transcribe.ListLanguageModelsInput{}

	_, err := conn.ListLanguageModels(context.Background(), input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccLanguageModelBaseConfig(rName string) string {
	return fmt.Sprintf(`
data "aws_iam_policy_document" "test" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["transcribe.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "test" {
  name               = %[1]q
  assume_role_policy = data.aws_iam_policy_document.test.json
}

resource "aws_iam_role_policy" "test_policy" {
  name = %[1]q
  role = aws_iam_role.test.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "s3:GetObject",
          "s3:ListBucket",
        ]
        Effect   = "Allow"
        Resource = ["*"]
      },
    ]
  })
}

resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_s3_object" "object" {
  bucket = aws_s3_bucket.test.id
  key    = "transcribe/test1.txt"
  source = "test-fixtures/language_model_test1.txt"
}
`, rName)
}

func testAccLanguageModelConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccLanguageModelBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_transcribe_language_model" "test" {
  model_name      = %[1]q
  base_model_name = "NarrowBand"

  input_data_config {
    data_access_role_arn = aws_iam_role.test.arn
    s3_uri               = "s3://${aws_s3_bucket.test.id}/transcribe/"
  }

  language_code = "en-US"

  tags = {
    tag1 = "value1"
  }
}
`, rName))
}

func testAccLanguageModelConfig_tags1(rName, key1, value1 string) string {
	return acctest.ConfigCompose(
		testAccLanguageModelBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_transcribe_language_model" "test" {
  model_name      = %[1]q
  base_model_name = "NarrowBand"

  input_data_config {
    data_access_role_arn = aws_iam_role.test.arn
    s3_uri               = "s3://${aws_s3_bucket.test.id}/transcribe/"
  }

  language_code = "en-US"

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, key1, value1))
}

func testAccLanguageModelConfig_tags2(rName, key1, value1, key2, value2 string) string {
	return acctest.ConfigCompose(
		testAccLanguageModelBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_transcribe_language_model" "test" {
  model_name      = %[1]q
  base_model_name = "NarrowBand"

  input_data_config {
    data_access_role_arn = aws_iam_role.test.arn
    s3_uri               = "s3://${aws_s3_bucket.test.id}/transcribe/"
  }

  language_code = "en-US"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, key1, value1, key2, value2))
}
