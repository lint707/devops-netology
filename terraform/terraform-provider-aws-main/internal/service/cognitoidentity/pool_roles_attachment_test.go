package cognitoidentity_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfcognitoidentity "github.com/hashicorp/terraform-provider-aws/internal/service/cognitoidentity"
)

func TestAccCognitoIdentityPoolRolesAttachment_basic(t *testing.T) {
	resourceName := "aws_cognito_identity_pool_roles_attachment.test"
	name := sdkacctest.RandString(10)
	updatedName := sdkacctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPoolRolesAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolRolesAttachmentConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolRolesAttachmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "identity_pool_id"),
					resource.TestCheckResourceAttrSet(resourceName, "roles.authenticated"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPoolRolesAttachmentConfig_basic(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolRolesAttachmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "identity_pool_id"),
					resource.TestCheckResourceAttrSet(resourceName, "roles.authenticated"),
				),
			},
		},
	})
}

func TestAccCognitoIdentityPoolRolesAttachment_roleMappings(t *testing.T) {
	resourceName := "aws_cognito_identity_pool_roles_attachment.test"
	name := sdkacctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPoolRolesAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolRolesAttachmentConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolRolesAttachmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "identity_pool_id"),
					resource.TestCheckResourceAttr(resourceName, "role_mapping.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "roles.authenticated"),
				),
			},
			{
				Config: testAccPoolRolesAttachmentConfig_roleMappings(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolRolesAttachmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "identity_pool_id"),
					resource.TestCheckResourceAttr(resourceName, "role_mapping.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "roles.authenticated"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPoolRolesAttachmentConfig_roleMappingsUpdated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolRolesAttachmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "identity_pool_id"),
					resource.TestCheckResourceAttr(resourceName, "role_mapping.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "roles.authenticated"),
				),
			},
			{
				Config: testAccPoolRolesAttachmentConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolRolesAttachmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "identity_pool_id"),
					resource.TestCheckResourceAttr(resourceName, "role_mapping.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "roles.authenticated"),
				),
			},
		},
	})
}

func TestAccCognitoIdentityPoolRolesAttachment_disappears(t *testing.T) {
	resourceName := "aws_cognito_identity_pool_roles_attachment.test"
	name := sdkacctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPoolRolesAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolRolesAttachmentConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolRolesAttachmentExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfcognitoidentity.ResourcePoolRolesAttachment(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCognitoIdentityPoolRolesAttachment_roleMappingsWithAmbiguousRoleResolutionError(t *testing.T) {
	name := sdkacctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPoolRolesAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccPoolRolesAttachmentConfig_roleMappingsWithAmbiguousRoleResolutionError(name),
				ExpectError: regexp.MustCompile(`Error validating ambiguous role resolution`),
			},
		},
	})
}

func TestAccCognitoIdentityPoolRolesAttachment_roleMappingsWithRulesTypeError(t *testing.T) {
	name := sdkacctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPoolRolesAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccPoolRolesAttachmentConfig_roleMappingsWithRulesTypeError(name),
				ExpectError: regexp.MustCompile(`mapping_rule is required for Rules`),
			},
		},
	})
}

func TestAccCognitoIdentityPoolRolesAttachment_roleMappingsWithTokenTypeError(t *testing.T) {
	name := sdkacctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPoolRolesAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccPoolRolesAttachmentConfig_roleMappingsWithTokenTypeError(name),
				ExpectError: regexp.MustCompile(`mapping_rule must not be set for Token based role mapping`),
			},
		},
	})
}

func testAccCheckPoolRolesAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Cognito Identity Pool Roles Attachment ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).CognitoIdentityConn

		_, err := conn.GetIdentityPoolRoles(&cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(rs.Primary.Attributes["identity_pool_id"]),
		})

		return err
	}
}

func testAccCheckPoolRolesAttachmentDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).CognitoIdentityConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cognito_identity_pool_roles_attachment" {
			continue
		}

		_, err := conn.GetIdentityPoolRoles(&cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(rs.Primary.Attributes["identity_pool_id"]),
		})

		if err != nil {
			if tfawserr.ErrCodeEquals(err, cognitoidentity.ErrCodeResourceNotFoundException) {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccPoolRolesAttachmentConfig(name string) string {
	return fmt.Sprintf(`
resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %[1]s"
  allow_unauthenticated_identities = false

  supported_login_providers = {
    "graph.facebook.com" = "7346241598935555"
  }
}

# Unauthenticated Role
resource "aws_iam_role" "unauthenticated" {
  name = "cognito_unauthenticated_%[1]s"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "cognito-identity.amazonaws.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "cognito-identity.amazonaws.com:aud": "${aws_cognito_identity_pool.main.id}"
        },
        "ForAnyValue:StringLike": {
          "cognito-identity.amazonaws.com:amr": "unauthenticated"
        }
      }
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "unauthenticated" {
  name = "unauthenticated_policy_%[1]s"
  role = aws_iam_role.unauthenticated.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "mobileanalytics:PutEvents",
        "cognito-sync:*"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
EOF
}

# Authenticated Role
resource "aws_iam_role" "authenticated" {
  name = "cognito_authenticated_%[1]s"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "cognito-identity.amazonaws.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "cognito-identity.amazonaws.com:aud": "${aws_cognito_identity_pool.main.id}"
        },
        "ForAnyValue:StringLike": {
          "cognito-identity.amazonaws.com:amr": "authenticated"
        }
      }
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "authenticated" {
  name = "authenticated_policy_%[1]s"
  role = aws_iam_role.authenticated.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "mobileanalytics:PutEvents",
        "cognito-sync:*",
        "cognito-identity:*"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
EOF
}
`, name)
}

func testAccPoolRolesAttachmentConfig_basic(name string) string {
	return fmt.Sprintf(testAccPoolRolesAttachmentConfig(name) + `
resource "aws_cognito_identity_pool_roles_attachment" "test" {
  identity_pool_id = aws_cognito_identity_pool.main.id

  roles = {
    "authenticated" = aws_iam_role.authenticated.arn
  }
}
`)
}

func testAccPoolRolesAttachmentConfig_roleMappings(name string) string {
	return fmt.Sprintf(testAccPoolRolesAttachmentConfig(name) + `
resource "aws_cognito_identity_pool_roles_attachment" "test" {
  identity_pool_id = aws_cognito_identity_pool.main.id

  role_mapping {
    identity_provider         = "graph.facebook.com"
    ambiguous_role_resolution = "AuthenticatedRole"
    type                      = "Rules"

    mapping_rule {
      claim      = "isAdmin"
      match_type = "Equals"
      role_arn   = aws_iam_role.authenticated.arn
      value      = "paid"
    }
  }

  roles = {
    "authenticated" = aws_iam_role.authenticated.arn
  }
}
`)
}

func testAccPoolRolesAttachmentConfig_roleMappingsUpdated(name string) string {
	return fmt.Sprintf(testAccPoolRolesAttachmentConfig(name) + `
resource "aws_cognito_identity_pool_roles_attachment" "test" {
  identity_pool_id = aws_cognito_identity_pool.main.id

  role_mapping {
    identity_provider         = "graph.facebook.com"
    ambiguous_role_resolution = "AuthenticatedRole"
    type                      = "Rules"

    mapping_rule {
      claim      = "isPaid"
      match_type = "Equals"
      role_arn   = aws_iam_role.authenticated.arn
      value      = "unpaid"
    }

    mapping_rule {
      claim      = "isFoo"
      match_type = "Equals"
      role_arn   = aws_iam_role.authenticated.arn
      value      = "bar"
    }
  }

  roles = {
    "authenticated" = aws_iam_role.authenticated.arn
  }
}
`)
}

func testAccPoolRolesAttachmentConfig_roleMappingsWithAmbiguousRoleResolutionError(name string) string {
	return fmt.Sprintf(testAccPoolRolesAttachmentConfig(name) + `
resource "aws_cognito_identity_pool_roles_attachment" "test" {
  identity_pool_id = aws_cognito_identity_pool.main.id

  role_mapping {
    identity_provider = "graph.facebook.com"
    type              = "Rules"

    mapping_rule {
      claim      = "isAdmin"
      match_type = "Equals"
      role_arn   = aws_iam_role.authenticated.arn
      value      = "paid"
    }
  }

  roles = {
    "authenticated" = aws_iam_role.authenticated.arn
  }
}
`)
}

func testAccPoolRolesAttachmentConfig_roleMappingsWithRulesTypeError(name string) string {
	return fmt.Sprintf(testAccPoolRolesAttachmentConfig(name) + `
resource "aws_cognito_identity_pool_roles_attachment" "test" {
  identity_pool_id = aws_cognito_identity_pool.main.id

  role_mapping {
    identity_provider         = "graph.facebook.com"
    ambiguous_role_resolution = "AuthenticatedRole"
    type                      = "Rules"
  }

  roles = {
    "authenticated" = aws_iam_role.authenticated.arn
  }
}
`)
}

func testAccPoolRolesAttachmentConfig_roleMappingsWithTokenTypeError(name string) string {
	return fmt.Sprintf(testAccPoolRolesAttachmentConfig(name) + `
resource "aws_cognito_identity_pool_roles_attachment" "test" {
  identity_pool_id = aws_cognito_identity_pool.main.id

  role_mapping {
    identity_provider         = "graph.facebook.com"
    ambiguous_role_resolution = "AuthenticatedRole"
    type                      = "Token"

    mapping_rule {
      claim      = "isAdmin"
      match_type = "Equals"
      role_arn   = aws_iam_role.authenticated.arn
      value      = "paid"
    }
  }

  roles = {
    "authenticated" = aws_iam_role.authenticated.arn
  }
}
`)
}
