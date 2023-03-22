package rolesanywhere_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rolesanywhere"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfrolesanywhere "github.com/hashicorp/terraform-provider-aws/internal/service/rolesanywhere"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccRolesAnywhereTrustAnchor_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	caCommonName := acctest.RandomDomainName()
	resourceName := "aws_rolesanywhere_trust_anchor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.RolesAnywhereEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustAnchorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustAnchorConfig_basic(rName, caCommonName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrustAnchorExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "source.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "source.0.source_data.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "source.0.source_data.0.acm_pca_arn", "aws_acmpca_certificate_authority.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "source.0.source_type", "AWS_ACM_PCA"),
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

func TestAccRolesAnywhereTrustAnchor_tags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	caCommonName := acctest.RandomDomainName()
	resourceName := "aws_rolesanywhere_trust_anchor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.RolesAnywhereEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustAnchorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustAnchorConfig_tags1(rName, caCommonName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrustAnchorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTrustAnchorConfig_tags2(rName, caCommonName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrustAnchorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccTrustAnchorConfig_tags1(rName, caCommonName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrustAnchorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccRolesAnywhereTrustAnchor_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	caCommonName := acctest.RandomDomainName()
	resourceName := "aws_rolesanywhere_trust_anchor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.RolesAnywhereEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustAnchorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustAnchorConfig_basic(rName, caCommonName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrustAnchorExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfrolesanywhere.ResourceTrustAnchor(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRolesAnywhereTrustAnchor_certificateBundle(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rolesanywhere_trust_anchor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.RolesAnywhereEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustAnchorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustAnchorConfig_certificateBundle(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrustAnchorExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "source.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "source.0.source_data.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "source.0.source_type", "CERTIFICATE_BUNDLE"),
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

func testAccCheckTrustAnchorDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).RolesAnywhereConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_rolesanywhere_trust_anchor" {
			continue
		}

		_, err := tfrolesanywhere.FindTrustAnchorByID(context.Background(), conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("RolesAnywhere Trust Anchor %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckTrustAnchorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No RolesAnywhere Trust Anchor ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).RolesAnywhereConn

		_, err := tfrolesanywhere.FindTrustAnchorByID(context.Background(), conn, rs.Primary.ID)

		return err
	}
}

func testAccTrustAnchorConfig_acmBase(caCommonName string) string {
	return fmt.Sprintf(`
resource "aws_acmpca_certificate_authority" "test" {
  permanent_deletion_time_in_days = 7
  type                            = "ROOT"
  certificate_authority_configuration {
    key_algorithm     = "RSA_4096"
    signing_algorithm = "SHA512WITHRSA"
    subject {
      common_name = %[1]q
    }
  }
}

data "aws_partition" "current" {}

resource "aws_acmpca_certificate" "test" {
  certificate_authority_arn   = aws_acmpca_certificate_authority.test.arn
  certificate_signing_request = aws_acmpca_certificate_authority.test.certificate_signing_request
  signing_algorithm           = "SHA512WITHRSA"

  template_arn = "arn:${data.aws_partition.current.partition}:acm-pca:::template/RootCACertificate/V1"

  validity {
    type  = "YEARS"
    value = 1
  }
}

resource "aws_acmpca_certificate_authority_certificate" "test" {
  certificate_authority_arn = aws_acmpca_certificate_authority.test.arn
  certificate               = aws_acmpca_certificate.test.certificate
  certificate_chain         = aws_acmpca_certificate.test.certificate_chain
}
`, caCommonName)
}

func testAccTrustAnchorConfig_basic(rName, caCommonName string) string {
	return acctest.ConfigCompose(
		testAccTrustAnchorConfig_acmBase(caCommonName),
		fmt.Sprintf(`
resource "aws_rolesanywhere_trust_anchor" "test" {
  name = %[1]q
  source {
    source_data {
      acm_pca_arn = aws_acmpca_certificate_authority.test.arn
    }
    source_type = "AWS_ACM_PCA"
  }
  depends_on = [aws_acmpca_certificate_authority_certificate.test]
}
`, rName))
}

func testAccTrustAnchorConfig_tags1(rName, caCommonName, tag, value string) string {
	return acctest.ConfigCompose(
		testAccTrustAnchorConfig_acmBase(caCommonName),
		fmt.Sprintf(`
resource "aws_rolesanywhere_trust_anchor" "test" {
  name = %[1]q
  source {
    source_data {
      acm_pca_arn = aws_acmpca_certificate_authority.test.arn
    }
    source_type = "AWS_ACM_PCA"
  }
  tags = {
    %[2]q = %[3]q
  }
  depends_on = [aws_acmpca_certificate_authority_certificate.test]
}
`, rName, tag, value))
}

func testAccTrustAnchorConfig_tags2(rName, caCommonName, tag1, value1, tag2, value2 string) string {
	return acctest.ConfigCompose(
		testAccTrustAnchorConfig_acmBase(caCommonName),
		fmt.Sprintf(`
resource "aws_rolesanywhere_trust_anchor" "test" {
  name = %[1]q
  source {
    source_data {
      acm_pca_arn = aws_acmpca_certificate_authority.test.arn
    }
    source_type = "AWS_ACM_PCA"
  }
  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
  depends_on = [aws_acmpca_certificate_authority_certificate.test]
}
`, rName, tag1, value1, tag2, value2))
}

func testAccTrustAnchorConfig_certificateBundle(rName string) string {
	caKey := acctest.TLSRSAPrivateKeyPEM(2048)
	caCertificate := acctest.TLSRSAX509SelfSignedCACertificateForRolesAnywhereTrustAnchorPEM(caKey)

	return fmt.Sprintf(`
resource "aws_rolesanywhere_trust_anchor" "test" {
  name = %[1]q
  source {
    source_data {
      x509_certificate_data = "%[2]s"
    }
    source_type = "CERTIFICATE_BUNDLE"
  }
}
`, rName, acctest.TLSPEMEscapeNewlines(caCertificate))
}

func testAccPreCheck(t *testing.T) {
	acctest.PreCheckPartitionHasService(names.RolesAnywhereEndpointID, t)

	conn := acctest.Provider.Meta().(*conns.AWSClient).RolesAnywhereConn

	input := &rolesanywhere.ListTrustAnchorsInput{}

	_, err := conn.ListTrustAnchors(context.Background(), input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}
