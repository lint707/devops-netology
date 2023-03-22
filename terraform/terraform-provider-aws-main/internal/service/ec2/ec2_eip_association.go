package ec2

import (
	"fmt"
	"log"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceEIPAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceEIPAssociationCreate,
		Read:   resourceEIPAssociationRead,
		Delete: resourceEIPAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"allow_reassociation": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"network_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceEIPAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	request := &ec2.AssociateAddressInput{}

	if v, ok := d.GetOk("allocation_id"); ok {
		request.AllocationId = aws.String(v.(string))
	}
	if v, ok := d.GetOk("allow_reassociation"); ok {
		request.AllowReassociation = aws.Bool(v.(bool))
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = aws.String(v.(string))
	}
	if v, ok := d.GetOk("network_interface_id"); ok {
		request.NetworkInterfaceId = aws.String(v.(string))
	}
	if v, ok := d.GetOk("private_ip_address"); ok {
		request.PrivateIpAddress = aws.String(v.(string))
	}
	if v, ok := d.GetOk("public_ip"); ok {
		request.PublicIp = aws.String(v.(string))
	}

	log.Printf("[DEBUG] EIP association configuration: %#v", request)

	var resp *ec2.AssociateAddressOutput
	err := resource.Retry(propagationTimeout, func() *resource.RetryError {
		var err error
		resp, err = conn.AssociateAddress(request)

		// EC2-VPC error for new addresses
		if tfawserr.ErrCodeEquals(err, "InvalidAllocationID.NotFound") {
			return resource.RetryableError(err)
		}

		// EC2-Classic error for new addresses
		if tfawserr.ErrMessageContains(err, "AuthFailure", "does not belong to you") {
			return resource.RetryableError(err)
		}

		if tfawserr.ErrMessageContains(err, "InvalidInstanceID", "pending instance") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if tfresource.TimedOut(err) {
		resp, err = conn.AssociateAddress(request)
	}
	if err != nil {
		return fmt.Errorf("Error associating EIP: %s", err)
	}

	log.Printf("[DEBUG] EIP Assoc Response: %s", resp)

	supportedPlatforms := meta.(*conns.AWSClient).SupportedPlatforms
	if len(supportedPlatforms) > 0 && !conns.HasEC2Classic(supportedPlatforms) && resp.AssociationId == nil {
		// We expect no association ID in EC2 Classic
		// but still error out if ID is missing and we _know_ it's NOT EC2 Classic
		return fmt.Errorf("Received no EIP Association ID in account that doesn't support EC2 Classic (%q): %s",
			supportedPlatforms, resp)
	}

	if resp.AssociationId != nil {
		d.SetId(aws.StringValue(resp.AssociationId))
	} else {
		// EC2-Classic
		d.SetId(aws.StringValue(request.PublicIp))

		if err := waitForAddressAssociationClassic(conn, aws.StringValue(request.PublicIp), aws.StringValue(request.InstanceId)); err != nil {
			return fmt.Errorf("error waiting for EC2 Address (%s) to associate with EC2-Classic Instance (%s): %w", aws.StringValue(request.PublicIp), aws.StringValue(request.InstanceId), err)
		}
	}

	return resourceEIPAssociationRead(d, meta)
}

func resourceEIPAssociationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	request, err := DescribeAddressesByID(d.Id(), meta.(*conns.AWSClient).SupportedPlatforms)
	if err != nil {
		return err
	}

	var response *ec2.DescribeAddressesOutput
	err = resource.Retry(propagationTimeout, func() *resource.RetryError {
		var err error
		response, err = conn.DescribeAddresses(request)

		if d.IsNewResource() && tfawserr.ErrCodeEquals(err, "InvalidAssociationID.NotFound") {
			return resource.RetryableError(err)
		}

		if d.IsNewResource() && (response.Addresses == nil || len(response.Addresses) == 0) {
			return resource.RetryableError(&resource.NotFoundError{})
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		response, err = conn.DescribeAddresses(request)
	}

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, "InvalidAssociationID.NotFound") {
		log.Printf("[WARN] EIP Association (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error reading EC2 Elastic IP %s: %#v", d.Get("allocation_id").(string), err)
	}

	if response.Addresses == nil || len(response.Addresses) == 0 {
		log.Printf("[INFO] EIP Association ID Not Found. Refreshing from state")
		d.SetId("")
		return nil
	}

	return readEIPAssociation(d, response.Addresses[0])
}

func resourceEIPAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	var opts *ec2.DisassociateAddressInput
	// We assume EC2 Classic if ID is a valid IPv4 address
	ip := net.ParseIP(d.Id())
	if ip != nil {
		supportedPlatforms := meta.(*conns.AWSClient).SupportedPlatforms
		if len(supportedPlatforms) > 0 && !conns.HasEC2Classic(supportedPlatforms) {
			return fmt.Errorf("Received IPv4 address as ID in account that doesn't support EC2 Classic (%q)",
				supportedPlatforms)
		}

		opts = &ec2.DisassociateAddressInput{
			PublicIp: aws.String(d.Id()),
		}
	} else {
		opts = &ec2.DisassociateAddressInput{
			AssociationId: aws.String(d.Id()),
		}
	}

	_, err := conn.DisassociateAddress(opts)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, "InvalidAssociationID.NotFound") {
			return nil
		}
		return fmt.Errorf("Error deleting Elastic IP association: %s", err)
	}

	return nil
}

func readEIPAssociation(d *schema.ResourceData, address *ec2.Address) error {
	if err := d.Set("allocation_id", address.AllocationId); err != nil {
		return err
	}
	if err := d.Set("instance_id", address.InstanceId); err != nil {
		return err
	}
	if err := d.Set("network_interface_id", address.NetworkInterfaceId); err != nil {
		return err
	}
	if err := d.Set("private_ip_address", address.PrivateIpAddress); err != nil {
		return err
	}
	if err := d.Set("public_ip", address.PublicIp); err != nil {
		return err
	}

	return nil
}

func DescribeAddressesByID(id string, supportedPlatforms []string) (*ec2.DescribeAddressesInput, error) {
	// We assume EC2 Classic if ID is a valid IPv4 address
	ip := net.ParseIP(id)
	if ip != nil {
		if len(supportedPlatforms) > 0 && !conns.HasEC2Classic(supportedPlatforms) {
			return nil, fmt.Errorf("Received IPv4 address as ID in account that doesn't support EC2 Classic (%q)",
				supportedPlatforms)
		}

		return &ec2.DescribeAddressesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("public-ip"),
					Values: []*string{aws.String(id)},
				},
				{
					Name:   aws.String("domain"),
					Values: []*string{aws.String("standard")},
				},
			},
		}, nil
	}

	return &ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("association-id"),
				Values: []*string{aws.String(id)},
			},
		},
	}, nil
}
