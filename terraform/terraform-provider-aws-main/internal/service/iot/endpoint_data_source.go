package iot

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEndpointRead,
		Schema: map[string]*schema.Schema{
			"endpoint_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"iot:CredentialProvider",
					"iot:Data",
					"iot:Data-ATS",
					"iot:Jobs",
				}, false),
			},
		},
	}
}

func dataSourceEndpointRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IoTConn
	input := &iot.DescribeEndpointInput{}

	if v, ok := d.GetOk("endpoint_type"); ok {
		input.EndpointType = aws.String(v.(string))
	}

	output, err := conn.DescribeEndpoint(input)
	if err != nil {
		return fmt.Errorf("error while describing iot endpoint: %w", err)
	}
	endpointAddress := aws.StringValue(output.EndpointAddress)
	d.SetId(endpointAddress)
	if err := d.Set("endpoint_address", endpointAddress); err != nil {
		return fmt.Errorf("error setting endpoint_address: %w", err)
	}
	return nil
}
