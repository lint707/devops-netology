package directconnect_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func testAccCheckVirtualInterfaceExists(name string, vif *directconnect.VirtualInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DirectConnectConn

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		resp, err := conn.DescribeVirtualInterfaces(&directconnect.DescribeVirtualInterfacesInput{
			VirtualInterfaceId: aws.String(rs.Primary.ID),
		})
		if err != nil {
			return err
		}

		for _, v := range resp.VirtualInterfaces {
			if aws.StringValue(v.VirtualInterfaceId) == rs.Primary.ID {
				*vif = *v

				return nil
			}
		}

		return fmt.Errorf("Direct Connect virtual interface (%s) not found", rs.Primary.ID)
	}
}

func testAccCheckVirtualInterfaceDestroy(s *terraform.State, t string) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).DirectConnectConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != t {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		resp, err := conn.DescribeVirtualInterfaces(&directconnect.DescribeVirtualInterfacesInput{
			VirtualInterfaceId: aws.String(rs.Primary.ID),
		})
		if tfawserr.ErrMessageContains(err, directconnect.ErrCodeClientException, "does not exist") {
			continue
		}
		if err != nil {
			return err
		}

		for _, v := range resp.VirtualInterfaces {
			if aws.StringValue(v.VirtualInterfaceId) == rs.Primary.ID && aws.StringValue(v.VirtualInterfaceState) != directconnect.VirtualInterfaceStateDeleted {
				return fmt.Errorf("[DESTROY ERROR] Direct Connect virtual interface (%s) not deleted", rs.Primary.ID)
			}
		}
	}

	return nil
}
