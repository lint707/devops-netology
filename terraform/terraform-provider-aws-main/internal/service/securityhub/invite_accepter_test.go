package securityhub_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func testAccInviteAccepter_basic(t *testing.T) {
	resourceName := "aws_securityhub_invite_accepter.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckAlternateAccount(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, securityhub.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(t),
		CheckDestroy:             testAccCheckInviteAccepterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInviteAccepterConfig_basic(acctest.DefaultEmailAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInviteAccepterExists(resourceName),
				),
			},
			{
				Config:            testAccInviteAccepterConfig_basic(acctest.DefaultEmailAddress),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckInviteAccepterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SecurityHubConn

		resp, err := conn.GetMasterAccount(&securityhub.GetMasterAccountInput{})

		if err != nil {
			return fmt.Errorf("error retrieving Security Hub master account: %w", err)
		}

		if resp == nil || resp.Master == nil || aws.StringValue(resp.Master.AccountId) == "" {
			return fmt.Errorf("Security Hub master account not found for: %s", resourceName)
		}

		return nil
	}
}

func testAccCheckInviteAccepterDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SecurityHubConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_securityhub_invite_accepter" {
			continue
		}

		resp, err := conn.GetMasterAccount(&securityhub.GetMasterAccountInput{})
		if tfawserr.ErrCodeEquals(err, securityhub.ErrCodeResourceNotFoundException) {
			continue
		}
		// If Security Hub is not enabled, the API returns "BadRequestException"
		if tfawserr.ErrCodeEquals(err, "BadRequestException") {
			continue
		}
		if err != nil {
			return fmt.Errorf("error retrieving Security Hub master account: %w", err)
		}

		if resp == nil || resp.Master == nil || aws.StringValue(resp.Master.AccountId) == "" {
			continue
		}

		return fmt.Errorf("Security Hub master account still configured: %s", aws.StringValue(resp.Master.AccountId))
	}
	return nil
}

func testAccInviteAccepterConfig_basic(email string) string {
	return acctest.ConfigCompose(
		acctest.ConfigAlternateAccountProvider(),
		fmt.Sprintf(`
resource "aws_securityhub_invite_accepter" "test" {
  master_id = aws_securityhub_member.source.master_id

  depends_on = [aws_securityhub_account.test]
}

resource "aws_securityhub_member" "source" {
  provider = awsalternate

  account_id = data.aws_caller_identity.test.account_id
  email      = %[1]q
  invite     = true

  depends_on = [aws_securityhub_account.source]
}

resource "aws_securityhub_account" "test" {}

resource "aws_securityhub_account" "source" {
  provider = awsalternate
}

data "aws_caller_identity" "test" {}
`, email))
}
