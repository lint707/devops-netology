package waf

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceWebACL() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceWebACLRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceWebACLRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn
	name := d.Get("name").(string)

	acls := make([]*waf.WebACLSummary, 0)
	// ListWebACLsInput does not have a name parameter for filtering
	input := &waf.ListWebACLsInput{}
	for {
		output, err := conn.ListWebACLs(input)
		if err != nil {
			return fmt.Errorf("error reading web ACLs: %w", err)
		}
		for _, acl := range output.WebACLs {
			if aws.StringValue(acl.Name) == name {
				acls = append(acls, acl)
			}
		}

		if output.NextMarker == nil {
			break
		}
		input.NextMarker = output.NextMarker
	}

	if len(acls) == 0 {
		return fmt.Errorf("web ACLs not found for name: %s", name)
	}

	if len(acls) > 1 {
		return fmt.Errorf("multiple web ACLs found for name: %s", name)
	}

	acl := acls[0]

	d.SetId(aws.StringValue(acl.WebACLId))

	return nil
}
