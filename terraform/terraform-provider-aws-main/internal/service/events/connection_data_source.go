package events

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConnectionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceConnectionRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("name").(string))

	conn := meta.(*conns.AWSClient).EventsConn

	input := &eventbridge.DescribeConnectionInput{
		Name: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Reading EventBridge connection (%s)", d.Id())
	output, err := conn.DescribeConnection(input)
	if err != nil {
		return fmt.Errorf("error getting EventBridge connection (%s): %w", d.Id(), err)
	}

	if output == nil {
		return fmt.Errorf("error getting EventBridge connection (%s): empty response", d.Id())
	}

	log.Printf("[DEBUG] Found EventBridge connection: %#v", *output)
	d.Set("arn", output.ConnectionArn)
	d.Set("secret_arn", output.SecretArn)
	d.Set("name", output.Name)
	d.Set("authorization_type", output.AuthorizationType)
	return nil
}
