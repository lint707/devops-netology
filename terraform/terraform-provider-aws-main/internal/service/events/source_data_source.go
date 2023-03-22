package events

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceSource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSourceRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSourceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EventsConn

	input := &eventbridge.ListEventSourcesInput{}
	if v, ok := d.GetOk("name_prefix"); ok {
		input.NamePrefix = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Listing EventBridge sources: %s", input)

	resp, err := conn.ListEventSources(input)
	if err != nil {
		return fmt.Errorf("error listing EventBridge sources: %w", err)
	}

	if resp == nil || len(resp.EventSources) == 0 {
		return fmt.Errorf("no matching partner event source")
	}
	if len(resp.EventSources) > 1 {
		return fmt.Errorf("multiple event sources matched; use additional constraints to reduce matches to a single event source")
	}

	es := resp.EventSources[0]

	d.SetId(aws.StringValue(es.Name))
	d.Set("arn", es.Arn)
	d.Set("created_by", es.CreatedBy)
	d.Set("name", es.Name)
	d.Set("state", es.State)

	return nil
}
