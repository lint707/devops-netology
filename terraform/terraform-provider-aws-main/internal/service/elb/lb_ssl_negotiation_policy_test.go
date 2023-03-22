package elb_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfelb "github.com/hashicorp/terraform-provider-aws/internal/service/elb"
)

func TestAccELBSSLNegotiationPolicy_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", sdkacctest.RandString(8)) // ELB name cannot be longer than 32 characters
	elbResourceName := "aws_elb.test"
	resourceName := "aws_lb_ssl_negotiation_policy.test"

	key := acctest.TLSRSAPrivateKeyPEM(2048)
	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(key, "example.com")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLBSSLNegotiationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBSSLNegotiationPolicyConfig_basic(rName, key, certificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBSSLNegotiationPolicy(elbResourceName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "attribute.#", "7"),
				),
			},
		},
	})
}

func TestAccELBSSLNegotiationPolicy_disappears(t *testing.T) {
	var loadBalancer elb.LoadBalancerDescription
	rName := fmt.Sprintf("tf-acc-test-%s", sdkacctest.RandString(8)) // ELB name cannot be longer than 32 characters
	elbResourceName := "aws_elb.test"
	resourceName := "aws_lb_ssl_negotiation_policy.test"

	key := acctest.TLSRSAPrivateKeyPEM(2048)
	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(key, "example.com")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, elb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLBSSLNegotiationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBSSLNegotiationPolicyConfig_basic(rName, key, certificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBSSLNegotiationPolicy(elbResourceName, resourceName),
					testAccCheckLoadBalancerExists(elbResourceName, &loadBalancer),
					testAccCheckLoadBalancerDisappears(&loadBalancer),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckLBSSLNegotiationPolicyDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ELBConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_elb" && rs.Type != "aws_lb_ssl_negotiation_policy" {
			continue
		}

		// Check that the ELB is destroyed
		if rs.Type == "aws_elb" {
			describe, err := conn.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{
				LoadBalancerNames: []*string{aws.String(rs.Primary.ID)},
			})

			if err == nil {
				if len(describe.LoadBalancerDescriptions) != 0 &&
					*describe.LoadBalancerDescriptions[0].LoadBalancerName == rs.Primary.ID {
					return fmt.Errorf("ELB still exists")
				}
			}

			// Verify the error
			providerErr, ok := err.(awserr.Error)
			if !ok {
				return err
			}

			if providerErr.Code() != "LoadBalancerNotFound" {
				return fmt.Errorf("Unexpected error: %s", err)
			}
		} else {
			// Check that the SSL Negotiation Policy is destroyed
			elbName, _, policyName, err := tfelb.SSLNegotiationPolicyParseID(rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = conn.DescribeLoadBalancerPolicies(&elb.DescribeLoadBalancerPoliciesInput{
				LoadBalancerName: aws.String(elbName),
				PolicyNames:      []*string{aws.String(policyName)},
			})

			if err == nil {
				return fmt.Errorf("ELB SSL Negotiation Policy still exists")
			}
		}
	}

	return nil
}

func testAccCheckLBSSLNegotiationPolicy(elbResource string, policyResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[elbResource]
		if !ok {
			return fmt.Errorf("Not found: %s", elbResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		policy, ok := s.RootModule().Resources[policyResource]
		if !ok {
			return fmt.Errorf("Not found: %s", policyResource)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ELBConn

		elbName, _, policyName, err := tfelb.SSLNegotiationPolicyParseID(policy.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.DescribeLoadBalancerPolicies(&elb.DescribeLoadBalancerPoliciesInput{
			LoadBalancerName: aws.String(elbName),
			PolicyNames:      []*string{aws.String(policyName)},
		})

		if err != nil {
			log.Printf("[ERROR] Problem describing load balancer policy '%s': %s", policyName, err)
			return err
		}

		if len(resp.PolicyDescriptions) != 1 {
			return fmt.Errorf("Unable to find policy %#v", resp.PolicyDescriptions)
		}

		attrmap := policyAttributesToMap(&resp.PolicyDescriptions[0].PolicyAttributeDescriptions)
		if attrmap["Protocol-TLSv1"] != "false" {
			return fmt.Errorf("Policy attribute 'Protocol-TLSv1' was of value %s instead of false!", attrmap["Protocol-TLSv1"])
		}
		if attrmap["Protocol-TLSv1.1"] != "false" {
			return fmt.Errorf("Policy attribute 'Protocol-TLSv1.1' was of value %s instead of false!", attrmap["Protocol-TLSv1.1"])
		}
		if attrmap["Protocol-TLSv1.2"] != "true" {
			return fmt.Errorf("Policy attribute 'Protocol-TLSv1.2' was of value %s instead of true!", attrmap["Protocol-TLSv1.2"])
		}
		if attrmap["Server-Defined-Cipher-Order"] != "true" {
			return fmt.Errorf("Policy attribute 'Server-Defined-Cipher-Order' was of value %s instead of true!", attrmap["Server-Defined-Cipher-Order"])
		}
		if attrmap["ECDHE-RSA-AES128-GCM-SHA256"] != "true" {
			return fmt.Errorf("Policy attribute 'ECDHE-RSA-AES128-GCM-SHA256' was of value %s instead of true!", attrmap["ECDHE-RSA-AES128-GCM-SHA256"])
		}
		if attrmap["AES128-GCM-SHA256"] != "true" {
			return fmt.Errorf("Policy attribute 'AES128-GCM-SHA256' was of value %s instead of true!", attrmap["AES128-GCM-SHA256"])
		}
		if attrmap["EDH-RSA-DES-CBC3-SHA"] != "false" {
			return fmt.Errorf("Policy attribute 'EDH-RSA-DES-CBC3-SHA' was of value %s instead of false!", attrmap["EDH-RSA-DES-CBC3-SHA"])
		}

		return nil
	}
}

func policyAttributesToMap(attributes *[]*elb.PolicyAttributeDescription) map[string]string {
	attrmap := make(map[string]string)

	for _, attrdef := range *attributes {
		attrmap[*attrdef.AttributeName] = *attrdef.AttributeValue
	}

	return attrmap
}

// Sets the SSL Negotiation policy with attributes.
func testAccLBSSLNegotiationPolicyConfig_basic(rName, key, certificate string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_iam_server_certificate" "test" {
  name             = "%[1]s"
  certificate_body = "%[2]s"
  private_key      = "%[3]s"
}

resource "aws_elb" "test" {
  name               = "%[1]s"
  availability_zones = [data.aws_availability_zones.available.names[0]]

  listener {
    instance_port      = 8000
    instance_protocol  = "https"
    lb_port            = 443
    lb_protocol        = "https"
    ssl_certificate_id = aws_iam_server_certificate.test.arn
  }
}

resource "aws_lb_ssl_negotiation_policy" "test" {
  name          = "foo-policy"
  load_balancer = aws_elb.test.id
  lb_port       = 443

  attribute {
    name  = "Protocol-TLSv1"
    value = "false"
  }

  attribute {
    name  = "Protocol-TLSv1.1"
    value = "false"
  }

  attribute {
    name  = "Protocol-TLSv1.2"
    value = "true"
  }

  attribute {
    name  = "Server-Defined-Cipher-Order"
    value = "true"
  }

  attribute {
    name  = "ECDHE-RSA-AES128-GCM-SHA256"
    value = "true"
  }

  attribute {
    name  = "AES128-GCM-SHA256"
    value = "true"
  }

  attribute {
    name  = "EDH-RSA-DES-CBC3-SHA"
    value = "false"
  }
}
`, rName, acctest.TLSPEMEscapeNewlines(certificate), acctest.TLSPEMEscapeNewlines(key))
}
