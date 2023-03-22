package sagemaker_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfsagemaker "github.com/hashicorp/terraform-provider-aws/internal/service/sagemaker"
)

func TestAccSageMakerModel_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "primary_container.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "primary_container.0.image", "data.aws_sagemaker_prebuilt_ecr_image.test", "registry_path"),
					resource.TestCheckResourceAttr(resourceName, "primary_container.0.mode", "SingleModel"),
					resource.TestCheckResourceAttr(resourceName, "primary_container.0.environment.%", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "execution_role_arn", "aws_iam_role.test", "arn"),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "sagemaker", fmt.Sprintf("model/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "enable_network_isolation", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "inference_execution_config.#", "0"),
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

func TestAccSageMakerModel_inferenceExecution(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_inferenceExecution(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "inference_execution_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "inference_execution_config.0.mode", "Serial"),
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

func TestAccSageMakerModel_tags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccModelConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccModelConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
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

func TestAccSageMakerModel_primaryContainerModelDataURL(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_primaryContainerDataURL(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "primary_container.0.model_data_url"),
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

func TestAccSageMakerModel_primaryContainerHostname(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_primaryContainerHostname(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "primary_container.0.container_hostname", "test"),
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

func TestAccSageMakerModel_primaryContainerImage(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_primaryContainerImage(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "primary_container.0.image_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "primary_container.0.image_config.0.repository_access_mode", "Platform"),
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

func TestAccSageMakerModel_primaryContainerEnvironment(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_primaryContainerEnvironment(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "primary_container.0.environment.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "primary_container.0.environment.test", "bar"),
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

func TestAccSageMakerModel_primaryContainerModeSingle(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_primaryContainerModeSingle(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "primary_container.0.mode", "SingleModel"),
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

func TestAccSageMakerModel_containers(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_containers(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "container.0.image", "data.aws_sagemaker_prebuilt_ecr_image.test", "registry_path"),
					resource.TestCheckResourceAttrPair(resourceName, "container.1.image", "data.aws_sagemaker_prebuilt_ecr_image.test", "registry_path"),
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

func TestAccSageMakerModel_vpc(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_vpcBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.subnets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.security_group_ids.#", "2"),
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

func TestAccSageMakerModel_primaryContainerPrivateDockerRegistry(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_primaryContainerPrivateDockerRegistry(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "primary_container.0.image_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "primary_container.0.image_config.0.repository_access_mode", "Vpc"),
					resource.TestCheckResourceAttr(resourceName, "primary_container.0.image_config.0.repository_auth_config.0.repository_credentials_provider_arn", "arn:aws:lambda:us-east-2:123456789012:function:my-function:1"), //lintignore:AWSAT003,AWSAT005
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

func TestAccSageMakerModel_networkIsolation(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_networkIsolation(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_network_isolation", "true"),
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

func TestAccSageMakerModel_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_model.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModelConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfsagemaker.ResourceModel(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckModelDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SageMakerConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_sagemaker_model" {
			continue
		}

		resp, err := conn.ListModels(&sagemaker.ListModelsInput{
			NameContains: aws.String(rs.Primary.ID),
		})
		if err == nil {
			if len(resp.Models) > 0 {
				return fmt.Errorf("SageMaker models still exist")
			}

			return nil
		}

		sagemakerErr, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if sagemakerErr.Code() != "ResourceNotFound" {
			return err
		}
	}

	return nil
}

func testAccCheckModelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No sagmaker model ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SageMakerConn
		DescribeModelOpts := &sagemaker.DescribeModelInput{
			ModelName: aws.String(rs.Primary.ID),
		}
		_, err := conn.DescribeModel(DescribeModelOpts)

		return err
	}
}

func testAccModelConfigBase(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name               = %[1]q
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.test.json
}

data "aws_iam_policy_document" "test" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["sagemaker.amazonaws.com"]
    }
  }
}

data "aws_sagemaker_prebuilt_ecr_image" "test" {
  repository_name = "kmeans"
}
`, rName)
}

func testAccModelConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.test.arn

  primary_container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
  }
}
`, rName))
}

func testAccModelConfig_inferenceExecution(rName string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.test.arn

  inference_execution_config {
    mode = "Serial"
  }

  container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
  }

  container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
  }
}
`, rName))
}

func testAccModelConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.test.arn

  primary_container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccModelConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.test.arn

  primary_container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccModelConfig_primaryContainerDataURL(rName string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.test.arn

  primary_container {
    image          = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
    model_data_url = "https://s3.amazonaws.com/${aws_s3_object.test.bucket}/${aws_s3_object.test.key}"
  }
}

resource "aws_iam_policy" "test" {
  name        = %[1]q
  description = "Allow SageMaker to create model"
  policy      = data.aws_iam_policy_document.policy.json
}

data "aws_iam_policy_document" "policy" {
  statement {
    effect = "Allow"

    actions = [
      "cloudwatch:PutMetricData",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "logs:CreateLogGroup",
      "logs:DescribeLogStreams",
      "ecr:GetAuthorizationToken",
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
    ]

    resources = [
      "*",
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${aws_s3_bucket.test.arn}/*",
    ]
  }
}

resource "aws_iam_role_policy_attachment" "test" {
  role       = aws_iam_role.test.name
  policy_arn = aws_iam_policy.test.arn
}

resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_s3_bucket_acl" "test" {
  bucket = aws_s3_bucket.test.id
  acl    = "private"
}

resource "aws_s3_object" "test" {
  bucket  = aws_s3_bucket.test.bucket
  key     = "model.tar.gz"
  content = "some-data"
}
`, rName))
}

func testAccModelConfig_primaryContainerHostname(rName string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.test.arn

  primary_container {
    image              = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
    container_hostname = "test"
  }
}
`, rName))
}

func testAccModelConfig_primaryContainerImage(rName string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.test.arn

  primary_container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path

    image_config {
      repository_access_mode = "Platform"
    }
  }
}
`, rName))
}

func testAccModelConfig_primaryContainerEnvironment(rName string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.test.arn

  primary_container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path

    environment = {
      test = "bar"
    }
  }
}
`, rName))
}

func testAccModelConfig_primaryContainerModeSingle(rName string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.test.arn

  primary_container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
    mode  = "SingleModel"
  }
}
`, rName))
}

func testAccModelConfig_containers(rName string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.test.arn

  container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
  }

  container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
  }
}
`, rName))
}

func testAccModelConfig_networkIsolation(rName string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name                     = %[1]q
  execution_role_arn       = aws_iam_role.test.arn
  enable_network_isolation = true

  primary_container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
  }
}
`, rName))
}

func testAccModelConfig_vpcBasic(rName string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name                     = %[1]q
  execution_role_arn       = aws_iam_role.test.arn
  enable_network_isolation = true

  primary_container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
  }

  vpc_config {
    subnets            = [aws_subnet.test.id, aws_subnet.bar.id]
    security_group_ids = [aws_security_group.test.id, aws_security_group.bar.id]
  }
}

resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.1.1.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "bar" {
  cidr_block        = "10.1.2.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = "%[1]s-1"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "bar" {
  name   = "%[1]s-2"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

// lintignore:AWSAT003,AWSAT005
func testAccModelConfig_primaryContainerPrivateDockerRegistry(rName string) string {
	return acctest.ConfigCompose(testAccModelConfigBase(rName), acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_sagemaker_model" "test" {
  name                     = %[1]q
  execution_role_arn       = aws_iam_role.test.arn
  enable_network_isolation = true

  primary_container {
    image = "registry.example.com/test-model"

    image_config {
      repository_access_mode = "Vpc"

      repository_auth_config {
        repository_credentials_provider_arn = "arn:aws:lambda:us-east-2:123456789012:function:my-function:1"
      }
    }
  }

  vpc_config {
    subnets            = [aws_subnet.test.id]
    security_group_ids = [aws_security_group.test.id]
  }
}

resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.1.1.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  name   = "%[1]s-1"
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName))
}
