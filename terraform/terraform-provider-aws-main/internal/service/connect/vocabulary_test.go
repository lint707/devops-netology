package connect_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/connect"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfconnect "github.com/hashicorp/terraform-provider-aws/internal/service/connect"
)

func testAccVocabulary_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}
	var v connect.DescribeVocabularyOutput
	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")

	content := "Phrase\tIPA\tSoundsLike\tDisplayAs\nLos-Angeles\t\t\tLos Angeles\nF.B.I.\tɛ f b i aɪ\t\tFBI\nEtienne\t\teh-tee-en\t"
	languageCode := "en-US"

	resourceName := "aws_connect_vocabulary.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, connect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVocabularyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVocabularyConfig_basic(rName, rName2, content, languageCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "content", content),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "aws_connect_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "language_code", languageCode),
					resource.TestCheckResourceAttrSet(resourceName, "last_modified_time"),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "Value1"),
					resource.TestCheckResourceAttrSet(resourceName, "vocabulary_id"),
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

func testAccVocabulary_disappears(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}
	var v connect.DescribeVocabularyOutput
	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")

	content := "Phrase\tIPA\tSoundsLike\tDisplayAs\nLos-Angeles\t\t\tLos Angeles\nF.B.I.\tɛ f b i aɪ\t\tFBI\nEtienne\t\teh-tee-en\t"
	languageCode := "en-US"

	resourceName := "aws_connect_vocabulary.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, connect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVocabularyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVocabularyConfig_basic(rName, rName2, content, languageCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &v),
					acctest.CheckResourceDisappears(acctest.Provider, tfconnect.ResourceVocabulary(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccVocabulary_updateTags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}
	var v connect.DescribeVocabularyOutput
	rName := sdkacctest.RandomWithPrefix("resource-test-terraform")
	rName2 := sdkacctest.RandomWithPrefix("resource-test-terraform")

	content := "Phrase\tIPA\tSoundsLike\tDisplayAs\nLos-Angeles\t\t\tLos Angeles\nF.B.I.\tɛ f b i aɪ\t\tFBI\nEtienne\t\teh-tee-en\t"
	languageCode := "en-US"

	resourceName := "aws_connect_vocabulary.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, connect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVocabularyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVocabularyConfig_basic(rName, rName2, content, languageCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "Value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVocabularyConfig_tags(rName, rName2, content, languageCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "Value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "Value2a"),
				),
			},
			{
				Config: testAccVocabularyConfig_tagsUpdate(rName, rName2, content, languageCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVocabularyExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "Value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "Value2b"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "Value3"),
				),
			},
		},
	})
}

func testAccCheckVocabularyExists(resourceName string, function *connect.DescribeVocabularyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Connect Vocabulary not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Connect Vocabulary ID not set")
		}
		instanceID, vocabularyID, err := tfconnect.VocabularyParseID(rs.Primary.ID)

		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ConnectConn

		params := &connect.DescribeVocabularyInput{
			InstanceId:   aws.String(instanceID),
			VocabularyId: aws.String(vocabularyID),
		}

		getFunction, err := conn.DescribeVocabulary(params)
		if err != nil {
			return err
		}

		*function = *getFunction

		return nil
	}
}

func testAccCheckVocabularyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_connect_vocabulary" {
			continue
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ConnectConn

		instanceID, vocabularyID, err := tfconnect.VocabularyParseID(rs.Primary.ID)

		if err != nil {
			return err
		}

		params := &connect.DescribeVocabularyInput{
			InstanceId:   aws.String(instanceID),
			VocabularyId: aws.String(vocabularyID),
		}

		resp, err := conn.DescribeVocabulary(params)

		if tfawserr.ErrCodeEquals(err, connect.ErrCodeResourceNotFoundException) {
			continue
		}

		if err != nil {
			return err
		}

		// API returns an empty list for Vocabulary if there are none
		if resp.Vocabulary == nil {
			continue
		}
	}

	return nil
}

func testAccVocabularyConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_connect_instance" "test" {
  identity_management_type = "CONNECT_MANAGED"
  inbound_calls_enabled    = true
  instance_alias           = %[1]q
  outbound_calls_enabled   = true
}
`, rName)
}

func testAccVocabularyConfig_basic(rName, rName2, content, languageCode string) string {
	return acctest.ConfigCompose(
		testAccVocabularyConfig_base(rName),
		fmt.Sprintf(`
resource "aws_connect_vocabulary" "test" {
  instance_id   = aws_connect_instance.test.id
  name          = %[1]q
  content       = %[2]q
  language_code = %[3]q

  tags = {
    "Key1" = "Value1"
  }
}
`, rName2, content, languageCode))
}

func testAccVocabularyConfig_tags(rName, rName2, content, languageCode string) string {
	return acctest.ConfigCompose(
		testAccVocabularyConfig_base(rName),
		fmt.Sprintf(`
resource "aws_connect_vocabulary" "test" {
  instance_id   = aws_connect_instance.test.id
  name          = %[1]q
  content       = %[2]q
  language_code = %[3]q

  tags = {
    "Key1" = "Value1"
    "Key2" = "Value2a"
  }
}
`, rName2, content, languageCode))
}

func testAccVocabularyConfig_tagsUpdate(rName, rName2, content, languageCode string) string {
	return acctest.ConfigCompose(
		testAccVocabularyConfig_base(rName),
		fmt.Sprintf(`
resource "aws_connect_vocabulary" "test" {
  instance_id   = aws_connect_instance.test.id
  name          = %[1]q
  content       = %[2]q
  language_code = %[3]q

  tags = {
    "Key1" = "Value1"
    "Key2" = "Value2b"
    "Key3" = "Value3"
  }
}
`, rName2, content, languageCode))
}
